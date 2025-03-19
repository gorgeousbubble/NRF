package app

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	. "nrf/conf"
	. "nrf/data"
	. "nrf/logs"
	"os"
	"sync"
)

var NRFService *NRF

type NRF struct {
	instances map[string][]NFInstance
	mutex     sync.RWMutex
}

type NFInstance struct {
	NFInstanceId   string      `json:"nfInstanceId" yaml:"nfInstanceId"`
	NFType         string      `json:"nfType" yaml:"nfType"`
	NFStatus       string      `json:"nfStatus" yaml:"nfStatus"`
	HeartBeatTimer int         `json:"heartBeatTimer" yaml:"heartBeatTimer"`
	NFServices     []NFService `json:"nfServices" yaml:"nfServices"`
}

func New() *NRF {
	return &NRF{
		instances: make(map[string][]NFInstance),
	}
}

func (nrf *NRF) Init() (err error) {
	err = InitLog()
	if err != nil {
		L.Error("Initialize NRF Logger failed:", err.Error())
		return err
	}
	L.Info("Initialize NRF Logger Success.")
	L.Info("Loading NRF Configuration...")
	err = LoadConf()
	if err != nil {
		L.Error("Loading NRF Configuration failed:", err.Error())
		return err
	}
	L.Info("Loading NRF Configuration Success.")
	L.Info("Initialize NRF Success.")
	return err
}

func (nrf *NRF) Start() {
	// create default Gin Engine instance
	router := gin.Default()
	// enable SBI TLS layer
	var tlsConfig *tls.Config
	tlsSettings := NRFConfigure.SBITLSSettings.TLSType
	if tlsSettings != "non-tls" {
		tlsConfig = &tls.Config{
			MinVersion: tls.VersionTLS13,
			CurvePreferences: []tls.CurveID{
				tls.X25519,
				tls.CurveP256,
			},
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			},
		}
		if tlsSettings == "mutual-tls" {
			// read CA certificate
			caCert, err := os.ReadFile("cert/ca.crt")
			if err != nil {
				log.Fatal(err)
			}
			// create CA certificate pool
			caCertPool := x509.NewCertPool()
			caCertPool.AppendCertsFromPEM(caCert)
			// specific CA certificate pool
			tlsConfig.ClientCAs = caCertPool
			tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
		}
	}
	// middleware handle functions
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(ContentEncodingMiddleware())
	router.Use(AcceptEncodingMiddleware())
	router.Use(SecurityHeadersMiddleware())
	router.Use(ETagMiddleware(defaultConfig))
	// OAuth2 protect
	// initialize OAuth2 public key
	/*oauthConfig.PublicKey = &oauthConfig.PrivateKey.PublicKey
	protected := router.Group("/nnrf-nfm/v1")
	protected.Use(AuthorizationMiddleware())
	{
		protected.PUT("nf-instances/:nfInstanceID", HandleNFRegisterOrNFProfileCompleteReplacement)
		protected.GET("nf-instances/:nfInstanceID", HandleNFProfileRetrieve)
	}*/
	// API route groups
	nfManagement := router.Group("/nnrf-nfm/v1")
	{
		nfManagement.PUT("nf-instances/:nfInstanceID", HandleNFRegisterOrNFProfileCompleteReplacement)
		nfManagement.GET("nf-instances/:nfInstanceID", HandleNFProfileRetrieve)
	}
	// start NRF services
	if tlsSettings == "non-tls" {
		// start http service
		server := &http.Server{
			Addr:    ":8080",
			Handler: router,
		}
		// listen and serve on http port
		fmt.Println("The NRF start http server on", server.Addr)
		L.Info("The NRF start http server on", server.Addr)
		err := server.ListenAndServe()
		if err != nil {
			fmt.Println("The NRF start http server failed:", err.Error())
			L.Error("The NRF start http server failed:", err.Error())
			os.Exit(1)
		}
	} else {
		// start https service
		server := &http.Server{
			Addr:      ":8443",
			Handler:   router,
			TLSConfig: tlsConfig,
		}
		// listen and serve on https port
		fmt.Println("The NRF start https server on", server.Addr)
		L.Info("The NRF start https server on", server.Addr)
		err := server.ListenAndServeTLS("cert/nrf.pem", "cert/nrf.key")
		if err != nil {
			fmt.Println("The NRF start https server failed:", err.Error())
			L.Error("The NRF start https server failed:", err.Error())
			os.Exit(1)
		}
	}
}

func (nrf *NRF) Stop() {

}

func (nrf *NRF) Status() {

}
