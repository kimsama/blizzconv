til_dump
========

til_dump is a tool for constructing squares, based on the information retrieved
from a given TIL file, and storing these squares as PNG images.

Installation
------------

	$ go get github.com/mewrnd/blizzconv/configs/cmd/til_dump

Usage
-----

	$ mkdir blizzdump/
	$ cd blizzdump/
	$ ln -s /path/to/extracted/diabdat_mpq/ mpqdump
	$ ln -s $GOPATH/src/github.com/mewrnd/blizzconv/mpq/mpq.ini
	$ ln -s $GOPATH/src/github.com/mewrnd/blizzconv/images/imgconf/cel.ini
	$ til_dump l1.til l2.til l3.til l4.til town.til
