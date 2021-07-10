# COPYFILETS - Squeeze69

## program to copy files' mtime on mtime/atime written in [GO](https://golang.org)

### License: GPLv3

**BTW**: Credits to "Squeeze69" would be appreciated if you use this code.

**NOTICE**: USE THIS CODE AT YOUR OWN RISK, NO WARRANTIES!

This could cause the end of the universe, or, worst, some bureaucratic nightmares (just kidding).

To get, this should work: go get github.com/squeeze69/copyfilets

Build: go build copyfilets.go

It's a simple tool to copy file mtime (modification time) to other file with the same name and same size, recusing subdirs.

WHY? Because I forgot to use --preserve=timestamps while using cp... with a HUGE number of files, it was faster to write this utility than to re-cp everything.

Usage: copyfilets [-u] sourcedir destdir

The "-u" unordered flag makes the program search in the destdir, without the necessity to match the source directory tree. If multiple files with the same name and same size are found, te mtime & ctime are set as for each one.
