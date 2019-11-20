package server

import (
	"net/http"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
)

func StaticFS(group *gin.RouterGroup, relativePath string, dir string) {
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
		panic("URL parameters can not be used when serving a static folder")
	}
	handler := createStaticHandler(group, relativePath, Dir(dir))
	urlPattern := path.Join(relativePath, "/*filepath")

	// Register GET and HEAD handlers
	group.GET(urlPattern, handler)
	group.HEAD(urlPattern, handler)
}

func createStaticHandler(group *gin.RouterGroup, relativePath string, fs FileSystem) gin.HandlerFunc {
	absolutePath := path.Join(group.BasePath(), relativePath)
	fileServer := http.StripPrefix(absolutePath, FileServer(fs))

	return func(c *gin.Context) {
		fileServer.ServeHTTP(c.Writer, c.Request)
	}
}
