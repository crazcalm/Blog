+++
draft = false
date = 2018-08-01
tags = ["golang", "file server"]
description = ""
title = "Using Embedding to Create a Custom File Server"
highlight = true
css = []
scripts = []
+++

## What is Embedding?

The [talk given by Sean Kelly on Embedding](https://www.youtube.com/watch?time_continue=1&v=-LzYjMzfGDQ) is the best resource I have found for explaining it.

In summary, Embedding, in golang, is the rules governing the ability to place one struct or interface within another struct or interface and, from the outer struct or interface, call the exported fields of the inner struct or interface.

## In terms of creating a custom file server, why is this important?

In order to create a custom file server we must implement the http.FileSystem:

	package http // import "net/http"

	type FileSystem interface {
	        Open(name string) (File, error)
	}
	    A FileSystem implements access to a collection of named files. The elements
	    in a file path are separated by slash ('/', U+002F) characters, regardless
	    of host operating system convention.

The easiest way to do this is to take the wanted interface and embed it in our own struct. Then we take an existing struct that implement this interfaces and pass it into our struct on creation. This will allow our struct to call the exported fields of the passed in struct that satisfies our embedded interface.

For example:

	package main

	import (
	        "log"
	        "net/http"
	        "os"
	)

	type myFileSystem struct {
	        http.FileSystem
	}

	func main() {
	        home := os.Getenv("HOME")
	        customFileSystem := myFileSystem{http.Dir(home)}
	        http.Handle("/", http.FileServer(customFileSystem))
	        log.Fatal(http.ListenAndServe(":12346", nil))
	}

The above code, creates a file server using our custom FileSystem. As specified in the interface, a FileSystem has one method called Open.

As of right now, myFileSystem does not have a method called Open. However, myFileSystem's embedded interface does have a method called Open and the struct we passed in (http.Dir) has a method called Open. Since myFileSystem does not have a method called Open and the embedded interface/struct does, the embedded method Open gets promoted up to myFileSystem, which means that myFileSystem.Open uses http.Dir.Open when called.

## How to deny access to dot files

The signature for http.Dir is:

	func (d Dir) Open(name string) (File, error)

The method Open takes in the name of a file (or directory) and returns the *os.File of that file (or directory).

If we want to deny access to dot files, we can do the following:

	package main

	import (
	        "log"
	        "net/http"
	        "os"
	        "strings"
	)

	//isDotFile -- checks to see if name is a dot file or in a dot directory
	func isDotFile(name string) (result bool) {
	        parts := strings.Split(name, "/")
	        for _, part := range parts {
	                if strings.HasPrefix(part, ".") {
	                        result = true
	                        return
	                }
	        }
	        return
	}

	type myFileSystem struct {
	        http.FileSystem
	}

	func (fs myFileSystem) Open(name string) (http.File, error) {
	        file, err := fs.FileSystem.Open(name)

	        if isDotFile(name) { //If dot file, return 403 response
	                return file, os.ErrPermission
	        }
	        return file, err
	}

	func main() {
	        home := os.Getenv("HOME")
	        customFileSystem := myFileSystem{http.Dir(home)}
	        http.Handle("/", http.FileServer(customFileSystem))
	        log.Fatal(http.ListenAndServe(":12346", nil))
	}

Now that myFileSystem has its own Open method, the embedded interface/struct's method Open cannot be promoted because that would cause ambiguity (essentially, we would not know which Open method we were calling).

In our Open method, we have called the Open method attached to the embedded interface/struct by using the name of the embedded interface/struct:

	file, err := fs.FileSystem.Open(name)
	
We do this so that we can use http.Dir's Open method, which we know already satisfies the http.FileSystem interface.

We then pass name into a function called isDotFile, which iterates  over the file name (the path to the file), and returns true if the file is, or is within, a dot file.

If isDotFile returns true, we then return an os.ErrPermission error, which will show the user "403 Forbidden" page.

If isDotFile returns false, then myFileSystem.Open works just like it did before.

## How to hide dot files

So far, we have denied access to dot files, but we can still see them. In order to hide the dot files so that they do not appear on the page, we must use the same trick we used to create a custom http.FileSystem to create a custom http.File.

	type myFile struct {
		http.File
	}

	func (f myFile) Readdir(n int) (wantedFiles []os.FileInfo, err error) {
		files, err := f.File.Readdir(n)
		for _, file := range files { // Filters out the dot files
			if !strings.HasPrefix(file.Name(), ".") {
				wantedFiles = append(wantedFiles, file)
			}
		}
		return
	}


We have embedded httpFile into our struct (called myFile) and created a Readdir method for it.

We created the Readdir method because the file server uses http.File.Readdir to obtain the contents of a directory. Dot files are typically included in these contents, but, in creating our own Readdir method, we can filter that content to exclude dot files. Thus, our Readdir method calls the embedded interface/structs Readdir method, filters their results to only include the content that we want and then returns the filtered content.

The finished program looks like this:

	package main

	import (
	        "log"
	        "net/http"
	        "os"
	        "strings"
	)

	//isDotFile -- checks to see if name is a dot file or in a dot directory
	func isDotFile(name string) (result bool) {
	        parts := strings.Split(name, "/")
	        for _, part := range parts {
	                if strings.HasPrefix(part, ".") {
	                        result = true
	                        return
	                }
	        }
	        return
	}

	type myFile struct {
	        http.File
	}

	func (f myFile) Readdir(n int) (wantedFiles []os.FileInfo, err error) {
	        files, err := f.File.Readdir(n)
	        for _, file := range files { // Filters out the dot files
	                if !strings.HasPrefix(file.Name(), ".") {
	                        wantedFiles = append(wantedFiles, file)
	                }
	        }
	        return
	}

	type myFileSystem struct {
	        http.FileSystem
	}

	func (fs myFileSystem) Open(name string) (http.File, error) {
	        file, err := fs.FileSystem.Open(name)

	        if isDotFile(name) { //If dot file, return 403 response
	                return file, os.ErrPermission
	        }
	        return myFile{file}, err
	}

	func main() {
	        home := os.Getenv("HOME")
	        customFileSystem := myFileSystem{http.Dir(home)}
	        http.Handle("/", http.FileServer(customFileSystem))
	        log.Fatal(http.ListenAndServe(":12346", nil))
	}



