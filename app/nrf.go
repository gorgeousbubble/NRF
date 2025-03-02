package app

import (
	"github.com/gin-gonic/gin"
	. "nrf/logs"
	"sync"
)

var NRFService *NRF

type NRF struct {
	instances map[string][]NFInstance
	mutex     sync.Mutex
}

type NFInstance struct {
	NFInstanceId string `json:"nfInstanceId" yaml:"nfInstanceId"`
	NFType       string `json:"nfType" yaml:"nfType"`
	NFStatus     string `json:"nfStatus" yaml:"nfStatus"`
}

func New() *NRF {
	return &NRF{
		instances: make(map[string][]NFInstance),
	}
}

func (nrf *NRF) Init() {
	InitLog()
	L.Info("Initialize NRF Successfully.")
}

func (nrf *NRF) Start() {
	router := gin.Default()
	// middleware handle functions
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	// API route groups
	nfManagement := router.Group("/nnrf-nfm/v1")
	{
		nfManagement.PUT("nf-instances/:nfInstanceID", HandleRegister)
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
