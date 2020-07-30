# Gox - Simple Go Cross Compilation

Gox is a simple, no-frills tool for Go cross compilation that behaves a
lot like standard `go build`. Gox will parallelize builds for multiple
platforms. Gox will also build the cross-compilation toolchain for you.

## Installation

To install Gox, please use `go get`. We tag versions so feel free to
checkout that tag and compile.

```
$ go get github.com/GoLandr/gox
...
$ gox -h
...
```

## Usage

If you know how to use `go build`, then you know how to use Gox. For
example, to build the current package, specify no parameters and just
call `gox`. Gox will parallelize based on the number of CPUs you have
by default and build for every platform by default:

```
$ gox
Number of parallel builds: 4

-->      darwin/386: github.com/GoLandr/gox
-->    darwin/amd64: github.com/GoLandr/gox
-->       linux/386: github.com/GoLandr/gox
-->     linux/amd64: github.com/GoLandr/gox
-->       linux/arm: github.com/GoLandr/gox
-->     freebsd/386: github.com/GoLandr/gox
-->   freebsd/amd64: github.com/GoLandr/gox
-->     openbsd/386: github.com/GoLandr/gox
-->   openbsd/amd64: github.com/GoLandr/gox
-->     windows/386: github.com/GoLandr/gox
-->   windows/amd64: github.com/GoLandr/gox
-->     freebsd/arm: github.com/GoLandr/gox
-->   android/arm64: github.com/GoLandr/gox
-->      netbsd/386: github.com/GoLandr/gox
-->    netbsd/amd64: github.com/GoLandr/gox
-->      netbsd/arm: github.com/GoLandr/gox
-->       plan9/386: github.com/GoLandr/gox
```

Or, if you want to build a package and sub-packages:

```
$ gox ./...
...
```

Or, if you want to build multiple distinct packages:

```
$ gox github.com/GoLandr/gox github.com/hashicorp/serf
...
```

Or if you want to just build for linux:

```
$ gox -os="linux"
...
```

Or maybe you just want to build for 64-bit linux:

```
$ gox -osarch="linux/amd64"
...
```

And more! Just run `gox -h` for help and additional information.
## Command User Guide

Usage: `gox [options] [packages]`

  Gox cross-compiles Go applications in parallel.

  If no specific operating systems or architectures are specified, Gox
  will build for all pairs supported by your version of Go.

Options:

      -arch=""            Space-separated list of architectures to build for
      -build-toolchain    Build cross-compilation toolchain
      -c-cross-compilers  Set custom C cross-compilers for platforms if CGO is enabled
      -cgo                Sets CGO_ENABLED=1, requires proper C toolchain (advanced)
      -gcflags=""         Additional '-gcflags' value to pass to go build
      -ldflags=""         Additional '-ldflags' value to pass to go build
      -asmflags=""        Additional '-asmflags' value to pass to go build
      -tags=""            Additional '-tags' value to pass to go build
      -mod=""             Additional '-mod' value to pass to go build
      -os=""              Space-separated list of operating systems to build for
      -osarch=""          Space-separated list of os/arch pairs to build for
      -armarch=""         Space-separated list of GOARM arch version to build for when arch is "arm"
      -osarch-list        List supported os/arch pairs for your Go version
      -output="foo"       Output path template. See below for more info
      -parallel=-1        Amount of parallelism, defaults to number of CPUs
      -gocmd="go"         Build command, defaults to Go
      -rebuild            Force rebuilding of package that were up to date
      -verbose            Verbose mode

Output path template:

  The output path for the compiled binaries is specified with the
  `-output` flag. The value is a string that is a Go text template.
  The default value is `{{.Dir}}_{{.OS}}_{{.Arch}}`. The variables and
  their values should be self-explanatory.

Platforms (OS/Arch):

  The operating systems and architectures to cross-compile for may be
  specified with the `-arch` and `-os` flags. These are space separated lists
  of valid GOOS/GOARCH values to build for, respectively. You may prefix an
  OS or Arch with `!` to negate and not build for that platform. If the list
  is made up of only negations, then the negations will come from the default
  list.

  Additionally, the `-osarch` flag may be used to specify complete os/arch
  pairs that should be built or ignored. The syntax for this is what you would
  expect: `darwin/amd64` would be a valid osarch value. Multiple can be space
  separated. An os/arch pair can begin with `!` to not build for that platform.

  The `-osarch` flag has the highest precedent when determing whether to
  build for a platform. If it is included in the `-osarch` list, it will be
  built even if the specific os and arch is negated in `-os` and `-arch`,
  respectively.

Platform Overrides:

  The `-gcflags`, `-ldflags` and `-asmflags` options can be overridden per-platform
  by using environment variables. Gox will look for environment variables
  in the following format and use those to override values if they exist:

    GOX_[OS]_[ARCH]_GCFLAGS
    GOX_[OS]_[ARCH]_LDFLAGS
    GOX_[OS]_[ARCH]_ASMFLAGS

C cross-compilers:

  It is possible to set C cross-compilers by platforms when CGO is enabled.
  The format of setting a compiler for a platform is the following:
  { platform }={ compiler }. To configure multiple compilers for multiple
  platforms separate each setting by a comma.
  Example: `-c-cross-compilers="linux/arm=arm-linux-gnueabi-gcc-6"`


## Versus Other Cross-Compile Tools

A big thanks to these other options for existing. They each paved the
way in many aspects to make Go cross-compilation approachable.

* [Dave Cheney's golang-crosscompile](https://github.com/davecheney/golang-crosscompile) -
  Gox compiles for multiple platforms and can therefore easily run on
  any platform Go supports, whereas Dave's scripts require a shell. Gox
  will also parallelize builds. Dave's scripts build sequentially. Gox has
  much easier to use OS/Arch filtering built in.

* [goxc](https://github.com/laher/goxc) -
  A very richly featured tool that can even do things such as build system
  packages, upload binaries, generate download webpages, etc. Gox is a
  super slim alternative that only cross-compiles binaries. Gox builds packages in parallel, whereas
  goxc doesn't. Gox doesn't enforce a specific output structure for built
  binaries.

