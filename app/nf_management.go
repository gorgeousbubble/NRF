package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	. "nrf/data"
)

func HandleRegister(context *gin.Context) {
	var request NFProfile
	// request body bind json
	err := context.ShouldBindJSON(&request)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	nfInstanceId := context.Param("nfInstanceID")
	fmt.Println("nfInstanceId:", nfInstanceId)
	// create instance from request body
	instance := NFInstance{
		NFInstanceId: request.NFInstanceId,
		NFType:       request.NFType,
		NFStatus:     request.NFStatus,
	}
	// store instance in NRF Service database
	func() {
		NRFService.mutex.Lock()
		defer NRFService.mutex.Unlock()
		NRFService.instances[request.NFType] = append(NRFService.instances[request.NFType], instance)
	}()
	// return success response
	context.JSON(http.StatusCreated, instance)
	return
}
