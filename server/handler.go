package server

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/chentanyi/go-utils/filehash"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func NewFileServer(router *gin.RouterGroup, relativePath, aliasPath string) {
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
		panic("URL parameters can not be used when serving a static folder")
	}

	urlPath := path.Join(router.BasePath(), relativePath)
	urlPattern := path.Join(relativePath, "/*filepath")

	handler := newFileServerHandler(urlPath, aliasPath)
	router.GET(urlPattern, handler)
	router.HEAD(urlPattern, handler)

	deleteHandler := deleteFileHandler(urlPath, aliasPath)
	router.DELETE(urlPattern, deleteHandler)

	postHandler := updateFileHandler(urlPath, aliasPath)
	router.POST(urlPattern, postHandler)
}

func newFileServerHandler(urlPath string, aliasPath string) gin.HandlerFunc {
	fileServer := http.StripPrefix(urlPath, FileServer(Dir(aliasPath)))

	return func(c *gin.Context) {
		defer func() {
			if err, ok := recover().(error); ok && err != nil {
				c.AbortWithStatusJSON(500, err.Error())
			}
		}()

		if c.Query("download") == "tar" {
			reqFilePath := filepath.Join(aliasPath, c.Param("filepath"))
			if _, err := os.Stat(reqFilePath); err != nil {
				panic(err)
			}

			go func() {
				if err := Compress(reqFilePath); err != nil {
					logrus.Errorf("Compress Error: %v", err)
				}
			}()

			query := c.Request.URL.Query()
			query.Del("download")
			c.Redirect(302, fmt.Sprintf("../?%s", query.Encode()))
			return
		}

		if c.Query("hash") != "" {
			reqFilePath := filepath.Join(aliasPath, c.Param("filepath"))
			file, err := os.Open(reqFilePath)
			if err != nil {
				panic(err)
			}
			defer file.Close()

			stat, err := file.Stat()
			if err != nil {
				panic(err)
			}

			if !stat.IsDir() {
				rangeHeader := c.GetHeader("Range")
				ranges, err := parseRange(rangeHeader, stat.Size())
				if err != nil {
					panic(err)
				}

				if len(ranges) == 0 {
					digest, err := filehash.HashFileWithFuncName(c.Query("hash"), file, 0, math.MaxInt64)
					if err != nil {
						panic(err)
					}
					c.String(200, fmt.Sprintf("%x", digest))
					return
				}

				allError := errors.New("")
				hashs := make([]string, 0, len(ranges))
				for _, fileRange := range ranges {
					digest, err := filehash.HashFileWithFuncName(c.Query("hash"), file, fileRange.start, fileRange.start+fileRange.length)
					if err != nil {
						allError = fmt.Errorf("%v\n%v", allError, err)
					} else {
						hashs = append(hashs, fmt.Sprintf("%x", digest))
					}
				}
				if allError = TrimEmptyError(allError); allError != nil {
					panic(allError)
				}

				c.String(200, strings.Join(hashs, ","))
				return
			}
		}

		fileServer.ServeHTTP(c.Writer, c.Request)
	}
}

func deleteFileHandler(urlPath string, aliasPath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err, ok := recover().(error); ok && err != nil {
				c.AbortWithStatusJSON(500, err.Error())
			}
		}()

		reqFilePath := filepath.Join(aliasPath, c.Param("filepath"))
		if err := os.RemoveAll(reqFilePath); err != nil {
			panic(err)
		}
	}
}

func updateFileHandler(urlPath string, aliasPath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err, ok := recover().(error); ok && err != nil {
				c.AbortWithStatusJSON(500, err.Error())
			}
		}()

		// body, err := c.GetRawData()
		// if err != nil {
		// 	panic(err)
		// }

		// logrus.Infof("Body: %s", string(body))

		switch c.ContentType() {
		case "application/json":
			postWithJSON(c, aliasPath)
		case "application/x-www-form-urlencoded":
			postWithForm(c, aliasPath)
		case "multipart/form-data":
			postWithMultiForm(c, aliasPath)
		default:
			c.Data(400, "text/plain", []byte(fmt.Sprintf("Unexpected content-type: %s", c.ContentType())))
		}
	}
}
