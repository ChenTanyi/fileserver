package server

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/chentanyi/go-utils/filehash"
)

var errNotFile = errors.New("Error: request content is not file")

func Hash(dir, hashFunc string) error {
	filename := dir
	if filename[len(filename)-1] == '/' {
		filename = filename[0 : len(filename)-1]
	}
	filename = fmt.Sprintf("%s.%s", filename, hashFunc)
	return HashToFile(dir, filename, hashFunc)
}

func HashToFile(inputDirpath, outputFilename, hashFunc string) (err error) {
	processingFile := outputFilename + ".processing"
	if _, err = os.Stat(processingFile); !os.IsNotExist(err) {
		if err == nil {
			return fmt.Errorf("File: %s is already exist", processingFile)
		}
		return fmt.Errorf("File: %s stat error, err = %v", processingFile, err)
	}

	if processing, err := os.OpenFile(processingFile, os.O_CREATE|os.O_RDWR, 0755); err != nil {
		return err
	} else {
		processing.Close()
	}

	defer func() {
		if removeErr := os.Remove(processingFile); removeErr != nil {
			if err == nil {
				err = removeErr
			} else {
				err = fmt.Errorf("%v\n%v", err, removeErr)
			}
		}
	}()

	outputFile, err := os.OpenFile(outputFilename, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return err
	}
	defer outputFile.Close()

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

		filename, err := filepath.Rel(inputDirpath, path)
		if err != nil {
			return err
		}
		filename = filepath.ToSlash(filepath.Join(folderName, filename))

		digest, err := filehash.HashAllFileWithFuncName(hashFunc, inputFile)
		if err != nil {
			return err
		}

		writeToFile := fmt.Sprintf("%x  %s\n", digest, filename)
		_, err = outputFile.Write([]byte(writeToFile))
		return err
	})
}

func HashOneFileToResponse(file *os.File, rangeHeader, hashFunc string) (string, error) {
	stat, err := file.Stat()
	if err != nil {
		return "", err
	}

	if stat.IsDir() {
		return "", errNotFile
	}

	ranges, err := parseRange(rangeHeader, stat.Size())
	if err != nil {
		return "", err
	}

	if len(ranges) == 0 {
		digest, err := filehash.HashAllFileWithFuncName(hashFunc, file)
		return fmt.Sprintf("%x", digest), err
	}

	allError := errors.New("")
	hashs := make([]string, 0, len(ranges))
	for _, fileRange := range ranges {
		digest, err := filehash.HashFileWithFuncName(hashFunc, file, fileRange.start, fileRange.start+fileRange.length)
		if err != nil {
			allError = fmt.Errorf("%v\n%v", allError, err)
		} else {
			hashs = append(hashs, fmt.Sprintf("%x", digest))
		}
	}
	if allError = TrimEmptyError(allError); allError != nil {
		return "", allError
	}

	return strings.Join(hashs, ","), nil
}
