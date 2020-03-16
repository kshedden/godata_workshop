
Some comments on why you might want to learn Go
-----------------------------------------------

Go is a relatively new open source programming language.  It originated
within Google around 2006, and became stable around 2012.  Although developed
and primarily funded by Google, the language specification and all Go tools
are free and open source.

One reason to learn Go (or any other programming language)
is that it is an intellectual challenge, and will make you a better
programmer in any other language.  Beyond that, Go is quite fast, and is arguably a good
tool for general-purpose programming.  Since it is compiled and strongly typed, you can write loops
in it.  It may be viewed as a modern "low level" language that delivers 80% of the performance
of C with 80% of the comfort of Python.

It's not easy to characterize a programming language in a few words, but here are some
of the features that describe Go:

* Minimalist syntax

* "Safe", garbage collected; no VM (virtual machine), minimal runtime

* Strongly typed, with a provision called "interfaces" for semi or full dynamic typing

* Source files compile to static native binaries (alse can create libraries)

* Simple, minimalist, robust, fast tools

* Strong support for concurrency

* Comprehensive and robust standard library

* No "metaprogramming", generics, parametric polymorphism, algebraic types,
and other concepts that some people consider to be too "magical".

* The language syntax is not extensible, so everyone's code should look similar.

* "Orthogonality" -- there is usually only one "right" way to do something

* Rigid source code formatting, naming, and commenting conventions (optionally enforced by tools)

* Go moves away from the idea of a scripting language that is "glue for C".  Go is fast enough
(in most cases) to implement algorithms, but easy enough to use for scripting.

## What is Go good for?  What was it designed for and how is it mainly used?

* "Systems" and "infrastructure" -- code that needs to run fast, reliably, without excessive
demand on the system (memory and CPU)

* Distributed "services" (microservices)

* Code that "less advanced" programmers can learn to read, write, and maintain -- most people
can become an "expert" in Go within a year, other lower-level languages may take much longer to master.

* Code that can be statically analyzed for correctness


## What is Go not good for (currently or ever)?

* Go is not a "domain specific language" for statistics, mathematics, engineering, etc.

* At present, Go is not widely used for academic data science and statistics

* Some excellent libraries exist for scientific and other work.  But there are many gaps
compared to, e.g. Python.
