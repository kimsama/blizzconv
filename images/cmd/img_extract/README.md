img_extract
===========

img_extract is a tool for extracting CEL and CL2 archives.

Installation
------------

    $ go get github.com/mewkiz/blizzconv/images/cmd/img_extract
    $ go install github.com/mewkiz/blizzconv/images/cmd/img_extract

Usage
-----

    $ mkdir blizzdump/
    $ cd blizzdump/
    $ ln -s /path/to/extracted/diabdat_mpq mpqdump
    $ ln -s $GOPATH/src/github.com/mewkiz/blizzconv/images/imgconf/cel.ini
    $ ln -s $GOPATH/src/github.com/mewkiz/blizzconv/images/imgconf/cl2.ini
    $ ln -s $GOPATH/src/github.com/mewkiz/blizzconv/mpq/mpq.ini
    $ img_extract cow.cel
