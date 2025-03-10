package app

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strings"
)

func ContentEncodingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		contentEncoding := strings.ToLower(c.GetHeader("Content-Encoding"))
		// skip if not fetch content-encoding
		if contentEncoding == "" {
			c.Next()
			return
		}
		// read raw request body
		rawBody, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
			return
		}
		// choose decompress algorithm according encoding type
		var decReader io.Reader
		switch contentEncoding {
		case "gzip":
			gzipReader, err := gzip.NewReader(bytes.NewBuffer(rawBody))
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid Gzip format"})
				return
			}
			defer func(gzipReader *gzip.Reader) {
				err = gzipReader.Close()
			}(gzipReader)
			decReader = gzipReader
		case "deflate":
			deflateReader, err := zlib.NewReader(bytes.NewBuffer(rawBody))
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid Zlib format"})
				return
			}
			defer func(deflateReader io.ReadCloser) {
				err = deflateReader.Close()
			}(deflateReader)
			decReader = deflateReader
		default:
			c.AbortWithStatusJSON(http.StatusUnsupportedMediaType, gin.H{"error": "Unsupported Content-Encoding"})
			return
		}
		// decompress request body
		decBody, err := io.ReadAll(decReader)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Failed to decompress request body"})
			return
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(decBody))
		c.Next()
	}
}

func AcceptEncodingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		acceptEncoding := strings.ToLower(c.GetHeader("Accept-Encoding"))
		// skip if not fetch accept-encoding
		if acceptEncoding == "" {
			c.Next()
			return
		}
		// choose compress algorithm according encoding type
		switch acceptEncoding {
		case "gzip":
			c.Writer = &GzipResponseWriter{c.Writer, gzip.NewWriter(c.Writer)}
		case "deflate":
			c.Writer = &DeflateResponseWriter{c.Writer, zlib.NewWriter(c.Writer)}
		default:
			c.AbortWithStatusJSON(http.StatusUnsupportedMediaType, gin.H{"error": "Unsupported Content-Encoding"})
			return
		}
		c.Header("Content-Encoding", acceptEncoding)
		c.Next()
	}
}

type GzipResponseWriter struct {
	gin.ResponseWriter
	compressWriter *gzip.Writer
}

func (w *GzipResponseWriter) Write(data []byte) (int, error) {
	return w.compressWriter.Write(data)
}

func (w *GzipResponseWriter) Close() {
	defer func(compressWriter *gzip.Writer) {
		err := compressWriter.Close()
		if err != nil {
			return
		}
	}(w.compressWriter)
}

type DeflateResponseWriter struct {
	gin.ResponseWriter
	compressWriter *zlib.Writer
}

func (w *DeflateResponseWriter) Write(data []byte) (int, error) {
	return w.compressWriter.Write(data)
}

func (w *DeflateResponseWriter) Close() {
	defer func(compressWriter *zlib.Writer) {
		err := compressWriter.Close()
		if err != nil {
			return
		}
	}(w.compressWriter)
}
