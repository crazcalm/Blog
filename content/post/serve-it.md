+++
draft = false
date = 2018-07-14
tags = ["golang", "file server", "rabbit hole"]
description = "My path and rabbit hole when creating serve-it."
title = "How to modify Golang's default FileServer to hide dot files"
highlight = true
css = []
scripts = []
+++
## Why? Why do such a thing

I have no good reason... I was poking around in the golang net/http package one night and came across triv.go. Triv.go (triv is short for trivial) is a file filled with trivial server examples.

One of those examples was a file server. When I pointed my browser at that file server, I saw all of my dot files. At the time, I was sharing my home directory, which is filled with dot files. In order to see everything else, I had to scroll down a couple of times.

So I thought to myself, "There should be a way to hide your dot files, right?"

That is when I started my search.

## Before we start, some background knowledge

A file server, in this context, is a server that allows other machines on the network access to the contents of its file system. Typically, you serve a directory, which give others access to everything in that directory (including nested directories).

Example:

	package main

	import (
		"log"
		"net/http"
	)

	func main() {
		http.Handle("/", http.FileServer(http.Dir(".")))
		log.Fatal(http.ListenAndServe(":12346", nil))
	}

The example will serve all of the contents of the current directory on localhost:12346

Back to the story!

# Is there a config option to hide dot files? To the docs!

go doc http.FileServer:

	package http // import "net/http"

	func FileServer(root FileSystem) Handler
	    FileServer returns a handler that serves HTTP requests with the contents of
	    the file system rooted at root.

	    To use the operating system's file system implementation, use http.Dir:

	    http.Handle("/", http.FileServer(http.Dir("/tmp")))

	    As a special case, the returned file server redirects any request ending in
	    "/index.html" to the same path, without the final "index.html".


There is no mention of a flag option or box that says, "Check here to hid dot files", but the http.Dir looks promising. Dir probably stands for directory, which probably controls the directory list, which is the thing that we want to modify.

Let's check it out!

go doc http.Dir:

	package http // import "net/http"

	type Dir string
	    A Dir implements FileSystem using the native file system restricted to a
	    specific directory tree.

	    While the FileSystem.Open method takes '/'-separated paths, a Dir's string
	    value is a filename on the native file system, not a URL, so it is separated
	    by filepath.Separator, which isn't necessarily '/'.

	    Note that Dir will allow access to files and directories starting with a
	    period, which could expose sensitive directories like a .git directory or
	    sensitive files like .htpasswd. To exclude files with a leading period,
	    remove the files/directories from the server or create a custom FileSystem
	    implementation.

	    An empty Dir is treated as ".".


Uh Oh... This explicitly states that, if I do not want to show my dot files, I either have to remove them from the directory in question or create a custom FileSystem Implementation.

I'll bite the bullet. What is this FileSystem thing.

go doc http.FileSystem:

	package http // import "net/http"

	type FileSystem interface {
	        Open(name string) (File, error)
	}
	    A FileSystem implements access to a collection of named files. The elements
	    in a file path are separated by slash ('/', U+002F) characters, regardless
	    of host operating system convention.

That is not very helpful... What is a file?

go doc http.File:

	package http // import "net/http"

	type File interface {
	        io.Closer
	        io.Reader
	        io.Seeker
	        Readdir(count int) ([]os.FileInfo, error)
	        Stat() (os.FileInfo, error)
	}
    		A File is returned by a FileSystem's Open method and can be served by the
    		FileServer implementation.

    		The methods should behave the same as those on an *os.File.

That is pretty cryptic. All in all, I am still really lost as to how this helps me.

## Time to look at the source code!

FileServer code is in net/http/fs.go

	func FileServer(root FileSystem) Handler {
 		return &fileHandler{root}
	}

That is a pointer to something that is not exported, which means that I will not be able to directly call any of their code when I start writing my own custom FileServer...

Rabbit hole time!

	type fileHandler struct {
        root FileSystem
	}

	func (f *fileHandler) ServeHTTP(w ResponseWriter, r *Request) {
        upath := r.URL.Path
        if !strings.HasPrefix(upath, "/") {
                upath = "/" + upath
                r.URL.Path = upath
        }
        serveFile(w, r, f.root, path.Clean(upath), true)
	}


Nothing about files here. Next stop, serveFile!


	// name is '/'-separated, not filepath.Separator.
	func serveFile(w ResponseWriter, r *Request, fs FileSystem, name string, 	redirect bool) {
        const indexPage = "/index.html"

        // redirect .../index.html to .../
        // can't use Redirect() because that would make the path absolute,
        // which would be a problem running under StripPrefix
        if strings.HasSuffix(r.URL.Path, indexPage) {
                localRedirect(w, r, "./")
                return
        }

        f, err := fs.Open(name)
        if err != nil {
                msg, code := toHTTPError(err)
                Error(w, msg, code)
                return
        }
        defer f.Close()

        d, err := f.Stat()
        if err != nil {
                msg, code := toHTTPError(err)
                Error(w, msg, code)
                return
        }

        if redirect {
                // redirect to canonical path: / at end of directory url
                // r.URL.Path always begins with /
                url := r.URL.Path
                if d.IsDir() {
                        if url[len(url)-1] != '/' {
                                localRedirect(w, r, path.Base(url)+"/")
                                return
                        }
                } else {
                        if url[len(url)-1] == '/' {
                                localRedirect(w, r, "../"+path.Base(url))
                                return
                        }
                }
        }

        // redirect if the directory name doesn't end in a slash
        if d.IsDir() {
                url := r.URL.Path
                if url[len(url)-1] != '/' {
                        localRedirect(w, r, path.Base(url)+"/")
                        return
                }
        }

        // use contents of index.html for directory, if present
        if d.IsDir() {
                index := strings.TrimSuffix(name, "/") + indexPage
                ff, err := fs.Open(index)
                if err == nil {
                        defer ff.Close()
                        dd, err := ff.Stat()
                        if err == nil {
                                name = index
                                d = dd
                                f = ff
                        }
                }
        }

        // Still a directory? (we didn't find an index.html file)
        if d.IsDir() {
                if checkIfModifiedSince(r, d.ModTime()) == condFalse {
                        writeNotModified(w)
                        return
                }
                w.Header().Set("Last-Modified", d.ModTime().UTC().Format(TimeFormat))
                dirList(w, r, f)
                return
        }

        // serveContent will check modification time
        sizeFunc := func() (int64, error) { return d.Size(), nil }
        serveContent(w, r, d.Name(), d.ModTime(), sizeFunc, f)
	}


This function checks for index.html (if present, serve it). It has a number of other checks to see if the requested url path is correct (if slightly malformed, it will redirect to the right path). And it handles our case of, if this is a directory with no index.html file (and the path is correct), what should we do?

The snippet we care about:

		// Still a directory? (we didn't find an index.html file)
        if d.IsDir() {
                if checkIfModifiedSince(r, d.ModTime()) == condFalse {
                        writeNotModified(w)
                        return
                }
                w.Header().Set("Last-Modified", d.ModTime().UTC().Format(TimeFormat))
                dirList(w, r, f)
                return
        }


dirList! That sounds like directory list, which sounds like a fresh glass of water in the desert (AKA Hope!).

	func dirList(w ResponseWriter, r *Request, f File) {
        dirs, err := f.Readdir(-1)
        if err != nil {
                logf(r, "http: error reading directory: %v", err)
                Error(w, "Error reading directory", StatusInternalServerError)
                return
        }
        sort.Slice(dirs, func(i, j int) bool { return dirs[i].Name() < dirs[j].Name() })

        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        fmt.Fprintf(w, "<pre>\n")
        for _, d := range dirs {
                name := d.Name()
                if d.IsDir() {
                        name += "/"
                }
                // name may contain '?' or '#', which must be escaped to remain
                // part of the URL path, and not indicate the start of a query
                // string or fragment.
                url := url.URL{Path: name}
                fmt.Fprintf(w, "<a href=\"%s\">%s</a>\n", url.String(), htmlReplacer.Replace(name))
        }
        fmt.Fprintf(w, "</pre>\n")
	}

Finally! These are the lines that we want to change!

We see that the list of files and directories comes from f.Readdir, Which means that if we were to modify that function to return non-dot file items, then we would be done!

However, we should note that http.File (which is what is being passed into dirList) is an interface, but there is nothing in the net/http/fs.go file that satisfies that interface. And f.Readdir() seems a lot like os.File.Readdir...

go doc os.File.Readdir:

	func (f *File) Readdir(n int) ([]FileInfo, error)
    		Readdir reads the contents of the directory associated with file and returns
    		a slice of up to n FileInfo values, as would be returned by Lstat, in
    		directory order. Subsequent calls on the same file will yield further
    		FileInfos.

    		If n > 0, Readdir returns at most n FileInfo structures. In this case, if
    		Readdir returns an empty slice, it will return a non-nil error explaining
    		why. At the end of a directory, the error is io.EOF.

    		If n <= 0, Readdir returns all the FileInfo from the directory in a single
    		slice. In this case, if Readdir succeeds (reads all the way to the end of
    		the directory), it returns the slice and a nil error. If it encounters an
    		error before the end of the directory, Readdir returns the FileInfo read
    		until that point and a non-nil error.

I am not touching that.

The next obvious choice is to modify the for loop by checking the name of the file/directory. If the name starts with a dot, skip the line that would write its name to the http response.

	func dirList(w http.ResponseWriter, r *http.Request, f http.File) {
        dirs, err := f.Readdir(-1)
        if err != nil {
                logf(r, "http: error reading directory: %v", err)
                http.Error(w, "Error reading directory", http.StatusInternalServerError)
                return
        }
        sort.Slice(dirs, func(i, j int) bool { return dirs[i].Name() < dirs[j].Name() })

        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        fmt.Fprintf(w, "<pre>\n")
        for _, d := range dirs {
                name := d.Name()

                //Added by Marcus
                if !*showDotFiles && strings.HasPrefix(name, ".") {
                        continue
                }

                if d.IsDir() {
                        name += "/"
                }
                // name may contain '?' or '#', which must be escaped to remain
                // part of the URL path, and not indicate the start of a query
                // string or fragment.
                url := url.URL{Path: name}
                fmt.Fprintf(w, "<a href=\"%s\">%s</a>\n", url.String(), htmlReplacer.Replace(name))
        }
        fmt.Fprintf(w, "</pre>\n")
	}

showDotFiles is a pointer to a boolean. If showDotFiles == false && name starts with ".", we continue on to the next item in the loop, which will skip to the next interation of the loop and thereby not write the dot file to the http response.

## (Theoreticaly) Done!

If all you want to do is create a FileServer that hides dot files AND you are okay with modifiying the standard library, add a boolean, an if statement, and a continue, and you are golden.

If you find it taboo to modify the standard library (like I do), you copy all the code needed for the default FileServer into one file and make your changes there.

That sounds easy enough until you realize that the FileServer implementation depends on 1,000 lines of code. Most of the code is from net/http/fs.go, but some of it comes from other files in that package.

Once you do that, you will end up with something like this -- [serve-it source code.](https://github.com/crazcalm/serve-it/blob/master/serve-it.go).

After all this digging and copy and pasting (I did it one function at a time to make sure I did not have any unneeded lines), I came to the realization that should of been obvious... Just because something is hidden does not mean you cannot access it.

But that has a straight forward fix.

	func (f *fileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
        upath := r.URL.Path
        fmt.Printf("Requested path: %s\n", upath)
        if !strings.HasPrefix(upath, "/") {
                upath = "/" + upath
                r.URL.Path = upath
        }

        //Added by Marcus
        if !*showDotFiles {
                pathParts := strings.Split(r.URL.Path, "/")
                for _, part := range pathParts {
                        if strings.HasPrefix(part, ".") {
                                http.Error(w, "403 Forbidden", http.StatusForbidden)
                                return
                        }
                }
        }

        serveFile(w, r, f.root, path.Clean(upath), true)
	}


Before we hit serveFile, we will check the requested url. If it has a dot file in it, return a 403 Forbidden response and call it a day.

<image src="/img/theEndDab.jpg">