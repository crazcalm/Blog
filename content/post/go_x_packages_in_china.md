+++
scripts = []
css = []
highlight = true
title = "Golang: Using Google's x packages in China"
description = ""
tags = []
draft = false
date = "2017-10-02T15:33:17+08:00"

+++
# The Issue

While in China, I do not have access to Google services, websites, or platforms. In terms of being a Golang user, this is annoying because a ton of Golang projects and packages use the `x` packages, which are a set of Golang packages that were built by Google. These packages are hosted on a Google server, so I end up seeing a lot of this:

	
	$ go get -u github.com/justjanne/powerline-go
	
	cd /home/crazcalm/.gvm/pkgsets/go1.8.3/global/src/golang.org/x/sys; git pull --ff-only
	
	fatal: unable to access 'https://go.googlesource.com/sys/': Failed to connect to go.googlesource.com port 443: Connection timed out
	
	package golang.org/x/sys/unix: exit status 1
	
	package golang.org/x/text/width: unrecognized import path "golang.org/x/text/width" (https fetch: Get https://golang.org/x/text/width?go-get=1: dial tcp 216.239.37.1:443: i/o timeout)




# The Solution

All of the `x` packages can be found on github. The below links will direct you to the Golang organization page on github (filtered for Go projects):

- [https://github.com/golang?utf8=%E2%9C%93&q=&type=&language=go](https://github.com/golang?utf8=%E2%9C%93&q=&type=&language=go)

Once there, you may search for the `x` package that you want. The README of each `x` package will have instructions on how to manually install the package. In summary, those instructions are to clone the github repo into the following directory on your computer:

	$GOPATH/src/golang.org/x/

