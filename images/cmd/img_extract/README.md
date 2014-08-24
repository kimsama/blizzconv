img_extract
===========

img_extract is a tool for extracting CEL and CL2 archives.

Installation
------------

	$ go get github.com/mewrnd/blizzconv/images/cmd/img_extract

Usage
-----

	$ mkdir blizzdump/
	$ cd blizzdump/
	$ ln -s /path/to/extracted/diabdat_mpq/ mpqdump
	$ ln -s $GOPATH/src/github.com/mewrnd/blizzconv/mpq/mpq.ini
	$ ln -s $GOPATH/src/github.com/mewrnd/blizzconv/images/imgconf/cel.ini
	$ ln -s $GOPATH/src/github.com/mewrnd/blizzconv/images/imgconf/cl2.ini
	$ img_extract cow.cel
