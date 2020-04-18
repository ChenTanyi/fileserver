package server

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func processGetQuery(c *gin.Context, reqFilePath string) (done bool, err error) {
	if c.Request.Method != "GET" {
		return false, nil
	}

	file, err := os.Open(reqFilePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return false, err
	}

	switch c.Query("download") {
	case "tar":
		go func() {
			if err := Compress(reqFilePath); err != nil {
				logrus.Errorf("Compress Error: %v", err)
			}
		}()
	case "hash":
		hashFunc := c.Query("hash")
		if hashFunc == "" {
			hashFunc = "sha256"
		}
		go func() {
			if err := Hash(reqFilePath, hashFunc); err != nil {
				logrus.Errorf("Hash Error: %v", err)
			}
		}()
	}

	if c.Query("download") != "" {
		if stat.IsDir() {
			c.Redirect(302, "../")
		} else {
			c.Redirect(302, "./")
		}
		return true, nil
	}

	if c.Query("hash") != "" {
		response, err := HashOneFileToResponse(file, c.GetHeader("Range"), c.Query("hash"))
		if err != nil {
			return false, err
		}
		c.String(200, response)
		return true, nil
	}

	return false, nil
}
