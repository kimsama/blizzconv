blizzconv
=========

This project attempts to provide functionality for converting different
proprietary formats into modern formats with open specifications.

Supported formats
-----------------

* cel
* cl2
* min
* til
* dun

Usage
-----

The following steps can be taken to convert all CEL, CL2, MIN, TIL and DUN files to PNG images.

1. Install Go [from a binary distribution](http://golang.org/doc/install) or [from source](http://golang.org/doc/install/source).
2. [Configure the GOPATH environment variable](http://golang.org/doc/code.html#GOPATH).

		$ mkdir $HOME/go
		$ export GOPATH=$HOME/go
		$ export PATH=$PATH:$GOPATH/bin

3. Extract `DIABDAT.MPQ` using Ladislav Zezula's [MPQ Editor](http://www.zezula.net/en/mpq/download.html).
4. Download and compile the `img_dump`, `min_dump`, `til_dump` and `dun_dump` commands by running:

		$ go get github.com/mewrnd/blizzconv/images/cmd/img_dump
		$ go get github.com/mewrnd/blizzconv/configs/cmd/min_dump
		$ go get github.com/mewrnd/blizzconv/configs/cmd/til_dump
		$ go get github.com/mewrnd/blizzconv/configs/cmd/dun_dump

5. Set up the environment required by `img_dump`, `min_dump`, `til_dump` and `dun_dump`:

		$ mkdir dump
		$ cd dump
		$ ln -s /path/to/extracted/diabdat_mpq/ mpqdump
		$ ln -s $GOPATH/src/github.com/mewrnd/blizzconv/mpq/mpq.ini mpq.ini
		$ ln -s $GOPATH/src/github.com/mewrnd/blizzconv/images/imgconf/cel.ini cel.ini
		$ ln -s $GOPATH/src/github.com/mewrnd/blizzconv/images/imgconf/cl2.ini cl2.ini
		$ ln -s $GOPATH/src/github.com/mewrnd/blizzconv/configs/dunconf/dun.ini dun.ini

6. Convert all CEL images to PNG images. The following command creates 12045 PNG images (57 MB) and takes about 1m20s to complete on my computer.

		$ time img_dump -imgini=cel.ini -a

7. Convert all CEL images to PNG images. The following command creates XXX PNG images (XXX MB) and takes about XXX to complete on my computer.

		$ time img_dump -imgini=cl2.ini -a

8. Convert all MIN files to PNG images. The following command creates 3286 PNG images (19 MB) and takes about 1m to complete on my computer.

		$ time min_dump l1.min l2.min l3.min l4.min town.min

9. Convert all TIL files to PNG images. The following command creates 1001 PNG images (14 MB) and takes about 40s to complete on my computer.

		$ time til_dump l1.til l2.til l3.til l4.til town.til

10. Convert all DUN files to PNG images. The following command creates 45 PNG images (62 MB) and takes about 4m20s to complete on my computer.

		$ time dun_dump -a

Note: Step 7 takes ages to finish, so I'll fill in the missing details (e.g. XXX) once it has completed.

public domain
-------------

Wherever possible the code of this project is hereby released into the
*[public domain][]*.

[public domain]: https://creativecommons.org/publicdomain/zero/1.0/
