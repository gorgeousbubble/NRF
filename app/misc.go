package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	. "nrf/data"
	. "nrf/logs"
	. "nrf/util"
	"strings"
)

func checkNFRegisterIEs(request *NFProfile) (b bool, err error) {
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

func handleNFRegisterIEs(request *NFProfile) (err error) {
	err = nil
	// handle NFInstanceId
	L.Debug("Start HandleNFInstanceId:", request.NFInstanceId)
	err = HandleNFInstanceId(&request.NFInstanceId)
	if err != nil {
		L.Error("HandleNFInstanceId failed:", err)
		return err
	}
	L.Debug("HandleNFInstanceId success:", request.NFInstanceId)
	// handle HeartBeatTimer
	L.Debug("Start HandleHeartBeatTimer:", request.HeartBeatTimer)
	err = HandleHeartBeatTimer(&request.HeartBeatTimer)
	if err != nil {
		L.Error("HandleHeartBeatTimer failed:", err)
		return err
	}
	L.Debug("HandleHeartBeatTimer success.")
	return err
}

func checkNFRegisterSharedDataIEs(request *SharedData) (b bool, err error) {
	b, err = true, nil
	// check mandatory IEs...
	// check SharedDataId
	L.Debug("Start CheckSharedDataId:", request.SharedDataId)
	b, err = CheckSharedDataId(request.SharedDataId)
	if err != nil {
		b = false
		L.Error("CheckSharedDataId failed:", err)
		return b, err
	}
	L.Debug("CheckSharedDataId success.")
	return b, err
}

func handleNFRegisterSharedDataIEs(request *SharedData) (err error) {
	err = nil
	// handle NFInstanceId
	L.Debug("Start HandleSharedDataId:", request.SharedDataId)
	err = HandleSharedDataId(&request.SharedDataId)
	if err != nil {
		L.Error("HandleSharedDataId failed:", err)
		return err
	}
	L.Debug("HandleSharedDataId success:", request.SharedDataId)
	return err
}

func matchFeatures(required, supported []string) bool {
	supportedSet := make(map[string]struct{})
	for _, s := range supported {
		supportedSet[s] = struct{}{}
	}
	for _, r := range required {
		if _, exist := supportedSet[r]; !exist {
			return false
		}
	}
	return true
}

func autodetectHttpProtocol(context *gin.Context) (protocol string) {
	// autodetect http protocol from header "X-Forwarded-Proto"
	if protocol = context.GetHeader("X-Forwarded-Proto"); protocol != "" {
		return protocol
	}
	// autodetect http protocol from header "X-Forwarded-Scheme"
	if protocol = context.GetHeader("X-Forwarded-Scheme"); protocol != "" {
		return protocol
	}
	// autodetect http protocol from request tls type
	if context.Request.TLS != nil {
		protocol = "https"
	} else {
		protocol = "http"
	}
	return protocol
}

func autodetectHttpHost(context *gin.Context) (host string) {
	// autodetect http host from header "X-Forwarded-Host"
	if host = context.GetHeader("X-Forwarded-Host"); host != "" {
		return host
	}
	// autodetect http host from request host
	host = context.Request.Host
	if !strings.Contains(host, ":") {
		switch autodetectHttpProtocol(context) {
		case "http":
			host += ":80"
		case "https":
			host += ":443"
		default:
			host += fetchListenPort(context)
		}
	}
	return host
}

func fetchListenPort(context *gin.Context) (port string) {
	addr := context.Request.Context().Value(http.LocalAddrContextKey).(net.Addr)
	_, port, _ = net.SplitHostPort(addr.String())
	return port
}

func formLocation(context *gin.Context, apiName string, apiVersion string, resource string, identity string) (location string) {
	return fmt.Sprintf("%s://%s/%s/%s/%s/%s", autodetectHttpProtocol(context), autodetectHttpHost(context), apiName, apiVersion, resource, identity)
}

func handleNFListRetrieveQuery(request *NFListRetrieveRequest) {
	// handle Limit
	L.Debug("Start HandleLimit", request.Limit)
	if request.Limit == 0 {
		request.Limit = 1
	}
	L.Debug("HandleLimit success:", request.Limit)
	// handle HandlePageNumber
	L.Debug("Start HandlePageNumber:", request.PageNumber)
	if request.PageNumber == 0 {
		request.PageNumber = 1
	}
	L.Debug("HandlePageSize success.")
	// handle HandlePageSize
	L.Debug("Start HandlePageSize:", request.PageSize)
	if request.PageSize == 0 {
		request.PageSize = 1
	}
	L.Debug("HandlePageSize success.")
	return
}
