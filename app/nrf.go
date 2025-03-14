package app

import (
	"github.com/gin-gonic/gin"
	. "nrf/conf"
	. "nrf/data"
	. "nrf/logs"
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
	router := gin.Default()
	// initialize OAuth2 public key
	oauthConfig.PublicKey = &oauthConfig.PrivateKey.PublicKey
	// middleware handle functions
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(ContentEncodingMiddleware())
	router.Use(AcceptEncodingMiddleware())
	router.Use(SecurityHeadersMiddleware())
	router.Use(ETagMiddleware(defaultConfig))
	// OAuth2 protect
	/*protected := router.Group("/nnrf-nfm/v1")
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
	err := router.Run(":8080")
	if err != nil {
		panic(err)
	}
}

func (nrf *NRF) Stop() {

}

func (nrf *NRF) Status() {

}
