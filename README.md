## chanstream

[![Build Status](https://travis-ci.org/gdamore/chanstream.svg?branch=master)](https://travis-ci.org/gdamore/chanstream) [![GoDoc](https://godoc.org/github.com/gdamore/chanstream?status.png)](https://godoc.org/github.com/gdamore/chanstream)

Package chanstream implements an API compatible with and similiar to the TCP
connection (and net.Conn as well) API, on top of Go channels.  This is in
pure Go, without any external dependencies.

The intention is to facilitate the use of channels for intra-program
communication, in a manner similiar to TCP or Unix Domain sockets, without
creating any externally visible addresses or service points.  (This can also
be done more efficiently, since data need not be copied to the kernel.)

An observer might wonder why not just utilize Go channels directly?  The
rationale here is that this allows abstraction layers to be built on top
of channels that can choose to use channels (via chanstream) or TCP or
other transports.  This can make it possible to eliminate certain special
cases in program handling.

In particular, this package was developed to support an effort to produce
a pure Go implementation of nanomsg and/or ZeroMQ.  This package can be used
as the underlying transport for the inproc: scheme.

Consider this a work-in-progress, and use at your own risk.

## Installing

### Using *go get*

    $ go get github.com/gdamore/chanstream

After this command *chanstream* is ready to use. Its source will be in:

    $GOROOT/src/pkg/github.com/gdamore/chanstream

You can use `go get -u -a` to update all installed packages.

## Documentation

For docs, see http://godoc.org/github.com/gdamore/chanstream or run:

    $ go doc github.com/gdamore/chanstream
