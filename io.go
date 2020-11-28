package main

import (
	"archive/zip"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func unzip() error {
	r, err := zip.OpenReader(filepath.Join(pwd, ipSrcFileName))
	if err != nil {
		return err
	}
	for _, z := range r.Reader.File {
		r, err := z.Open()
		if err != nil {
			return err
		}
		if z.FileInfo().IsDir() {
			_ = os.MkdirAll(z.Name, os.ModePerm)
			continue
		}
		NewFile, err := os.Create(filepath.Join(pwd, z.Name))
		if err != nil {
			return err
		}
		_, _ = io.Copy(NewFile, r)
		_ = NewFile.Close()
		_ = r.Close()
	}
	return r.Close()
}

func getCountryList() ([]string, error) {
	fis, err := ioutil.ReadDir(filepath.Join(pwd, ipData))
	if err != nil {
		return nil, err
	}
	cys := make([]string, len(fis))
	for i, f := range fis {
		cys[i] = strings.TrimSuffix(f.Name(), filepath.Ext(f.Name()))
	}
	sort.Strings(cys)
	return cys, nil
}