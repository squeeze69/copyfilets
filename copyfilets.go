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

type destut interface {
	Destfile(finfo os.FileInfo, name, fullname string) ([]string, error)
	Dorecurse(d string) error
}

type research struct {
	dfiles map[int64][]fn
}

func (r research) Destfile(finfo os.FileInfo, name string, fullname string) ([]string, error) {
	ret := make([]string, 0, 1)
	if b, ok := r.dfiles[finfo.Size()]; ok {
		for i := range b {
			if b[i].Name == name {
				ret = append(ret, b[i].Fullname)
			}
		}
	}
	if len(ret) > 0 {
		return ret, nil
	}
	return ret, errors.New("file not found")
}

func NewResearch(d string) (*research, error) {
	df := make(map[int64][]fn)
	err := recursedir(d, df)
	return &research{dfiles: df}, err
}

//Dorecurse : return nil if it's ok to recurse in SOURCE tree
func (r research) Dorecurse(d string) error {
	return nil
}

//name and fullname for unordered search
type fn struct {
	Name     string
	Fullname string
}

//recupera tutti nomi e li inserisce in una map
func recursedir(d string, df map[int64][]fn) error {
	dirs, err := os.ReadDir(d)
	if err != nil {
		return err
	}
	var fd os.FileInfo
	for _, v := range dirs {
		jd := filepath.Join(d, v.Name())
		fd, err = os.Lstat(jd)
		if err != nil {
			return err
		}
		if fd.IsDir() {
			err = recursedir(jd, df)
			if err != nil {
				return err
			}
		} else {
			size := fd.Size()
			if b, ok := df[size]; ok {
				df[size] = append(b, fn{Name: v.Name(), Fullname: jd})
			} else {
				bucket := make([]fn, 0, 1)
				df[size] = append(bucket, fn{Name: v.Name(), Fullname: jd})
			}
		}
	}
	return nil
}

// copymtime : copies the src file mtime on the dst file mtime/atime
func copymtime(src, dst string) error {
	fin, err := os.Lstat(src)
	if err != nil {
		return err
	}
	err = os.Chtimes(dst, fin.ModTime(), fin.ModTime())
	if err != nil {
		return err
	}
	return nil
}

type classic struct {
}

func (r classic) Dorecurse(d string) error {
	fdinfo, err := os.Lstat(d)
	if err != nil {
		return err
	}
	if fdinfo.IsDir() {
		return nil
	}
	return errors.New("not a dir")
}

func (r classic) Destfile(finfo os.FileInfo, name string, fullname string) ([]string, error) {
	fdinfo, err1 := os.Lstat(fullname)
	if err1 != nil {
		return nil, err1
	}
	if !fdinfo.IsDir() {
		if fdinfo.Size() == finfo.Size() {
			return []string{fullname}, nil
		}
	}
	return nil, errors.New("file not found")
}

func recurseandcopyts(s, d string, dstu destut) error {
	dirs, err := os.ReadDir(s)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var fi os.FileInfo
	var err1 error
	var fullname []string

	for _, v := range dirs {
		jd := filepath.Join(d, v.Name())
		js := filepath.Join(s, v.Name())
		fi, _ = os.Lstat(js)
		if fi.IsDir() {
			if dstu.Dorecurse(jd) == nil {
				err := recurseandcopyts(js, jd, dstu)
				if err != nil {
					return err
				}
			}
		} else {
			fullname, err1 = dstu.Destfile(fi, fi.Name(), jd)
			if err1 != nil {
				continue
			}
			for i := range fullname {
				copymtime(js, fullname[i])
			}
		}
	}
	return nil
}

var unordered = false

var du destut

func main() {
	flag.BoolVar(&unordered, "u", false, "find files in the destination dir")
	flag.Parse()
	if flag.NArg() < 2 {
		fmt.Println("Usage: copyfilets [-u] sourcedir destdir")
		flag.Usage()
		os.Exit(2)
	}
	if unordered {
		du, _ = NewResearch(flag.Arg(1))
		if du == nil {
			os.Exit(1)
		}
	} else {
		du = new(classic)
	}
	err := recurseandcopyts(flag.Arg(0), flag.Arg(1), du)
	if err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
