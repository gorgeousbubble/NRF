package app

import (
	"net/http"
	. "nrf/data"
	. "nrf/logs"
	. "nrf/util"
)

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

func handleNFProfileRetrieveResponseBadRequest(err error) (problemDetails ProblemDetails) {
	// handle NFProfileRetrieve Response Bad Request (400)
	problemDetails.Title = "Bad Request"
	problemDetails.Status = http.StatusBadRequest
	problemDetails.Detail = err.Error()
	return problemDetails
}

func handleNFProfileRetrieveResponseNotFound(err error) (problemDetails ProblemDetails) {
	// handle NFProfileRetrieve Response Not Found (404)
	problemDetails.Title = "Not Found"
	problemDetails.Status = http.StatusNotFound
	problemDetails.Detail = err.Error()
	return problemDetails
}

func validateRequesterFeatures(nfFeatures, reqFeatures []string) bool {
	featureSet := make(map[string]struct{})
	for _, f := range nfFeatures {
		featureSet[f] = struct{}{}
	}

	for _, rf := range reqFeatures {
		if _, exists := featureSet[rf]; !exists {
			return false
		}
	}
	return true
}
