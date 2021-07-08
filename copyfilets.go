// copyfilets : recurse dirs to copy mtime from every matching file on matching files in destination dir
// directory structure MUST be the same, use os.Chtimes to set mtime and atime to source mtime
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func copymtime(src, dst string) error {
	fin, err := os.Stat(src)
	if err != nil {
		return err
	}
	err = os.Chtimes(dst, fin.ModTime(), fin.ModTime())
	if err != nil {
		return err
	}
	return nil
}

func recursedir(s, d string) error {
	dirs, err := os.ReadDir(s)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var fi os.FileInfo
	for _, v := range dirs {
		jd := filepath.Join(d, v.Name())
		js := filepath.Join(s, v.Name())
		fdinf, err1 := os.Lstat(jd)
		if err1 != nil {
			continue
		}
		fi, _ = os.Lstat(js)
		if fi.IsDir() {
			if fdinf.IsDir() {
				err := recursedir(js, jd)
				if err != nil {
					return err
				}
			}
		} else {
			if fi.Size() == fdinf.Size() {
				copymtime(js, jd)
			}
		}
	}
	return nil
}
func main() {
	flag.Parse()
	if flag.NArg() < 2 {
		fmt.Println("Usage: copyfilets sourcedir destdir")
		os.Exit(2)
	}
	err := recursedir(flag.Arg(0), flag.Arg(1))
	if err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
