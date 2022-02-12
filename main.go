package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	var files []*GoProVideoFile

	root, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		file, err := ParseGoProVideoFile(path)
		if err != nil {
			return err
		}
		if file != nil {
			files = append(files, file)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	for _, v := range files {
		// fmt.Printf("%+v\n", v)
		fmt.Printf("%s -> %s\n", v.Path, v.NewPath())

		err := os.MkdirAll(filepath.Dir(v.NewPath()), os.ModePerm)
		if err != nil {
			panic(err)
		}
		err = os.Rename(v.Path, v.NewPath())
		if err != nil {
			panic(err)
		}
	}

}
