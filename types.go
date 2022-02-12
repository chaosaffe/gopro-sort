package main

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
)

type FileType int

const (
	Single FileType = iota
	Chaptered
	Looping
)

type VideoEncoding rune

const (
	AVC    VideoEncoding = 'H'
	HEVC   VideoEncoding = 'X'
	LOWRES VideoEncoding = 'L'
)

type GoProVideoFile struct {
	Path          string
	Type          FileType
	Encoding      VideoEncoding
	ChapterNumber int
	FileNumber    int
	Extension     string
}

func (g GoProVideoFile) NewPath() string {
	return fmt.Sprintf("%s/%04d/G%s%04d%02d.%s", filepath.Dir(g.Path), g.FileNumber, string(g.Encoding), g.FileNumber, g.ChapterNumber, g.Extension)
}

const VideoRegex = `(?i)G(?P<Encoding>[HXL])(?P<ChapterNumber>[0-9]{2})(?P<FileNumber>[0-9]{4})\.(?P<Extension>mp4|wav|thm|lrv)`

var PathRegex = `.*/(?:[0-9]{4})/` + VideoRegex

func ParseGoProVideoFile(path string) (*GoProVideoFile, error) {

	vr := *regexp.MustCompile(VideoRegex)
	pr := *regexp.MustCompile(PathRegex)

	filename := filepath.Base(path)

	if !vr.MatchString(filename) || pr.MatchString(path) {
		return nil, nil
	}

	out := GoProVideoFile{
		Path: path,
	}
	match := vr.FindStringSubmatch(filename)
	for i, name := range vr.SubexpNames() {
		switch name {
		case "Encoding":
			switch match[i] {
			case "H":
				out.Encoding = AVC
			case "X":
				out.Encoding = HEVC
			case "L":
				out.Encoding = LOWRES
			}
		case "ChapterNumber":
			cn, err := strconv.Atoi(match[i])
			if err != nil {
				return nil, err
			}
			out.ChapterNumber = cn
		case "FileNumber":
			fn, err := strconv.Atoi(match[i])
			if err != nil {
				return nil, err
			}
			out.FileNumber = fn
		case "Extension":
			out.Extension = match[i]
		}

	}
	return &out, nil

}
