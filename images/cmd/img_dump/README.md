img_dump
========

img_dump is a tool for converting CEL and CL2 images into png images.

Installation
------------

    $ go get github.com/mewkiz/blizzconv/images/cmd/img_dump
    $ go install github.com/mewkiz/blizzconv/images/cmd/img_dump

Usage
-----

    $ mkdir blizzdump/
    $ cd blizzdump/
    $ ln -s /path/to/extracted/diabdat_mpq mpqdump
    $ ln -s $GOPATH/src/github.com/mewkiz/blizzconv/images/imgconf/cel.ini
    $ ln -s $GOPATH/src/github.com/mewkiz/blizzconv/images/imgconf/cl2.ini
    $ ln -s $GOPATH/src/github.com/mewkiz/blizzconv/mpq/mpq.ini
    $ img_dump -a
