package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleRegister(t *testing.T) {
	NRFService = New()
	NRFService.Init()
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	// simulate POST request
	url := "http://localhost:8080/nnrf-nfm/v1/nf-instances/331a1fb2-3ac1-43df-a7d0-882d0ee44b7d"
	body := strings.NewReader(`{"nfInstanceID":"331a1fb2-3ac1-43df-a7d0-882d0ee44b7d", "nfType":"AMF", "nfStatus":"REGISTERED"}`)
	context.Request, _ = http.NewRequest(http.MethodPut, url, body)
	// test HandleRegister function
	HandleNFRegister(context)
}
