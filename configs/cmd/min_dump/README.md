min_dump
========

min_dump is a tool for constructing pillars, based on the information retrieved
from a given MIN file, and storing these pillars as png images.

Installation
------------

    $ go get github.com/mewrnd/blizzconv/configs/cmd/min_dump
    $ go install github.com/mewrnd/blizzconv/configs/cmd/min_dump

Usage
-----

    $ mkdir blizzdump/
    $ cd blizzdump/
    $ ln -s /path/to/extracted/diabdat_mpq mpqdump
    $ ln -s $GOPATH/src/github.com/mewrnd/blizzconv/images/imgconf/cel.ini
    $ ln -s $GOPATH/src/github.com/mewrnd/blizzconv/mpq/mpq.ini
    $ min_dump town.min
