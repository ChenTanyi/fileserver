package server

import (
	"fmt"
	"net/url"
	"os"
)

const timeFormat = "02-Jan-2006 15:04"

type FileInfo struct {
	Info os.FileInfo
}

type DirListHtmlTemplate struct {
	Title string
	Files []*FileInfo
}

func (fileInfo *FileInfo) Name() string {
	name := fileInfo.Info.Name()
	if fileInfo.Info.IsDir() {
		name += "/"
	}
	return name
}

func (fileInfo *FileInfo) HtmlName() string {
	return htmlReplacer.Replace(fileInfo.Name())
}

func (fileInfo *FileInfo) Link() string {
	uri := url.URL{Path: fileInfo.Name()}
	return uri.String()
}

func (fileInfo *FileInfo) ModTime() string {
	return fileInfo.Info.ModTime().Format(timeFormat)
}

func (fileInfo *FileInfo) SizeReadable() string {
	size := fileInfo.Info.Size()
	if size == 0 {
		return "-"
	}

	suffixes := []string{" B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}
	index := 0
	sz := float64(size)
	for sz >= 1024 {
		index++
		sz /= 1024
	}
	return fmt.Sprintf("%.2f%s", sz, suffixes[index])
}
