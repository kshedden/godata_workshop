Go for data processing
======================

This site contains materials for a workshop on using
[Go](http://golang.org) for data processing.

Here are some subjective [comments](why.md) on why you might want to learn Go.

Go is an excellent "utility language" that is being adopted by organizations that
manage large volumes of data.  It is also being used to develop tools for data
management and system maintenance.  Important utilities for maintaining cloud
infrastructure such as [Kubernetes](http://kubernetes.io) and
[Docker](http://docker.com) are written in Go.

Every programming language involves tradeoffs.
Go is stricter and simpler than Python, but less strict than Rust and less
complex than either Rust or Julia.  The simplicity of Go makes it easy
to maintain and refactor Go code, and generally makes Go code easy to
read and understand.  On the other hand, this simplicity means that
it can be necessary to write more "boilerplate code" than in other languages.
Some people feel that it is harder to develop an an "abstraction" of the problem
that you are trying to solve when using Go, but other people argue that Go supports
abstractions, just in a different way than, say, C++ or Java.
While
Go is much faster than Python, it is not quite as fast as C, C++, or Rust.

At present, Go is one of the 10 most used programming lanuages,
according to Tiobe.  However it is not commonly encountered in academic or research settings,
especially in data science.  This may be due to the success of the Python/C++
ecosystem (using C++ to implement algorithms and Python to script them).
However it is not uncommon to encounter tasks where no appropriate C++ library routine exists,
and where Python would be very slow.  At present, this may be the best
use-case for Go.  More important than that, learning Go will make anyone a better
overall programmer.  Since it is somewhat closer to the hardware than, say, Python,
it will help you understand how computers and compilers work,
and will reveal some of the trade-offs in programming language design.


Compiling and running a simple Go program
-----------------------------------------

Go scripts can be placed anywhere in your file system.  Go packages
that are able to be used by other Go programs should be placed in your
GOHOME directory, which defaults to the directory named "go" in the
top level of your home directory, i.e. `~/go`.

Go source files are text files with suffix ".go".  A very simple Go
program is:

```
package main

import "fmt"

func main() {
    fmt.Printf("This is a Go program...\n")
}
```

Suppose you save this in a file called `simple.go` in your current
working directory.  You should be able to run it using the command `go
run simple.go`.  Alternatively, you can compile it to an executable
using `go build simple.go`.  This will produce a binary executable
called `simple` in your local directory, which you can then run
using `./simple`.


Go resources
------------

* Use Go in your browser at the [Go playground](https://play.golang.org/)

* Use Go in your browser at [Repl.It](https://repl.it/repls)

* [Go project web site](http://golang.org)

* [Effective Go](https://golang.org/doc/effective_go.html)

* [Go tour](https://tour.golang.org/list)

* The Go [standard library](https://golang.org/pkg/)

* [Godoc](https://godoc.org) module documentation

* The Go project [wiki](https://github.com/golang/go/wiki)

* [High performance Go](https://dave.cheney.net/high-performance-go)

* Rob Pike Go talks on youtube: [Another Golang at language
design](https://www.youtube.com/watch?v=7VcArS4Wpqk), [Simplicity is
complicated](https://www.youtube.com/watch?v=rFejpH_tAHM), [Go
concurrency patterns](https://www.youtube.com/watch?v=f6kdp27TYZs),
[Go proverbs](https://www.youtube.com/watch?v=PAAkCSZUG1c),
[Concurrency is not
parallelism](https://www.youtube.com/watch?v=B9lP-E4J_lc),
[Go 2](https://www.youtube.com/watch?v=RIvL2ONhFBI)

* [Installation and configuration of the Go tools](install.md)
