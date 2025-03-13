package app

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strings"
)

type ETagConfig struct {
	WeakValidation bool // 是否使用弱验证 (W/)
	CacheMaxAge    int  // 缓存最大年龄（秒）
}

var defaultConfig = ETagConfig{
	WeakValidation: false,
	CacheMaxAge:    3600,
}

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

func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-NRF-API-Version", "1.3.0")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Next()
	}
}

func ETagMiddleware(config ETagConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// skip if request method is not GET
		if c.Request.Method != "GET" {
			c.Next()
			return
		}
		// pre-check client ETag
		clientETag := c.GetHeader("If-None-Match")
		c.Next()
		// get response content
		respBody, exists := c.Get("responseBody")
		if !exists || respBody == nil {
			return
		}
		// generate ETag
		content := respBody.([]byte)
		etag := generateETag(content, config.WeakValidation)
		// set response header
		c.Header("ETag", etag)
		c.Header("Cache-Control", fmt.Sprintf("public, max-age=%d", config.CacheMaxAge))
		// verify client ETag
		if clientETag != "" {
			if compareETags(clientETag, etag, config.WeakValidation) {
				c.AbortWithStatus(http.StatusNotModified)
				return
			}
		}
		// send response
		c.Data(http.StatusOK, c.Writer.Header().Get("Content-Type"), content)
	}
}

func generateETag(data []byte, weak bool) string {
	hash := fmt.Sprintf("%x", sha1.Sum(data))
	if weak {
		return fmt.Sprintf("W/\"%s\"", hash)
	}
	return fmt.Sprintf("\"%s\"", hash)
}

func compareETags(clientTag, serverTag string, weakCompare bool) bool {
	// clean ETag
	cleanTag := func(t string) string {
		if len(t) >= 2 && t[0] == '"' {
			t = t[1 : len(t)-1]
		}
		if len(t) > 2 && t[:2] == "W/" {
			t = t[2:]
		}
		return t
	}
	// clean client and server ETag
	client := cleanTag(clientTag)
	server := cleanTag(serverTag)
	// compare ETag
	if weakCompare {
		return client == server
	}
	return client == server && !strings.HasPrefix(clientTag, "W/")
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
