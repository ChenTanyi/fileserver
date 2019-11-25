package server

import (
	"fmt"
	"net/http"
	"path"
	"path/filepath"
	"strings"

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
}

func newFileServerHandler(urlPath string, aliasPath string) gin.HandlerFunc {
	fileServer := http.StripPrefix(urlPath, FileServer(Dir(aliasPath)))

	return func(c *gin.Context) {
		defer func() {
			if err, ok := recover().(error); ok && err != nil {
				c.Error(err)
			}
		}()

		if strings.Contains(c.Query("download"), "tar") {
			reqFilePath := filepath.Join(aliasPath, c.Param("filepath"))
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

		fileServer.ServeHTTP(c.Writer, c.Request)
	}
}
