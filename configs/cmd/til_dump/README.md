til_dump
========

til_dump is a tool for constructing squares, based on the information retrieved
from a given TIL file, and storing these squares as png images.

Installation
------------

    $ go get github.com/mewkiz/blizzconv/configs/cmd/til_dump
    $ go install github.com/mewkiz/blizzconv/configs/cmd/til_dump

Usage
-----

    $ mkdir blizzdump/
    $ cd blizzdump/
    $ ln -s /path/to/extracted/diabdat_mpq mpqdump
    $ ln -s $GOPATH/src/github.com/mewkiz/blizzconv/images/imgconf/cel.ini
    $ ln -s $GOPATH/src/github.com/mewkiz/blizzconv/mpq/mpq.ini
    $ til_dump town.til
