package app

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	. "nrf/data"
	. "nrf/logs"
	"strings"
)

type RetrieveRequest struct {
	RequesterFeatures []string `form:"requester-features" binding:"omitempty,dive,oneof=ipv4 ipv6 tls http2 service-auth"`
	TargetNFType      string   `form:"target-nf-type"`
}

func HandleNFRegister(context *gin.Context) {
	var request NFProfile
	// record context in logs
	L.Info("NFRegister request:", context.Request)
	// check request body bind json
	L.Debug("Start bind NFRegister request body to json:", context.Request.Body)
	err := context.ShouldBindJSON(&request)
	if err != nil {
		registrationError := handleRegisterResponseBadRequest(err)
		context.Header("Content-Type", "application/problem+json")
		context.JSON(http.StatusBadRequest, registrationError)
		L.Error("NFRegister request body bind json failed:", err)
		return
	}
	L.Debug("NFRegister request body bind json success.")
	// check request body IEs
	b, err := checkRegisterIEs(&request)
	if b == false {
		registrationError := handleRegisterResponseBadRequest(err)
		context.Header("Content-Type", "application/problem+json")
		context.JSON(http.StatusBadRequest, registrationError)
		L.Error("NFRegister request check failed:", err)
		return
	}
	// handle request body IEs
	err = handleRegisterIEs(&request)
	if err != nil {
		problemDetails := handleRegisterResponseInternalServerError(err)
		context.Header("Content-Type", "application/problem+json")
		context.JSON(http.StatusInternalServerError, problemDetails)
		L.Error("NFRegister request body handle failed:", err)
		return
	}
	// extract nfInstanceId from request uri
	nfInstanceId := strings.ToLower(context.Param("nfInstanceID"))
	fmt.Println("nfInstanceId:", nfInstanceId)
	// create instance from request body
	instance := NFInstance{
		NFInstanceId:   request.NFInstanceId,
		NFType:         request.NFType,
		NFStatus:       request.NFStatus,
		HeartBeatTimer: request.HeartBeatTimer,
	}
	// store instance in NRF Service database
	func() {
		NRFService.mutex.Lock()
		defer NRFService.mutex.Unlock()
		NRFService.instances[request.NFType] = append(NRFService.instances[request.NFType], instance)
	}()
	// return success response
	response := handleRegisterResponseCreated(request)
	context.Header("Content-Type", "application/json")
	context.Header("Location", "http://localhost:8000/nnrf-nfm/v1/nf-instances/"+request.NFInstanceId)
	context.JSON(http.StatusCreated, response)
	return
}

func HandleNFProfileRetrieve(context *gin.Context) {
	var request RetrieveRequest
	// record context in logs
	L.Info("NFProfileRetrieve request:", context.Request)
	// check request body bind json
	L.Debug("Start bind NFProfileRetrieve request body to json:", context.Request.Body)
	err := context.ShouldBindJSON(&request)
	if err != nil {
		problemDetails := handleNFProfileRetrieveResponseBadRequest(err)
		context.Header("Content-Type", "application/problem+json")
		context.JSON(http.StatusBadRequest, problemDetails)
		L.Error("NFProfileRetrieve request body bind json failed:", err)
		return
	}
	L.Debug("NFProfileRetrieve request body bind json success.")
	// extract nfInstanceId from request uri
	nfInstanceId := strings.ToLower(context.Param("nfInstanceID"))
	fmt.Println("nfInstanceId:", nfInstanceId)
	// store instance in NRF Service database
	var response NFInstance
	exists := func(ins *NFInstance) bool {
		NRFService.mutex.RLock()
		defer NRFService.mutex.RUnlock()
		for _, instances := range NRFService.instances {
			for _, v := range instances {
				if v.NFInstanceId == nfInstanceId {
					*ins = v
					return true
				}
			}
		}
		return false
	}(&response)
	if !exists {
		problemDetails := handleNFProfileRetrieveResponseNotFound(err)
		context.Header("Content-Type", "application/problem+json")
		context.JSON(http.StatusNotFound, problemDetails)
		L.Error("NFProfileRetrieve request NFInstance not found:", err)
		return
	}
	// check validate request features
	// check NFType match
	if request.TargetNFType != "" && request.TargetNFType != response.NFType {
		err = errors.New("NFProfileRetrieve request NFType does not match actual NFType")
		problemDetails := handleNFProfileRetrieveResponseBadRequest(err)
		context.Header("Content-Type", "application/problem+json")
		context.JSON(http.StatusBadRequest, problemDetails)
		L.Error("NFProfileRetrieve request body bind json failed:", err)
		return
	}
	// return success response
	context.Header("Content-Type", "application/json")
	context.Header("Cache-Control", "no-cache")
	context.JSON(http.StatusOK, response)
	return
}
