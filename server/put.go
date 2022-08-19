package server

import (
	"io"
	"math"
	"os"

	"github.com/gin-gonic/gin"
)

func uploadFile(c *gin.Context, reqFilePath string) {
	ranges, err := parseRange(c.Request.Header.Get("Range"), math.MaxInt64)
	if err != nil {
		panic(err)
	}
	if len(ranges) != 1 {
		c.String(400, "Unable to write file with ranges %v", ranges)
		return
	}
	ra := ranges[0]

	f, err := os.OpenFile(reqFilePath, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		panic(err)
	}
	size, err := f.Seek(0, io.SeekEnd)
	if err != nil {
		panic(err)
	}
	if size < ra.start {
		c.String(400, "File size %d has gap with range %v", size, ra)
	}
	_, err = f.Seek(ra.start, io.SeekStart)
	if err != nil {
		panic(err)
	}
	_, err = io.CopyN(f, c.Request.Body, ra.length)
	if err != nil {
		panic(err)
	}
	size, err = f.Seek(0, io.SeekEnd)
	if err != nil {
		panic(err)
	}
	c.String(200, "%d", size)
}
