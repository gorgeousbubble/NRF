package app

import (
	"crypto/rand"
	"crypto/rsa"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

type OAuthConfig struct {
	ClientID     string
	ClientSecret string
	PrivateKey   *rsa.PrivateKey
	PublicKey    *rsa.PublicKey
}

type ProtectedResource struct {
	NFInstanceID string `json:"nfInstanceId"`
	IPAddress    string `json:"ipAddress"`
}

var oauthConfig = &OAuthConfig{
	ClientID:     "NRF_Service",
	ClientSecret: "123456",
	PrivateKey:   generateRSAKey(),
	PublicKey:    nil,
}

func generateRSAKey() *rsa.PrivateKey {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic("failed to generate RSA key")
	}
	return privateKey
}

func HandleAccessToken(context *gin.Context) {
	clientID := context.PostForm("client_id")
	clientSecret := context.PostForm("client_secret")
	grantType := context.PostForm("grant_type")
	// verify client credentials
	if clientID != oauthConfig.ClientID || clientSecret != oauthConfig.ClientSecret {
		context.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid_client",
		})
		return
	}
	// verify grant_type
	if grantType != "client_credentials" {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "unsupported_grant_type",
		})
		return
	}
	// creat JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss": "nrf-oauth-server",
		"sub": clientID,
		"aud": []string{"nrf-service"},
		"exp": time.Now().Add(time.Hour * 1).Unix(),
		"iat": time.Now().Unix(),
	})
	// convert token to string
	tokenString, err := token.SignedString(oauthConfig.PrivateKey)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed_to_generate_token",
		})
		return
	}
	// return access token
	context.JSON(http.StatusOK, gin.H{
		"access_token": tokenString,
		"token_type":   "Bearer",
		"expires_in":   3600,
	})
	return
}
