package server

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Compress(dir string) error {
	filename := dir
	if filename[len(filename)-1] == '/' {
		filename = filename[0 : len(filename)-1]
	}
	filename += ".tar.gz"
	return TarGzCompress(dir, filename)
}

func TarGzCompress(inputDirpath string, outputFilename string) (err error) {
	processingFile := outputFilename + ".processing"
	if _, err = os.Stat(processingFile); !os.IsNotExist(err) {
		if err == nil {
			return fmt.Errorf("File: %s is already exist", processingFile)
		}
		return fmt.Errorf("File: %s stat error, err = %v", processingFile, err)
	}

	if _, err = os.OpenFile(processingFile, os.O_CREATE|os.O_RDWR, 0755); err != nil {
		return err
	}

	defer func() {
		err = os.Remove(processingFile)
	}()

	outputFile, err := os.OpenFile(outputFilename, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return err
	}

	defer outputFile.Close()

	gzWriter := gzip.NewWriter(outputFile)
	defer gzWriter.Close()

	tarWriter := tar.NewWriter(gzWriter)
	defer tarWriter.Close()

	folderName := filepath.Base(inputDirpath)

	return filepath.Walk(inputDirpath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		inputFile, err := os.Open(path)
		if err != nil {
			return err
		}

		defer inputFile.Close()

		link := ""
		if info.Mode()&os.ModeSymlink != 0 {
			link, _ = os.Readlink(path)
		}

		tarHeader, err := tar.FileInfoHeader(info, link)
		if err != nil {
			return err
		}

		filename := strings.TrimPrefix(strings.Replace(path, inputDirpath, "", -1), string(filepath.Separator))
		tarHeader.Name = filepath.ToSlash(filepath.Join(folderName, filename))

		return compress(tarWriter, inputFile, tarHeader)
	})
}

func compress(w *tar.Writer, r io.Reader, header *tar.Header) error {
	if err := w.WriteHeader(header); err != nil {
		return err
	}

	_, err := io.Copy(w, r)
	if err != nil {
		return err
	}

	return nil
}
