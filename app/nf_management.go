package app

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	. "nrf/conf"
	. "nrf/data"
	. "nrf/logs"
	"strings"
)

type NFProfileRetrieveRequest struct {
	RequesterFeatures []string `form:"requester-features" binding:"omitempty,dive,ipv4,ipv6,tls,http2,service-auth"`
}

type NFListRetrieveRequest struct {
	NFType     string `form:"nf-type" binding:"omitempty,oneof=NRF UDM AMF SMF AUSF NEF PCF SMSF NSSF UDR LMF GMLC 5G_EIR SEPP UPF N3IWF AF UDSF BSF CHF NWDAF PCSCF CBCF UCMF HSS SOR_AF SPAF MME SCSAS SCEF SCP NSSAAF ICSCF SCSCF DRA IMS_AS AANF 5G_DDNMF NSACF MFAF EASDF DCCF MB_SMF TSCTSF ADRF GBA_BSF CEF MB_UPF NSWOF PKMF MNPF SMS_GMSC SMS_IWMSC MBSF MBSTF PANF IP_SM_GW SMS_ROUTER DCSF MRF MRFP MF SLPKMF RH"`
	Limit      int    `form:"limit" binding:"omitempty,min=1"`
	PageNumber int    `form:"page-number" binding:"omitempty,min=1"`
	PageSize   int    `form:"page-size" binding:"omitempty,min=1"`
}

func HandleNFRegisterOrNFProfileCompleteReplacement(context *gin.Context) {
	// extract nfInstanceId from request uri
	nfInstanceId := strings.ToLower(context.Param("nfInstanceID"))
	fmt.Println("nfInstanceId:", nfInstanceId)
	// found nfInstanceId in database
	exists := func() bool {
		NRFService.mutex.RLock()
		defer NRFService.mutex.RUnlock()
		for _, instances := range NRFService.instances {
			for _, v := range instances {
				if v.NFInstanceId == nfInstanceId {
					return true
				}
			}
		}
		return false
	}()
	if !exists {
		// NFRegister
		HandleNFRegister(context)
	} else {
		// NFUpdate (Profile Complete Replacement)
		HandleNFProfileCompleteReplacement(context)
	}
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
	context.Header("Location", formLocation(context, "nnrf-nfm", "v1", "nf-instances", nfInstanceId))
	context.JSON(http.StatusCreated, response)
	return
}

func HandleNFProfileCompleteReplacement(context *gin.Context) {
	var request NFProfile
	L.Info("NFProfileCompleteReplacement request:", context.Request)
	// check request body bind json
	L.Debug("Start bind NFProfileCompleteReplacement request body to json:", context.Request.Body)
	err := context.ShouldBindJSON(&request)
	if err != nil {
		var registrationError NFProfileRegistrationError
		registrationError.ProblemDetails.Title = "Bad Request"
		registrationError.ProblemDetails.Status = http.StatusBadRequest
		registrationError.ProblemDetails.Detail = err.Error()
		context.Header("Content-Type", "application/problem+json")
		context.JSON(http.StatusBadRequest, registrationError)
		L.Error("NFProfileCompleteReplacement request body bind json failed:", err)
		return
	}
	L.Debug("NFProfileCompleteReplacement request body bind json success.")
	// check request body IEs
	b, err := checkNFRegisterIEs(&request)
	if b == false && err != nil {
		var registrationError NFProfileRegistrationError
		registrationError.ProblemDetails.Title = "Bad Request"
		registrationError.ProblemDetails.Status = http.StatusBadRequest
		registrationError.ProblemDetails.Detail = err.Error()
		context.Header("Content-Type", "application/problem+json")
		context.JSON(http.StatusBadRequest, registrationError)
		L.Error("NFProfileCompleteReplacement request check failed:", err)
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
		L.Error("NFProfileCompleteReplacement request body handle failed:", err)
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
	err = func(instance *NFInstance) (err error) {
		NRFService.mutex.RLock()
		defer NRFService.mutex.RUnlock()
		for _, instances := range NRFService.instances {
			for k, v := range instances {
				if v.NFInstanceId == nfInstanceId {
					instances[k], err = *instance, nil
					return err
				}
			}
		}
		err = errors.New("NFInstance not found")
		return err
	}(&instance)
	if err != nil {
		var problemDetails ProblemDetails
		problemDetails.Title = "Internal Server Error"
		problemDetails.Status = http.StatusInternalServerError
		problemDetails.Detail = err.Error()
		context.Header("Content-Type", "application/problem+json")
		context.JSON(http.StatusInternalServerError, problemDetails)
		L.Error("NFProfileCompleteReplacement profile complete replacement failed:", err)
		return
	}
	// return success response
	context.Header("Content-Type", "application/json")
	context.JSON(http.StatusOK, response)
	return
}

func HandleNFProfileRetrieve(context *gin.Context) {
	var request NFProfileRetrieveRequest
	var requestFeatureFilter bool
	// record context in logs
	L.Info("NFProfileRetrieve request:", context.Request)
	// check request body bind json
	L.Debug("Start bind NFProfileRetrieve request body to json:", context.Request.Body)
	err := context.ShouldBindJSON(&request)
	if err != nil {
		requestFeatureFilter = false
		L.Debug("NFProfileRetrieve request body bind json failed:", err.Error())
		L.Debug("NFProfileRetrieve request-feature filter not allowed.")
	} else {
		requestFeatureFilter = true
		L.Debug("NFProfileRetrieve request body bind json success.")
		L.Debug("NFProfileRetrieve request-feature filter allowed.")
	}
	// extract nfInstanceId from request uri
	nfInstanceId := strings.ToLower(context.Param("nfInstanceID"))
	fmt.Println("nfInstanceId:", nfInstanceId)
	// found instance in NRF Service database
	var response NFInstance
	exists := func(instance *NFInstance) bool {
		NRFService.mutex.RLock()
		defer NRFService.mutex.RUnlock()
		for _, instances := range NRFService.instances {
			for _, v := range instances {
				if v.NFInstanceId == nfInstanceId {
					*instance = v
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
	// check match request features (request-feature filter allowed)
	if requestFeatureFilter {
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
	}
	// return success response
	context.Header("Content-Type", "application/json")
	context.Header("Cache-Control", "no-cache")
	context.JSON(http.StatusOK, response)
	return
}

func HandleNFRegisterOrNFSharedDataCompleteReplacement(context *gin.Context) {
	// check allowedSharedData feature enable
	if !NRFConfigure.AllowedSharedData {
		context.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "SharedData feature not allowed"})
		L.Info("NFRegisterOrNFSharedDataCompleteReplacement abort caused by SharedData feature not allowed:", context.Request)
		return
	}
	// extract sharedDataId from request uri
	sharedDataId := strings.ToLower(context.Param("sharedDataId"))
	fmt.Println("sharedDataId:", sharedDataId)
	// found sharedDataId in database
	exists := func() bool {
		NRFService.mutex.RLock()
		defer NRFService.mutex.RUnlock()
		for _, repositories := range NRFService.repositories {
			for _, v := range repositories {
				if v.SharedDataId == sharedDataId {
					return true
				}
			}
		}
		return false
	}()
	if !exists {
		// NFRegister (SharedData)
		HandleNFRegisterSharedData(context)
	} else {
		// NFUpdate (SharedData Complete Replacement)
		HandleNFSharedDataCompleteReplacement(context)
	}
}

func HandleNFRegisterSharedData(context *gin.Context) {
	var request SharedData
	// record context in logs
	L.Info("NFRegister (SharedData) request:", context.Request)
	// check request body bind json
	L.Debug("Start bind NFRegister (SharedData) request body to json:", context.Request.Body)
	err := context.ShouldBindJSON(&request)
	if err != nil {
		var problemDetails ProblemDetails
		problemDetails.Title = "Bad Request"
		problemDetails.Status = http.StatusBadRequest
		problemDetails.Detail = err.Error()
		context.Header("Content-Type", "application/problem+json")
		context.JSON(http.StatusBadRequest, problemDetails)
		L.Error("NFRegister (SharedData) request body bind json failed:", err)
		return
	}
	L.Debug("NFRegister (SharedData) request body bind json success.")
	// check request body IEs
	b, err := checkNFRegisterSharedDataIEs(&request)
	if b == false && err != nil {
		var problemDetails ProblemDetails
		problemDetails.Title = "Bad Request"
		problemDetails.Status = http.StatusBadRequest
		problemDetails.Detail = err.Error()
		context.Header("Content-Type", "application/problem+json")
		context.JSON(http.StatusBadRequest, problemDetails)
		L.Error("NFRegister (SharedData) request check failed:", err)
		return
	}
	// handle request body IEs
	response := request
	err = handleNFRegisterSharedDataIEs(&response)
	if err != nil {
		var problemDetails ProblemDetails
		problemDetails.Title = "Internal Server Error"
		problemDetails.Status = http.StatusInternalServerError
		problemDetails.Detail = err.Error()
		context.Header("Content-Type", "application/problem+json")
		context.JSON(http.StatusInternalServerError, problemDetails)
		L.Error("NFRegister (SharedData) request body handle failed:", err)
		return
	}
	// extract sharedDataId from request uri
	sharedDataId := strings.ToLower(context.Param("sharedDataId"))
	fmt.Println("sharedDataId:", sharedDataId)
	// create repository from request body
	repository := SharedRepository{
		SharedDataId:      sharedDataId,
		SharedProfileData: response.SharedProfileData,
		SharedServiceData: response.SharedServiceData,
	}
	// store repository in NRF Service database
	func() {
		NRFService.mutex.Lock()
		defer NRFService.mutex.Unlock()
		NRFService.repositories[response.SharedDataId] = append(NRFService.repositories[response.SharedDataId], repository)
	}()
	// return success response
	context.Header("Content-Type", "application/json")
	context.Header("Location", formLocation(context, "nnrf-nfm", "v1", "shared-data", sharedDataId))
	context.JSON(http.StatusCreated, response)
	return
}

func HandleNFSharedDataCompleteReplacement(context *gin.Context) {
	var request SharedData
	L.Info("NFSharedDataCompleteReplacement request:", context.Request)
	// check request body bind json
	L.Debug("Start bind NFSharedDataCompleteReplacement request body to json:", context.Request.Body)
	err := context.ShouldBindJSON(&request)
	if err != nil {
		var problemDetails ProblemDetails
		problemDetails.Title = "Bad Request"
		problemDetails.Status = http.StatusBadRequest
		problemDetails.Detail = err.Error()
		context.Header("Content-Type", "application/problem+json")
		context.JSON(http.StatusBadRequest, problemDetails)
		L.Error("NFSharedDataCompleteReplacement request body bind json failed:", err)
		return
	}
	L.Debug("NFSharedDataCompleteReplacement request body bind json success.")
	// check request body IEs
	b, err := checkNFRegisterSharedDataIEs(&request)
	if b == false && err != nil {
		var problemDetails ProblemDetails
		problemDetails.Title = "Bad Request"
		problemDetails.Status = http.StatusBadRequest
		problemDetails.Detail = err.Error()
		context.Header("Content-Type", "application/problem+json")
		context.JSON(http.StatusBadRequest, problemDetails)
		L.Error("NFSharedDataCompleteReplacement request check failed:", err)
		return
	}
	// handle request body IEs
	response := request
	err = handleNFRegisterSharedDataIEs(&response)
	if err != nil {
		var problemDetails ProblemDetails
		problemDetails.Title = "Internal Server Error"
		problemDetails.Status = http.StatusInternalServerError
		problemDetails.Detail = err.Error()
		context.Header("Content-Type", "application/problem+json")
		context.JSON(http.StatusInternalServerError, problemDetails)
		L.Error("NFSharedDataCompleteReplacement request body handle failed:", err)
		return
	}
	// extract sharedDataId from request uri
	sharedDataId := strings.ToLower(context.Param("sharedDataId"))
	fmt.Println("sharedDataId:", sharedDataId)
	// create repository from request body
	repository := SharedRepository{
		SharedDataId:      sharedDataId,
		SharedProfileData: response.SharedProfileData,
		SharedServiceData: response.SharedServiceData,
	}
	// store repository in SharedRepositories database
	err = func(repo *SharedRepository) (err error) {
		NRFService.mutex.RLock()
		defer NRFService.mutex.RUnlock()
		for _, repositories := range NRFService.repositories {
			for k, v := range repositories {
				if v.SharedDataId == sharedDataId {
					repositories[k], err = *repo, nil
					return err
				}
			}
		}
		err = errors.New("SharedRepositories not found")
		return err
	}(&repository)
	if err != nil {
		var problemDetails ProblemDetails
		problemDetails.Title = "Internal Server Error"
		problemDetails.Status = http.StatusInternalServerError
		problemDetails.Detail = err.Error()
		context.Header("Content-Type", "application/problem+json")
		context.JSON(http.StatusInternalServerError, problemDetails)
		L.Error("NFSharedDataCompleteReplacement profile complete replacement failed:", err)
		return
	}
	// return success response
	context.Header("Content-Type", "application/json")
	context.JSON(http.StatusOK, response)
	return
}

func HandleNFSharedDataRetrieve(context *gin.Context) {
	var request NFProfileRetrieveRequest
	var requestFeatureFilter bool
	// record context in logs
	L.Info("NFSharedDataRetrieve request:", context.Request)
	// check request body bind json
	L.Debug("Start bind NFSharedDataRetrieve request body to json:", context.Request.Body)
	err := context.ShouldBindJSON(&request)
	if err != nil {
		requestFeatureFilter = false
		L.Debug("NFSharedDataRetrieve request body bind json failed:", err.Error())
		L.Debug("NFSharedDataRetrieve request-feature filter not allowed.")
	} else {
		requestFeatureFilter = true
		L.Debug("NFSharedDataRetrieve request body bind json success.")
		L.Debug("NFSharedDataRetrieve request-feature filter allowed.")
	}
	// extract sharedDataId from request uri
	sharedDataId := strings.ToLower(context.Param("sharedDataId"))
	fmt.Println("sharedDataId:", sharedDataId)
	// store instance in NRF Service database
	var response SharedRepository
	exists := func(repo *SharedRepository) bool {
		NRFService.mutex.RLock()
		defer NRFService.mutex.RUnlock()
		for _, repositories := range NRFService.repositories {
			for _, v := range repositories {
				if v.SharedDataId == sharedDataId {
					*repo = v
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
		problemDetails.Detail = errors.New("SharedDataId not found").Error()
		context.Header("Content-Type", "application/problem+json")
		context.JSON(http.StatusNotFound, problemDetails)
		L.Error("NFSharedDataRetrieve request SharedData not found:", err)
		return
	}
	// check match request features (request-feature filter allowed)
	if requestFeatureFilter {
		var supported []string
		// SharedProfileData and SharedServiceData are conditional...
		for _, v := range response.SharedProfileData.NFServices {
			supported = append(supported, v.SupportedFeatures)
		}
		if !matchFeatures(request.RequesterFeatures, supported) {
			var problemDetails ProblemDetails
			problemDetails.Title = "Forbidden"
			problemDetails.Status = http.StatusForbidden
			problemDetails.Detail = errors.New("request Features not supported").Error()
			context.Header("Content-Type", "application/problem+json")
			context.JSON(http.StatusForbidden, problemDetails)
			L.Error("NFSharedDataRetrieve request features not supported:", err)
			return
		}
	}
	// return success response
	context.Header("Content-Type", "application/json")
	context.Header("Cache-Control", "no-cache")
	context.JSON(http.StatusOK, response)
	return
}

func HandleNFListRetrieve(context *gin.Context) {
	var request NFListRetrieveRequest
	// record context in logs
	L.Info("NFListRetrieve request:", context.Request)
	// check request body bind json
	L.Debug("Start bind NFListRetrieve request body to json:", context.Request.Body)
	err := context.ShouldBindQuery(&request)
	if err != nil {
		var problemDetails ProblemDetails
		problemDetails.Title = "Bad Request"
		problemDetails.Status = http.StatusBadRequest
		problemDetails.Detail = err.Error()
		context.Header("Content-Type", "application/problem+json")
		context.JSON(http.StatusBadRequest, problemDetails)
		L.Error("NFListRetrieve request body bind json failed:", err)
		return
	}
	// handle query parameters
	handleNFListRetrieveQuery(&request)
	// get instances in NRF Service database
	response, err := func(request NFListRetrieveRequest) (uriList UriList, err error) {
		NRFService.mutex.RLock()
		defer NRFService.mutex.RUnlock()
		// get instances according to nfType
		if request.NFType != "" {
			// search specific nfType
			for k, v := range NRFService.instances {
				if k == request.NFType {
					// start and end points in slices
					start := (request.PageNumber - 1) * request.PageSize
					end := request.PageNumber + request.PageSize
					// check validation of slices
					if start >= len(v) {
						err = errors.New("NFListRetrieveRequest start index out of bounds")
						return
					}
					if end > len(v) {
						end = len(v)
					}
					if (end - start) > request.Limit {
						end = start + request.Limit
					}
					// retrieve NFs uri list
					uriList.TotalItemCount = end - start
					for _, j := range v[start:end] {
						uriList.Links = append(uriList.Links, formLocation(context, "nnrf-nfm", "v1", "nf-instances", j.NFInstanceId))
					}
					break
				}
			}
		}
		return uriList, err
	}(request)
	if err != nil {
		var problemDetails ProblemDetails
		problemDetails.Title = "Not Found"
		problemDetails.Status = http.StatusNotFound
		problemDetails.Detail = errors.New("UriList not found").Error() + ":" + err.Error()
		context.Header("Content-Type", "application/problem+json")
		context.JSON(http.StatusNotFound, problemDetails)
		L.Error("NFListRetrieve request query UriList not found:", err)
	}
	// return success response
	context.Header("Content-Type", "application/3gppHal+json")
	context.JSON(http.StatusOK, response)
	return
}
