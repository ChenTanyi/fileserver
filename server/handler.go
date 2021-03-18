package server

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
)

func NewFileServer(router *gin.RouterGroup, relativePath, aliasPath string) {
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
		panic("URL parameters can not be used when serving a static folder")
	}

	urlPath := path.Join(router.BasePath(), relativePath)
	urlPattern := path.Join(relativePath, "/*filepath")

	handler := newFileServerHandler(urlPath, Dir(aliasPath))
	router.GET(urlPattern, handler)
	router.HEAD(urlPattern, handler)

	deleteHandler := deleteFileHandler(urlPath, Dir(aliasPath))
	router.DELETE(urlPattern, deleteHandler)

	postHandler := updateFileHandler(urlPath, Dir(aliasPath))
	router.POST(urlPattern, postHandler)
}

func newFileServerHandler(urlPath string, aliasPath Dir) gin.HandlerFunc {
	fileServer := http.StripPrefix(urlPath, FileServer(aliasPath))

	return func(c *gin.Context) {
		defer func() {
			if err, ok := recover().(error); ok && err != nil {
				c.AbortWithStatusJSON(500, err.Error())
			}
		}()

		done, err := processGetQuery(c, aliasPath.Join(c.Param("filepath")))
		if err != nil {
			panic(err)
		}
		if done {
			return
		}

		fileServer.ServeHTTP(c.Writer, c.Request)
	}
}

func deleteFileHandler(urlPath string, aliasPath Dir) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err, ok := recover().(error); ok && err != nil {
				c.AbortWithStatusJSON(500, err.Error())
			}
		}()

		if err := os.RemoveAll(aliasPath.Join(c.Param("filepath"))); err != nil {
			panic(err)
		}
	}
}

func updateFileHandler(urlPath string, aliasPath Dir) gin.HandlerFunc {
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
