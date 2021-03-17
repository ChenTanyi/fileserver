package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type PostJSONParams struct {
	Action   string
	Property Property
}

type Property struct {
	Name string
}

type FileProperty struct {
	Name string
	Path string
	Size int64
}

func postWithJSON(c *gin.Context, aliasPath string) {
	body, err := c.GetRawData()
	if err != nil {
		panic(err)
	}
	logrus.Infof("Body: %s", string(body))

	postParams := &PostJSONParams{}
	if err := json.Unmarshal(body, postParams); err != nil {
		panic(err)
	}

	if postParams.Action == "rename" {
		// if strings.Contains(postParams.Property.Name, "/") {
		// 	c.Data(400, "text/plain", []byte("Rename with illegel \"/\""))
		// 	c.Abort()
		// }

		reqFilePath := filepath.Join(aliasPath, c.Param("filepath"))
		resultFilePath := filepath.Join(aliasPath, path.Clean("/"+filepath.Join(c.Param("filepath"), "..", postParams.Property.Name)))
		if err := os.Rename(reqFilePath, resultFilePath); err != nil {
			panic(err)
		}
	} else if postParams.Action == "mkdir" {
		reqDir := filepath.Join(aliasPath, c.Param("filepath"))
		if stat, err := os.Stat(reqDir); err != nil {
			panic(err)
		} else {
			if !stat.IsDir() {
				panic(fmt.Errorf("Parent Dir '%s' is not a dir", reqDir))
			}
			err := os.MkdirAll(filepath.Join(reqDir, path.Clean("/"+postParams.Property.Name)), 0755)
			if err != nil {
				panic(err)
			}
		}
	} else {
		c.Data(400, "text/plain", []byte("Unexpected action"))
	}
}

func postWithForm(c *gin.Context, aliasPath string) {
	err := c.Request.ParseForm()
	if err != nil {
		panic(err)
	}

	c.JSON(200, c.Request.Form)
}

func postWithMultiForm(c *gin.Context, aliasPath string) {
	requestDir := filepath.Join(aliasPath, c.Param("filepath"))
	if stat, err := os.Stat(requestDir); err != nil || !stat.IsDir() {
		tmpDir := filepath.Join(requestDir, "..")
		if stat, err = os.Stat(tmpDir); err != nil || !stat.IsDir() {
			panic(fmt.Errorf("Request path is not dir, err = %v, path = %s", err, requestDir))
		} else {
			requestDir = tmpDir
		}
	}

	multiForm, err := c.MultipartForm()
	if err != nil {
		panic(err)
	}

	allError := errors.New("")
	if err := dealWithNginxUpload(multiForm.Value, requestDir); err != nil {
		allError = fmt.Errorf("%v\n%v", allError, err)
	}

	if err := dealWithFileUpload(multiForm.File, requestDir); err != nil {
		allError = fmt.Errorf("%v\n%v", allError, err)
	}

	if allError = TrimEmptyError(allError); allError == nil {
		c.JSON(200, "Success")
	} else {
		panic(allError)
	}
}

func dealWithNginxUpload(MultiFormValue map[string][]string, requestDir string) error {
	allError := errors.New("")
	err := errors.New("")
	files := make(map[string]*FileProperty)

	for key, value := range MultiFormValue {
		keys := strings.Split(key, ".")
		if len(keys) != 2 && len(value) != 1 {
			msg := fmt.Sprintf("multi form error: key = %s, value = %s", key, value)
			allError = fmt.Errorf("%v\n%s", allError, msg)
			continue
		}
		if keys[0] == "" {
			keys[0] = "default"
		}

		if _, ok := files[keys[0]]; !ok {
			files[keys[0]] = &FileProperty{}
		}
		file := files[keys[0]]
		switch keys[1] {
		case "name":
			file.Name = value[0]
		case "path":
			file.Path = value[0]
		case "size":
			file.Size, err = strconv.ParseInt(value[0], 10, 64)
			if err != nil {
				msg := fmt.Sprintf("Parse size error: %v", err)
				allError = fmt.Errorf("%v\n%s", allError, msg)
			}
		}
	}

	for _, file := range files {
		if stat, err := os.Stat(file.Path); err != nil || stat == nil || stat.Size() != file.Size {
			msg := fmt.Sprintf("Stat file error: err = %v, file = %+v, stat = %+v", err, file, stat)
			allError = fmt.Errorf("%v\n%s", allError, msg)
			continue
		}

		targetFile := filepath.Join(requestDir, file.Name)
		logrus.Infof("Rename file: %s -> %s", file.Path, targetFile)
		if err := os.Rename(file.Path, targetFile); err != nil {
			msg := fmt.Sprintf("Rename err: current path = %s, target path = %s, err = %v", file.Path, targetFile, err)
			allError = fmt.Errorf("%v\n%s", allError, msg)
		}
	}

	return TrimEmptyError(allError)
}

func dealWithFileUpload(MultiFormFile map[string][]*multipart.FileHeader, requestDir string) error {
	allError := errors.New("")
	for _, files := range MultiFormFile {
		for _, f := range files {
			logrus.Infof("Receive file: name = %s, size = %d", f.Filename, f.Size)
			file, err := f.Open()
			if err != nil {
				msg := fmt.Sprintf("Open form file error: err = %v", err)
				allError = fmt.Errorf("%v\n%s", allError, msg)
				continue
			}

			defer file.Close()
			filename := filepath.Join(requestDir, f.Filename)
			targetFile, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0644)
			if err != nil {
				msg := fmt.Sprintf("Create file error: filename = %s, err = %v", filename, err)
				allError = fmt.Errorf("%v\n%s", allError, msg)
				continue
			}

			defer targetFile.Close()
			n, err := io.Copy(targetFile, file)
			if err != nil {
				msg := fmt.Sprintf("Write file error: filename = %s, copy size = %d, err = %v", filename, n, err)
				allError = fmt.Errorf("%v\n%s", allError, msg)
			}
		}
	}
	return TrimEmptyError(allError)
}

func TrimEmptyError(err error) error {
	err = errors.New(strings.TrimSpace(err.Error()))
	if err.Error() == "" {
		return nil
	}
	return err
}
