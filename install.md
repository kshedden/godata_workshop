Go toolchain installation and configuration
===========================================

## Installation and setup

The Go project and tools are open source and community-driven.
To compile and run Go code on your computer, you will need to download
and install the "Go tool", which can be found
[here](https://golang.org/dl).  More detailed installation
instructions are [here](https://golang.org/doc/install).

Go source files can be written in a text editor, or using various
IDE's (integrated development environments).  There is no official or
predominant IDE or editor for Go.  [This
screencast](https://www.youtube.com/watch?v=XCsL89YtqCs) walks through
the standard uses of the Go tool.


Greatlakes
----------

* Connect to Greatlakes (more information
  [here](http://arc-ts.umich.edu/greatlakes/user-guide))

* Type `module load go`

* Create a directory called `go` at the top of your home
 directory.

Your own machine
----------------

* Download the latest version of the Go tools
  [here](https://golang.org/dl) and read the [installation
  instructions](https://golang.org/doc/install)

Goimports utility and editor configuration (optional)
-----------------------------------------------------

* Type `go get golang.org/x/tools/cmd/goimports` to get the
  [goimports](https://godoc.org/golang.org/x/tools/cmd/goimports) tool

See the [goimports page](https://godoc.org/golang.org/x/tools/cmd/goimports) for
instructions on configuring in your editor.
