package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
)

// ServiceInstance Network Function Service Instance
type ServiceInstance struct {
	NFType     string `json:"nfType"`     // Network Function Type (eg. AMF, SMF,...)
	InstanceID string `json:"instanceId"` // Service Instance Identity
	IPAddress  string `json:"IPAddress"`  // Service IP Address
	Port       int    `json:"port"`       // Service Port
	Priority   int    `json:"priority"`   // Priority
	Weight     int    `json:"weight"`     // Weight
}

// NRF Network Repository Function
type NRF struct {
	serviceRegistry map[string][]ServiceInstance
	mutex           sync.RWMutex
}

func NewNRF() *NRF {
	return &NRF{
		serviceRegistry: make(map[string][]ServiceInstance),
	}
}

func (nrf *NRF) Register(ctx *gin.Context) {
	var instance ServiceInstance
	if err := ctx.ShouldBind(&instance); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Mutex Lock Resource
	nrf.mutex.Lock()
	defer nrf.mutex.Unlock()
	// Append New NF Service Instance to Register Map
	nrf.serviceRegistry[instance.NFType] = append(nrf.serviceRegistry[instance.NFType], instance)
	// Return Successful Response
	ctx.JSON(http.StatusCreated, gin.H{"status": "Registered", "instance": instance})
}

func (nrf *NRF) Discovery(ctx *gin.Context) {
	nfType := ctx.Query("nfType")
	if nfType == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'nfType' is required"})
		return
	}
	// Mutex Lock Resource
	nrf.mutex.Lock()
	defer nrf.mutex.Unlock()
	// Query NF Service Instance
	instances, exists := nrf.serviceRegistry[nfType]
	if !exists || len(instances) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "service not found"})
		return
	}
	// Return Successful Response
	ctx.JSON(http.StatusOK, gin.H{"nfType": nfType, "instances": instances})
}

func main() {
	router := gin.Default()
	nrf := NewNRF()
	fmt.Println("The Network Repository Function (NRF) Service.")
	// Register API Route
	api := router.Group("/nrf/v1")
	{
		api.POST("register", nrf.Register)
		api.POST("discovery", nrf.Discovery)
	}
	// Start NRF Service
	router.Run(":8080")
}
