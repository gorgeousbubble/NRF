package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	. "nrf/data"
	. "nrf/logs"
	. "nrf/util"
)

func HandleRegister(context *gin.Context) {
	var request NFProfile
	// record context in logs
	L.Info("Register request:", context.Request)
	// request body bind json
	L.Debug("Start bind Register request body to json:", context.Request.Body)
	err := context.ShouldBindJSON(&request)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		L.Error("Register request body bind json failed:", err)
		return
	}
	L.Debug("Register request body bind json success.")
	// check request body mandatory IEs
	b, err := checkRegisterMandatoryIEs(&request)
	if b == false {
		context.JSON(http.StatusBadRequest, gin.H{"error": "register request lack of mandatory IEs."})
		L.Error("Register request lack of mandatory IEs:", err)
		return
	}
	// extract nfInstanceId from request uri
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

func checkRegisterMandatoryIEs(request *NFProfile) (b bool, err error) {
	b, err = true, nil
	// check NFInstanceId
	L.Debug("Start CheckNFInstanceId:", request.NFInstanceId)
	b, err = CheckNFInstanceId(request.NFInstanceId)
	if err != nil {
		b = false
		L.Error("CheckNFInstanceId failed:", err)
		return b, err
	}
	L.Debug("CheckNFInstanceId success.")
	// check NFType
	L.Debug("Start checkNFType:", request.NFType)
	b, err = CheckNFType(request.NFType)
	if err != nil {
		b = false
		L.Error("CheckNFType failed:", err)
		return b, err
	}
	L.Debug("CheckNFType success.")
	// check NFStatus
	L.Debug("Start checkNFStatus:", request.NFStatus)
	b, err = CheckNFStatus(request.NFStatus)
	if err != nil {
		b = false
		L.Error("CheckNFStatus failed:", err)
		return b, err
	}
	L.Debug("CheckNFStatus success.")
	return b, err
}
