// copyfilets : recurse dirs to copy mtime from every matching file on matching files in destination dir
// directory structure MUST be the same, use os.Chtimes to set mtime and atime to source mtime
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

// copymtime : copies the src file mtime on the dst file mtime/atime
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

func destfile(finfo os.FileInfo, name string, fullname string) (string, error) {
	if unordered {
		if b, ok := dfiles[finfo.Size()]; ok {
			for _, v := range b {
				if v.Name == name {
					return v.Fullname, nil
				}
			}
		}
	} else {
		fdinfo, err1 := os.Lstat(fullname)
		if err1 != nil {
			return "", err1
		}
		if !fdinfo.IsDir() {
			return fullname, nil
		}
	}

	return "", errors.New("file not found")
}

func recurseandcopyts(s, d string) error {
	dirs, err := os.ReadDir(s)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var fi, fdinf os.FileInfo
	var err1 error
	var fullname string

	for _, v := range dirs {
		jd := filepath.Join(d, v.Name())
		js := filepath.Join(s, v.Name())
		if !unordered {
			fdinf, err1 = os.Lstat(jd)
			if err1 != nil {
				continue
			}
		}
		fi, _ = os.Lstat(js)
		if fi.IsDir() {
			if fdinf.IsDir() {
				err := recurseandcopyts(js, jd)
				if err != nil {
					return err
				}
			}
		} else {
			if unordered {
				fullname, err1 = destfile(fi, fi.Name(), jd)
				if err1 != nil {
					continue
				}
				copymtime(js, fullname)
			} else {
				if fi.Size() == fdinf.Size() {
					copymtime(js, jd)
				}
			}
		}
	}
	return nil
}

type fn struct {
	Name     string
	Fullname string
}

//recupera tutti nomi e li inserisce in una map
func recursedir(mp map[int64][]fn, d string) error {
	dirs, err := os.ReadDir(d)
	if err != nil {
		return err
	}
	var fd os.FileInfo
	var err1 error
	for _, v := range dirs {
		jd := filepath.Join(d, v.Name())
		fd, err1 = os.Lstat(jd)
		if err1 != nil {
			return err1
		}
		if fd.IsDir() {
			err = recursedir(mp, jd)
			if err != nil {
				return err
			}
		} else {
			if b, ok := mp[fd.Size()]; ok {
				b = append(b, fn{Name: v.Name(), Fullname: jd})
				mp[fd.Size()] = b
			} else {
				bucket := make([]fn, 1)
				bucket = append(bucket, fn{Name: v.Name(), Fullname: jd})
				mp[fd.Size()] = bucket
			}
		}
	}
	return nil
}

var unordered = false

var dfiles map[int64][]fn

func main() {
	flag.BoolVar(&unordered, "u", false, "find files in the destination dir")
	flag.Parse()
	if flag.NArg() < 2 {
		fmt.Println("Usage: copyfilets [-u] sourcedir destdir")
		flag.Usage()
		os.Exit(2)
	}
	if unordered {
		dfiles = make(map[int64][]fn)
		err := recursedir(dfiles, flag.Arg(1))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	err := recurseandcopyts(flag.Arg(0), flag.Arg(1))
	if err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
