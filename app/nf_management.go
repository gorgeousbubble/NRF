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
	RequesterFeatures []string `form:"requester-features" binding:"dive,ipv4,ipv6,tls,http2,service-auth"`
}

func HandleNFRegister(context *gin.Context) {
	var request NFProfile
	// record context in logs
	L.Info("NFRegister request:", context.Request)
	// check request body bind json
	L.Debug("Start bind NFRegister request body to json:", context.Request.Body)
	err := context.ShouldBindJSON(&request)
	if err != nil {
		var registrationError NFProfileRegistrationError
		registrationError.ProblemDetails.Title = "Bad Request"
		registrationError.ProblemDetails.Status = http.StatusBadRequest
		registrationError.ProblemDetails.Detail = err.Error()
		context.Header("Content-Type", "application/problem+json")
		context.JSON(http.StatusBadRequest, registrationError)
		L.Error("NFRegister request body bind json failed:", err)
		return
	}
	L.Debug("NFRegister request body bind json success.")
	// check request body IEs
	b, err := checkNFRegisterIEs(&request)
	if b == false && err != nil {
		var registrationError NFProfileRegistrationError
		registrationError.ProblemDetails.Title = "Bad Request"
		registrationError.ProblemDetails.Status = http.StatusBadRequest
		registrationError.ProblemDetails.Detail = err.Error()
		context.Header("Content-Type", "application/problem+json")
		context.JSON(http.StatusBadRequest, registrationError)
		L.Error("NFRegister request check failed:", err)
		return
	}
	// handle request body IEs
	response := request
	err = handleNFRegisterIEs(&response)
	if err != nil {
		var problemDetails ProblemDetails
		problemDetails.Title = "Internal Server Error"
		problemDetails.Status = http.StatusInternalServerError
		problemDetails.Detail = err.Error()
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
		NFInstanceId:   nfInstanceId,
		NFType:         response.NFType,
		NFStatus:       response.NFStatus,
		HeartBeatTimer: response.HeartBeatTimer,
		NFServices:     response.NFServices,
	}
	// store instance in NRF Service database
	func() {
		NRFService.mutex.Lock()
		defer NRFService.mutex.Unlock()
		NRFService.instances[response.NFType] = append(NRFService.instances[response.NFType], instance)
	}()
	// return success response
	context.Header("Content-Type", "application/json")
	context.Header("Location", "http://localhost:8000/nnrf-nfm/v1/nf-instances/"+nfInstanceId)
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
		var problemDetails ProblemDetails
		problemDetails.Title = "Bad Request"
		problemDetails.Status = http.StatusBadRequest
		problemDetails.Detail = err.Error()
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
		var problemDetails ProblemDetails
		problemDetails.Title = "Not Found"
		problemDetails.Status = http.StatusNotFound
		problemDetails.Detail = errors.New("NFInstanceId not found").Error()
		context.Header("Content-Type", "application/problem+json")
		context.JSON(http.StatusNotFound, problemDetails)
		L.Error("NFProfileRetrieve request NFInstance not found:", err)
		return
	}
	// check match request features...
	var supported []string
	for _, v := range response.NFServices {
		supported = append(supported, v.SupportedFeatures)
	}
	if !matchFeatures(request.RequesterFeatures, supported) {
		var problemDetails ProblemDetails
		problemDetails.Title = "Forbidden"
		problemDetails.Status = http.StatusForbidden
		problemDetails.Detail = errors.New("request Features not supported").Error()
		context.Header("Content-Type", "application/problem+json")
		context.JSON(http.StatusForbidden, problemDetails)
		L.Error("NFProfileRetrieve request features not supported:", err)
		return
	}
	// return success response
	context.Header("Content-Type", "application/json")
	context.Header("Cache-Control", "no-cache")
	context.JSON(http.StatusOK, response)
	return
}
