package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
)

func writer2File(header []string, name string, items ...[]string) error {
	buff := bytes.NewBuffer([]byte{})
	for _, h := range header {
		buff.WriteString(h + "\n")
	}
	for _, item := range items {
		for _, im := range item {
			buff.WriteString(im + "\n")
		}
	}
	return ioutil.WriteFile(name, buff.Bytes(), os.ModePerm)
}

func formatIP(v string, d ...[]string) []string {
	a := make([]string, 0)
	for _, s := range d {
		for _, sv := range s {
			a = append(a, fmt.Sprintf(v, sv))
		}
	}
	return a
}