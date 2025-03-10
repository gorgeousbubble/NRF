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
	// check request body bind json
	L.Debug("Start bind Register request body to json:", context.Request.Body)
	err := context.ShouldBindJSON(&request)
	if err != nil {
		registrationError := handleRegisterResponseBadRequest(err)
		context.Header("Content-Type", "application/problem+json")
		context.JSON(http.StatusBadRequest, registrationError)
		L.Error("Register request body bind json failed:", err)
		return
	}
	L.Debug("Register request body bind json success.")
	// check request body IEs
	b, err := checkRegisterIEs(&request)
	if b == false {
		registrationError := handleRegisterResponseBadRequest(err)
		context.Header("Content-Type", "application/problem+json")
		context.JSON(http.StatusBadRequest, registrationError)
		L.Error("Register request check failed:", err)
		return
	}
	// handle request body IEs
	err = handleRegisterIEs(&request)
	if err != nil {
		problemDetails := handleRegisterResponseInternalServerError(err)
		context.Header("Content-Type", "application/problem+json")
		context.JSON(http.StatusInternalServerError, problemDetails)
		L.Error("Register request body handle failed:", err)
		return
	}
	// extract nfInstanceId from request uri
	nfInstanceId := context.Param("nfInstanceID")
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

func handleRegisterResponseCreated(request NFProfile) (response NFProfile) {
	// handle Register Response Created (201)
	response = request
	return response
}

func handleRegisterResponseBadRequest(err error) (registrationError NFProfileRegistrationError) {
	// handle Register Response Bad Request (400)
	registrationError.ProblemDetails.Title = "Bad Request"
	registrationError.ProblemDetails.Status = http.StatusBadRequest
	registrationError.ProblemDetails.Detail = err.Error()
	return registrationError
}

func handleRegisterResponseInternalServerError(err error) (problemDetails ProblemDetails) {
	// handle Register Response Internal Server Error (500)
	problemDetails.Title = "Internal Server Error"
	problemDetails.Status = http.StatusInternalServerError
	problemDetails.Detail = err.Error()
	return problemDetails
}

func checkRegisterIEs(request *NFProfile) (b bool, err error) {
	b, err = true, nil
	// check mandatory IEs...
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
	L.Debug("Start CheckNFType:", request.NFType)
	b, err = CheckNFType(request.NFType)
	if err != nil {
		b = false
		L.Error("CheckNFType failed:", err)
		return b, err
	}
	L.Debug("CheckNFType success.")
	// check NFStatus
	L.Debug("Start CheckNFStatus:", request.NFStatus)
	b, err = CheckNFStatus(request.NFStatus)
	if err != nil {
		b = false
		L.Error("CheckNFStatus failed:", err)
		return b, err
	}
	L.Debug("CheckNFStatus success.")
	// check conditional IEs...
	// check HeartBeatTimer
	L.Debug("Start CheckHeartBeatTimer:", request.HeartBeatTimer)
	if request.HeartBeatTimer != 0 {
		b, err = CheckHeartBeatTimer(request.HeartBeatTimer)
		if err != nil {
			b = false
			L.Error("CheckHeartBeatTimer failed:", err)
			return b, err
		}
	}
	L.Debug("CheckHeartBeatTimer success.")
	return b, err
}

func handleRegisterIEs(request *NFProfile) (err error) {
	err = nil
	// handle NFInstanceId
	L.Debug("Start HandleNFInstanceId:", request.NFStatus)
	err = HandleNFInstanceId(&request.NFInstanceId)
	if err != nil {
		L.Error("HandleNFInstanceId failed:", err)
	}
	L.Debug("HandleNFInstanceId success:", request.NFType)
	// handle HeartBeatTimer
	L.Debug("Start HandleHeartBeatTimer:", request.HeartBeatTimer)
	err = HandleHeartBeatTimer(&request.HeartBeatTimer)
	if err != nil {
		L.Error("HandleHeartBeatTimer failed:", err)
	}
	L.Debug("HandleHeartBeatTimer success.")
	return err
}
