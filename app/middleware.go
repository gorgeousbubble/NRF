package app

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

func AuthorizationMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		// get auth token header
		authHeader := context.GetHeader("Authorization")
		if authHeader == "" {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "authorization_header_missing",
			})
			return
		}
		// extract bearer token
		tokenString := extractBearerToken(authHeader)
		if tokenString == "" {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid_authorization_header",
			})
			return
		}
		// parse and verify token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return oauthConfig.PublicKey, nil
		})
		if err != nil || !token.Valid {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid_token",
			})
			return
		}
		// set token context information
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			context.Set("clientID", claims["sub"])
		}
		context.Next()
	}
}

func ContentEncodingMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		contentEncoding := strings.ToLower(context.GetHeader("Content-Encoding"))
		// skip if not fetch content-encoding
		if contentEncoding == "" {
			context.Next()
			return
		}
		// read raw request body
		rawBody, err := io.ReadAll(context.Request.Body)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
			return
		}
		// choose decompress algorithm according encoding type
		var decReader io.Reader
		switch contentEncoding {
		case "gzip":
			gzipReader, err := gzip.NewReader(bytes.NewBuffer(rawBody))
			if err != nil {
				context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid Gzip format"})
				return
			}
			defer func(gzipReader *gzip.Reader) {
				err = gzipReader.Close()
			}(gzipReader)
			decReader = gzipReader
		case "deflate":
			deflateReader, err := zlib.NewReader(bytes.NewBuffer(rawBody))
			if err != nil {
				context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid Zlib format"})
				return
			}
			defer func(deflateReader io.ReadCloser) {
				err = deflateReader.Close()
			}(deflateReader)
			decReader = deflateReader
		default:
			context.AbortWithStatusJSON(http.StatusUnsupportedMediaType, gin.H{"error": "Unsupported Content-Encoding"})
			return
		}
		// decompress request body
		decBody, err := io.ReadAll(decReader)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Failed to decompress request body"})
			return
		}
		context.Request.Body = io.NopCloser(bytes.NewBuffer(decBody))
		context.Next()
	}
}

func AcceptEncodingMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		acceptEncoding := strings.ToLower(context.GetHeader("Accept-Encoding"))
		// skip if not fetch accept-encoding
		if acceptEncoding == "" {
			context.Next()
			return
		}
		// choose compress algorithm according encoding type
		switch acceptEncoding {
		case "gzip":
			context.Writer = &GzipResponseWriter{context.Writer, gzip.NewWriter(context.Writer)}
		case "deflate":
			context.Writer = &DeflateResponseWriter{context.Writer, zlib.NewWriter(context.Writer)}
		default:
			context.AbortWithStatusJSON(http.StatusUnsupportedMediaType, gin.H{"error": "Unsupported Content-Encoding"})
			return
		}
		context.Header("Content-Encoding", acceptEncoding)
		context.Next()
	}
}

func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Header("X-NRF-API-Version", "1.3.0")
		context.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		context.Next()
	}
}

func ETagMiddleware(config ETagConfig) gin.HandlerFunc {
	return func(context *gin.Context) {
		// skip if request method is not GET
		if context.Request.Method != "GET" {
			context.Next()
			return
		}
		// pre-check client ETag
		clientETag := context.GetHeader("If-None-Match")
		context.Next()
		// get response content
		respBody, exists := context.Get("responseBody")
		if !exists || respBody == nil {
			return
		}
		// generate ETag
		content := respBody.([]byte)
		etag := generateETag(content, config.WeakValidation)
		// set response header
		context.Header("ETag", etag)
		context.Header("Cache-Control", fmt.Sprintf("public, max-age=%d", config.CacheMaxAge))
		// verify client ETag
		if clientETag != "" {
			if compareETags(clientETag, etag, config.WeakValidation) {
				context.AbortWithStatus(http.StatusNotModified)
				return
			}
		}
		// send response
		context.Data(http.StatusOK, context.Writer.Header().Get("Content-Type"), content)
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

func extractBearerToken(header string) string {
	if len(header) > 7 && header[:7] == "Bearer " {
		return header[7:]
	}
	return ""
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
