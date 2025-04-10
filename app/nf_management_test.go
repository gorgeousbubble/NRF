package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	. "nrf/data"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() *gin.Engine {
	// initialize Gin framework
	gin.SetMode(gin.TestMode)
	router := gin.New()
	// middleware handle functions
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(ContentEncodingMiddleware())
	router.Use(AcceptEncodingMiddleware())
	router.Use(SecurityHeadersMiddleware())
	router.Use(ETagMiddleware(defaultConfig))
	// API route groups
	nfManagement := router.Group("/nnrf-nfm/v1")
	{
		nfManagement.GET("nf-instances", HandleNFListRetrieve)
		nfManagement.PUT("nf-instances/:nfInstanceID", HandleNFRegisterOrNFProfileCompleteReplacement)
		nfManagement.GET("nf-instances/:nfInstanceID", HandleNFProfileRetrieve)
		nfManagement.PUT("shared-data/:sharedDataId", HandleNFRegisterOrNFSharedDataCompleteReplacement)
		nfManagement.GET("shared-data/:sharedDataId", HandleNFSharedDataRetrieve)
	}
	return router
}

func startTestServer() (*httptest.Server, *gin.Engine) {
	router := setupTestRouter()
	return httptest.NewServer(router), router
}

func TestHandleNFRegister(t *testing.T) {
	// initialize NRF Service
	NRFService = New()
	err := NRFService.Init()
	if err != nil {
		t.Error(err)
	}
	// start http test service
	server, router := startTestServer()
	defer server.Close()
	// construct network function request content
	url := server.URL + "/nnrf-nfm/v1/nf-instances"
	nfInstanceId := uuid.New().String()
	nfType := "AMF"
	nfStatus := "REGISTERED"
	// assemble network function http request
	profile := NFProfile{
		NFInstanceId: nfInstanceId,
		NFType:       nfType,
		NFStatus:     nfStatus,
	}
	body, err := json.Marshal(profile)
	if err != nil {
		t.Errorf("Error marshalling profile: %v", err)
	}
	// http request NFRegister
	w := httptest.NewRecorder()
	request, err := http.NewRequest("PUT", url+"/"+nfInstanceId, bytes.NewReader(body))
	if err != nil {
		t.Errorf("Error creating request: %v", err)
	}
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, request)
	var response NFProfile
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error unmarshalling response: %v", err)
	}
	// assert http response
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, url+"/"+nfInstanceId, w.Header().Get("Location"))
	assert.Equal(t, nfInstanceId, response.NFInstanceId)
	assert.Equal(t, nfType, response.NFType)
	assert.Equal(t, nfStatus, response.NFStatus)
}

func BenchmarkHandleNFRegister(b *testing.B) {
	// initialize NRF Service
	NRFService = New()
	err := NRFService.Init()
	if err != nil {
		b.Error(err)
	}
	// start http test service
	server, router := startTestServer()
	defer server.Close()
	// start benchmark test...
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// construct network function request content
		rand.New(rand.NewSource(time.Now().UnixNano()))
		url := server.URL + "/nnrf-nfm/v1/nf-instances"
		nfInstanceId := uuid.New().String()
		nfType := NetworkFunctionType[rand.Intn(len(NetworkFunctionType))]
		nfStatus := NetworkFunctionStatus[rand.Intn(len(NetworkFunctionStatus))]
		// assemble network function http request
		profile := NFProfile{
			NFInstanceId: nfInstanceId,
			NFType:       nfType,
			NFStatus:     nfStatus,
		}
		body, err := json.Marshal(profile)
		if err != nil {
			b.Errorf("Error marshalling profile: %v", err)
		}
		// http request NFRegister
		w := httptest.NewRecorder()
		request, err := http.NewRequest("PUT", url+"/"+nfInstanceId, bytes.NewReader(body))
		if err != nil {
			b.Errorf("Error creating request: %v", err)
		}
		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, request)
		var response NFProfile
		err = json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			b.Errorf("Error unmarshalling response: %v", err)
		}
		// assert http response
		assert.Equal(b, http.StatusCreated, w.Code)
		assert.Equal(b, "application/json", w.Header().Get("Content-Type"))
		assert.Equal(b, url+"/"+nfInstanceId, w.Header().Get("Location"))
		assert.Equal(b, nfInstanceId, response.NFInstanceId)
		assert.Equal(b, nfType, response.NFType)
		assert.Equal(b, nfStatus, response.NFStatus)
	}
}

func BenchmarkHandleNFRegisterParallel(b *testing.B) {
	// initialize NRF Service
	NRFService = New()
	err := NRFService.Init()
	if err != nil {
		b.Error(err)
	}
	// start http test service
	server, router := startTestServer()
	defer server.Close()
	// start benchmark test...
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// construct network function request content
			rand.New(rand.NewSource(time.Now().UnixNano()))
			url := server.URL + "/nnrf-nfm/v1/nf-instances"
			nfInstanceId := uuid.New().String()
			nfType := NetworkFunctionType[rand.Intn(len(NetworkFunctionType))]
			nfStatus := NetworkFunctionStatus[rand.Intn(len(NetworkFunctionStatus))]
			// assemble network function http request
			profile := NFProfile{
				NFInstanceId: nfInstanceId,
				NFType:       nfType,
				NFStatus:     nfStatus,
			}
			body, err := json.Marshal(profile)
			if err != nil {
				b.Errorf("Error marshalling profile: %v", err)
			}
			// http request NFRegister
			w := httptest.NewRecorder()
			request, err := http.NewRequest("PUT", url+"/"+nfInstanceId, bytes.NewReader(body))
			if err != nil {
				b.Errorf("Error creating request: %v", err)
			}
			request.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, request)
			var response NFProfile
			err = json.Unmarshal(w.Body.Bytes(), &response)
			if err != nil {
				b.Errorf("Error unmarshalling response: %v", err)
			}
			// assert http response
			assert.Equal(b, http.StatusCreated, w.Code)
			assert.Equal(b, "application/json", w.Header().Get("Content-Type"))
			assert.Equal(b, url+"/"+nfInstanceId, w.Header().Get("Location"))
			assert.Equal(b, nfInstanceId, response.NFInstanceId)
			assert.Equal(b, nfType, response.NFType)
			assert.Equal(b, nfStatus, response.NFStatus)
		}
	})
}

func FuzzHandleNFRegister(f *testing.F) {
	// initialize NRF Service
	NRFService = New()
	err := NRFService.Init()
	if err != nil {
		f.Error(err)
	}
	// start http test service
	server, router := startTestServer()
	defer server.Close()
	// construct network function request content
	url := server.URL + "/nnrf-nfm/v1/nf-instances"
	nfInstanceId := uuid.New().String()
	nfType := "AMF"
	nfStatus := "REGISTERED"
	// add fuzzy test database
	f.Add([]byte(`{"nfInstanceId":"` + nfInstanceId + `","nfType":"` + nfType + `","nfStatus":"` + nfStatus + `"}`))
	f.Fuzz(func(t *testing.T, data []byte) {
		// unexpected result
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("panic: %v", r)
			}
		}()
		// http request NFRegister
		w := httptest.NewRecorder()
		request, err := http.NewRequest("PUT", url+"/"+nfInstanceId, bytes.NewReader(data))
		if err != nil {
			t.Errorf("Error creating request: %v", err)
		}
		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, request)
		var response NFProfile
		err = json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("Error unmarshalling response: %v", err)
		}
		// assert http response
		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
		assert.Equal(t, url+"/"+nfInstanceId, w.Header().Get("Location"))
		assert.Equal(t, nfInstanceId, response.NFInstanceId)
		assert.Equal(t, nfType, response.NFType)
		assert.Equal(t, nfStatus, response.NFStatus)
	})
}

func ExampleHandleNFRegister() {
	// initialize NRF Service
	NRFService = New()
	_ = NRFService.Init()
	// start http test service
	server, router := startTestServer()
	defer server.Close()
	// construct network function request content
	url := server.URL + "/nnrf-nfm/v1/nf-instances"
	nfInstanceId := uuid.New().String()
	nfType := "AMF"
	nfStatus := "REGISTERED"
	// assemble network function http request
	profile := NFProfile{
		NFInstanceId: nfInstanceId,
		NFType:       nfType,
		NFStatus:     nfStatus,
	}
	body, _ := json.Marshal(profile)
	// http request NFRegister
	w := httptest.NewRecorder()
	request, _ := http.NewRequest("PUT", url+"/"+nfInstanceId, bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, request)
	var response NFProfile
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	// print http response
	fmt.Println("response status code:", w.Code)
	fmt.Println("response content-type:", w.Header().Get("Content-Type"))
	fmt.Println("response location:", w.Header().Get("Location"))
	fmt.Println("response body:", string(w.Body.Bytes()))
	fmt.Println("response nfInstanceId:", response.NFInstanceId)
	fmt.Println("response nfType:", response.NFType)
	fmt.Println("response nfStatus:", response.NFStatus)
}

func TestHandleNFRegisterWithUpperNFInstanceID(t *testing.T) {
	// initialize NRF Service
	NRFService = New()
	err := NRFService.Init()
	if err != nil {
		t.Error(err)
	}
	// start http test service
	server, router := startTestServer()
	defer server.Close()
	// construct network function request content
	url := server.URL + "/nnrf-nfm/v1/nf-instances"
	nfInstanceId := strings.ToUpper(uuid.New().String())
	nfType := "UPF"
	nfStatus := "REGISTERED"
	// assemble network function http request
	profile := NFProfile{
		NFInstanceId: nfInstanceId,
		NFType:       nfType,
		NFStatus:     nfStatus,
	}
	body, err := json.Marshal(profile)
	if err != nil {
		t.Errorf("Error marshalling profile: %v", err)
	}
	// http request NFRegister
	w := httptest.NewRecorder()
	request, err := http.NewRequest("PUT", url+"/"+nfInstanceId, bytes.NewReader(body))
	if err != nil {
		t.Errorf("Error creating request: %v", err)
	}
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, request)
	var response NFProfile
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error unmarshalling response: %v", err)
	}
	// assert http response
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, url+"/"+strings.ToLower(nfInstanceId), w.Header().Get("Location"))
	assert.Equal(t, strings.ToLower(nfInstanceId), response.NFInstanceId)
	assert.Equal(t, nfType, response.NFType)
	assert.Equal(t, nfStatus, response.NFStatus)
}

func BenchmarkHandleNFRegisterWithUpperNFInstanceID(b *testing.B) {
	// initialize NRF Service
	NRFService = New()
	err := NRFService.Init()
	if err != nil {
		b.Error(err)
	}
	// start http test service
	server, router := startTestServer()
	defer server.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// construct network function request content
		rand.New(rand.NewSource(time.Now().UnixNano()))
		url := server.URL + "/nnrf-nfm/v1/nf-instances"
		nfInstanceId := strings.ToUpper(uuid.New().String())
		nfType := NetworkFunctionType[rand.Intn(len(NetworkFunctionType))]
		nfStatus := NetworkFunctionStatus[rand.Intn(len(NetworkFunctionStatus))]
		// assemble network function http request
		profile := NFProfile{
			NFInstanceId: nfInstanceId,
			NFType:       nfType,
			NFStatus:     nfStatus,
		}
		body, err := json.Marshal(profile)
		if err != nil {
			b.Errorf("Error marshalling profile: %v", err)
		}
		// http request NFRegister
		w := httptest.NewRecorder()
		request, err := http.NewRequest("PUT", url+"/"+nfInstanceId, bytes.NewReader(body))
		if err != nil {
			b.Errorf("Error creating request: %v", err)
		}
		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, request)
		var response NFProfile
		err = json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			b.Errorf("Error unmarshalling response: %v", err)
		}
		// assert http response
		assert.Equal(b, http.StatusCreated, w.Code)
		assert.Equal(b, "application/json", w.Header().Get("Content-Type"))
		assert.Equal(b, url+"/"+strings.ToLower(nfInstanceId), w.Header().Get("Location"))
		assert.Equal(b, strings.ToLower(nfInstanceId), response.NFInstanceId)
		assert.Equal(b, nfType, response.NFType)
		assert.Equal(b, nfStatus, response.NFStatus)
	}
}

func BenchmarkHandleNFRegisterWithUpperNFInstanceIDParallel(b *testing.B) {
	// initialize NRF Service
	NRFService = New()
	err := NRFService.Init()
	if err != nil {
		b.Error(err)
	}
	// start http test service
	server, router := startTestServer()
	defer server.Close()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// construct network function request content
			rand.New(rand.NewSource(time.Now().UnixNano()))
			url := server.URL + "/nnrf-nfm/v1/nf-instances"
			nfInstanceId := strings.ToUpper(uuid.New().String())
			nfType := NetworkFunctionType[rand.Intn(len(NetworkFunctionType))]
			nfStatus := NetworkFunctionStatus[rand.Intn(len(NetworkFunctionStatus))]
			// assemble network function http request
			profile := NFProfile{
				NFInstanceId: nfInstanceId,
				NFType:       nfType,
				NFStatus:     nfStatus,
			}
			body, err := json.Marshal(profile)
			if err != nil {
				b.Errorf("Error marshalling profile: %v", err)
			}
			// http request NFRegister
			w := httptest.NewRecorder()
			request, err := http.NewRequest("PUT", url+"/"+nfInstanceId, bytes.NewReader(body))
			if err != nil {
				b.Errorf("Error creating request: %v", err)
			}
			request.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, request)
			var response NFProfile
			err = json.Unmarshal(w.Body.Bytes(), &response)
			if err != nil {
				b.Errorf("Error unmarshalling response: %v", err)
			}
			// assert http response
			assert.Equal(b, http.StatusCreated, w.Code)
			assert.Equal(b, "application/json", w.Header().Get("Content-Type"))
			assert.Equal(b, url+"/"+strings.ToLower(nfInstanceId), w.Header().Get("Location"))
			assert.Equal(b, strings.ToLower(nfInstanceId), response.NFInstanceId)
			assert.Equal(b, nfType, response.NFType)
			assert.Equal(b, nfStatus, response.NFStatus)
		}
	})
}

func TestHandleNFRegisterWithoutNFType(t *testing.T) {
	// initialize NRF Service
	NRFService = New()
	err := NRFService.Init()
	if err != nil {
		t.Error(err)
	}
	// start http test service
	server, router := startTestServer()
	defer server.Close()
	// construct network function request content
	url := server.URL + "/nnrf-nfm/v1/nf-instances"
	nfInstanceId := uuid.New().String()
	nfStatus := "REGISTERED"
	// assemble network function http request
	profile := NFProfile{
		NFInstanceId: nfInstanceId,
		NFStatus:     nfStatus,
	}
	body, err := json.Marshal(profile)
	if err != nil {
		t.Errorf("Error marshalling profile: %v", err)
	}
	// http request NFRegister
	w := httptest.NewRecorder()
	request, err := http.NewRequest("PUT", url+"/"+nfInstanceId, bytes.NewReader(body))
	if err != nil {
		t.Errorf("Error creating request: %v", err)
	}
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, request)
	var response NFProfileRegistrationError
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error unmarshalling response: %v", err)
	}
	// assert http response
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "application/problem+json", w.Header().Get("Content-Type"))
	assert.Equal(t, "Bad Request", response.ProblemDetails.Title)
	assert.Equal(t, "Key: 'NFProfile.NFType' Error:Field validation for 'NFType' failed on the 'required' tag", response.ProblemDetails.Detail)
}

func BenchmarkHandleNFRegisterWithoutNFType(b *testing.B) {
	// initialize NRF Service
	NRFService = New()
	err := NRFService.Init()
	if err != nil {
		b.Error(err)
	}
	// start http test service
	server, router := startTestServer()
	defer server.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// construct network function request content
		rand.New(rand.NewSource(time.Now().UnixNano()))
		url := server.URL + "/nnrf-nfm/v1/nf-instances"
		nfInstanceId := strings.ToUpper(uuid.New().String())
		nfStatus := NetworkFunctionStatus[rand.Intn(len(NetworkFunctionStatus))]
		// assemble network function http request
		profile := NFProfile{
			NFInstanceId: nfInstanceId,
			NFStatus:     nfStatus,
		}
		body, err := json.Marshal(profile)
		if err != nil {
			b.Errorf("Error marshalling profile: %v", err)
		}
		// http request NFRegister
		w := httptest.NewRecorder()
		request, err := http.NewRequest("PUT", url+"/"+nfInstanceId, bytes.NewReader(body))
		if err != nil {
			b.Errorf("Error creating request: %v", err)
		}
		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, request)
		var response NFProfileRegistrationError
		err = json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			b.Errorf("Error unmarshalling response: %v", err)
		}
		// assert http response
		assert.Equal(b, http.StatusBadRequest, w.Code)
		assert.Equal(b, "application/problem+json", w.Header().Get("Content-Type"))
		assert.Equal(b, "Bad Request", response.ProblemDetails.Title)
		assert.Equal(b, "Key: 'NFProfile.NFType' Error:Field validation for 'NFType' failed on the 'required' tag", response.ProblemDetails.Detail)
	}
}

func BenchmarkHandleNFRegisterWithoutNFTypeParallel(b *testing.B) {
	// initialize NRF Service
	NRFService = New()
	err := NRFService.Init()
	if err != nil {
		b.Error(err)
	}
	// start http test service
	server, router := startTestServer()
	defer server.Close()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// construct network function request content
			rand.New(rand.NewSource(time.Now().UnixNano()))
			url := server.URL + "/nnrf-nfm/v1/nf-instances"
			nfInstanceId := strings.ToUpper(uuid.New().String())
			nfStatus := NetworkFunctionStatus[rand.Intn(len(NetworkFunctionStatus))]
			// assemble network function http request
			profile := NFProfile{
				NFInstanceId: nfInstanceId,
				NFStatus:     nfStatus,
			}
			body, err := json.Marshal(profile)
			if err != nil {
				b.Errorf("Error marshalling profile: %v", err)
			}
			// http request NFRegister
			w := httptest.NewRecorder()
			request, err := http.NewRequest("PUT", url+"/"+nfInstanceId, bytes.NewReader(body))
			if err != nil {
				b.Errorf("Error creating request: %v", err)
			}
			request.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, request)
			var response NFProfileRegistrationError
			err = json.Unmarshal(w.Body.Bytes(), &response)
			if err != nil {
				b.Errorf("Error unmarshalling response: %v", err)
			}
			// assert http response
			assert.Equal(b, http.StatusBadRequest, w.Code)
			assert.Equal(b, "application/problem+json", w.Header().Get("Content-Type"))
			assert.Equal(b, "Bad Request", response.ProblemDetails.Title)
			assert.Equal(b, "Key: 'NFProfile.NFType' Error:Field validation for 'NFType' failed on the 'required' tag", response.ProblemDetails.Detail)
		}
	})
}

func TestHandleNFProfileCompleteReplacement(t *testing.T) {
	// initialize NRF Service
	NRFService = New()
	err := NRFService.Init()
	if err != nil {
		t.Error(err)
	}
	// start http test service
	server, router := startTestServer()
	defer server.Close()
	// construct network function request content
	url := server.URL + "/nnrf-nfm/v1/nf-instances"
	nfInstanceId := uuid.New().String()
	nfType := "AMF"
	nfStatus := "REGISTERED"
	// assemble network function http request
	profile := NFProfile{
		NFInstanceId: nfInstanceId,
		NFType:       nfType,
		NFStatus:     nfStatus,
	}
	body, err := json.Marshal(profile)
	if err != nil {
		t.Errorf("Error marshalling profile: %v", err)
	}
	// http request NFRegister
	w := httptest.NewRecorder()
	request, err := http.NewRequest("PUT", url+"/"+nfInstanceId, bytes.NewReader(body))
	if err != nil {
		t.Errorf("Error creating request: %v", err)
	}
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, request)
	var response NFProfile
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error unmarshalling response: %v", err)
	}
	// assert http response
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, url+"/"+nfInstanceId, w.Header().Get("Location"))
	assert.Equal(t, nfInstanceId, response.NFInstanceId)
	assert.Equal(t, nfType, response.NFType)
	assert.Equal(t, nfStatus, response.NFStatus)
	// construct network function request content
	nfStatusNew := "SUSPENDED"
	// assemble network function http request
	profileNew := NFProfile{
		NFInstanceId: nfInstanceId,
		NFType:       nfType,
		NFStatus:     nfStatusNew,
	}
	bodyNew, err := json.Marshal(profileNew)
	if err != nil {
		t.Errorf("Error marshalling profile: %v", err)
	}
	// http request NFProfileCompleteReplacement
	w = httptest.NewRecorder()
	request, err = http.NewRequest("PUT", url+"/"+nfInstanceId, bytes.NewReader(bodyNew))
	if err != nil {
		t.Errorf("Error creating request: %v", err)
	}
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, request)
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error unmarshalling response: %v", err)
	}
	// assert http response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, nfInstanceId, response.NFInstanceId)
	assert.Equal(t, nfType, response.NFType)
	assert.Equal(t, nfStatusNew, response.NFStatus)
}

func BenchmarkHandleNFProfileCompleteReplacement(b *testing.B) {
	// initialize NRF Service
	NRFService = New()
	err := NRFService.Init()
	if err != nil {
		b.Error(err)
	}
	// start http test service
	server, router := startTestServer()
	defer server.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// construct network function request content
		rand.New(rand.NewSource(time.Now().UnixNano()))
		url := server.URL + "/nnrf-nfm/v1/nf-instances"
		nfInstanceId := uuid.New().String()
		nfType := NetworkFunctionType[rand.Intn(len(NetworkFunctionType))]
		nfStatus := NetworkFunctionStatus[rand.Intn(len(NetworkFunctionStatus))]
		// assemble network function http request
		profile := NFProfile{
			NFInstanceId: nfInstanceId,
			NFType:       nfType,
			NFStatus:     nfStatus,
		}
		body, err := json.Marshal(profile)
		if err != nil {
			b.Errorf("Error marshalling profile: %v", err)
		}
		// http request NFRegister
		w := httptest.NewRecorder()
		request, err := http.NewRequest("PUT", url+"/"+nfInstanceId, bytes.NewReader(body))
		if err != nil {
			b.Errorf("Error creating request: %v", err)
		}
		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, request)
		var response NFProfile
		err = json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			b.Errorf("Error unmarshalling response: %v", err)
		}
		// assert http response
		assert.Equal(b, http.StatusCreated, w.Code)
		assert.Equal(b, "application/json", w.Header().Get("Content-Type"))
		assert.Equal(b, url+"/"+nfInstanceId, w.Header().Get("Location"))
		assert.Equal(b, nfInstanceId, response.NFInstanceId)
		assert.Equal(b, nfType, response.NFType)
		assert.Equal(b, nfStatus, response.NFStatus)
		// construct network function request content
		nfStatusNew := NetworkFunctionStatus[rand.Intn(len(NetworkFunctionStatus))]
		// assemble network function http request
		profileNew := NFProfile{
			NFInstanceId: nfInstanceId,
			NFType:       nfType,
			NFStatus:     nfStatusNew,
		}
		bodyNew, err := json.Marshal(profileNew)
		if err != nil {
			b.Errorf("Error marshalling profile: %v", err)
		}
		// http request NFProfileCompleteReplacement
		w = httptest.NewRecorder()
		request, err = http.NewRequest("PUT", url+"/"+nfInstanceId, bytes.NewReader(bodyNew))
		if err != nil {
			b.Errorf("Error creating request: %v", err)
		}
		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, request)
		err = json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			b.Errorf("Error unmarshalling response: %v", err)
		}
		// assert http response
		assert.Equal(b, http.StatusOK, w.Code)
		assert.Equal(b, "application/json", w.Header().Get("Content-Type"))
		assert.Equal(b, nfInstanceId, response.NFInstanceId)
		assert.Equal(b, nfType, response.NFType)
		assert.Equal(b, nfStatusNew, response.NFStatus)
	}
}

func BenchmarkHandleNFProfileCompleteReplacementParallel(b *testing.B) {
	// initialize NRF Service
	NRFService = New()
	err := NRFService.Init()
	if err != nil {
		b.Error(err)
	}
	// start http test service
	server, router := startTestServer()
	defer server.Close()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// construct network function request content
			rand.New(rand.NewSource(time.Now().UnixNano()))
			url := server.URL + "/nnrf-nfm/v1/nf-instances"
			nfInstanceId := uuid.New().String()
			nfType := NetworkFunctionType[rand.Intn(len(NetworkFunctionType))]
			nfStatus := NetworkFunctionStatus[rand.Intn(len(NetworkFunctionStatus))]
			// assemble network function http request
			profile := NFProfile{
				NFInstanceId: nfInstanceId,
				NFType:       nfType,
				NFStatus:     nfStatus,
			}
			body, err := json.Marshal(profile)
			if err != nil {
				b.Errorf("Error marshalling profile: %v", err)
			}
			// http request NFRegister
			w := httptest.NewRecorder()
			request, err := http.NewRequest("PUT", url+"/"+nfInstanceId, bytes.NewReader(body))
			if err != nil {
				b.Errorf("Error creating request: %v", err)
			}
			request.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, request)
			var response NFProfile
			err = json.Unmarshal(w.Body.Bytes(), &response)
			if err != nil {
				b.Errorf("Error unmarshalling response: %v", err)
			}
			// assert http response
			assert.Equal(b, http.StatusCreated, w.Code)
			assert.Equal(b, "application/json", w.Header().Get("Content-Type"))
			assert.Equal(b, url+"/"+nfInstanceId, w.Header().Get("Location"))
			assert.Equal(b, nfInstanceId, response.NFInstanceId)
			assert.Equal(b, nfType, response.NFType)
			assert.Equal(b, nfStatus, response.NFStatus)
			// construct network function request content
			nfStatusNew := NetworkFunctionStatus[rand.Intn(len(NetworkFunctionStatus))]
			// assemble network function http request
			profileNew := NFProfile{
				NFInstanceId: nfInstanceId,
				NFType:       nfType,
				NFStatus:     nfStatusNew,
			}
			bodyNew, err := json.Marshal(profileNew)
			if err != nil {
				b.Errorf("Error marshalling profile: %v", err)
			}
			// http request NFProfileCompleteReplacement
			w = httptest.NewRecorder()
			request, err = http.NewRequest("PUT", url+"/"+nfInstanceId, bytes.NewReader(bodyNew))
			if err != nil {
				b.Errorf("Error creating request: %v", err)
			}
			request.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, request)
			err = json.Unmarshal(w.Body.Bytes(), &response)
			if err != nil {
				b.Errorf("Error unmarshalling response: %v", err)
			}
			// assert http response
			assert.Equal(b, http.StatusOK, w.Code)
			assert.Equal(b, "application/json", w.Header().Get("Content-Type"))
			assert.Equal(b, nfInstanceId, response.NFInstanceId)
			assert.Equal(b, nfType, response.NFType)
			assert.Equal(b, nfStatusNew, response.NFStatus)
		}
	})
}

func TestHandleNFProfileRetrieve(t *testing.T) {
	// initialize NRF Service
	NRFService = New()
	err := NRFService.Init()
	if err != nil {
		t.Error(err)
	}
	// start http test service
	server, router := startTestServer()
	defer server.Close()
	// construct network function request content
	url := server.URL + "/nnrf-nfm/v1/nf-instances"
	nfInstanceId := uuid.New().String()
	nfType := "SMF"
	nfStatus := "REGISTERED"
	// assemble network function http request
	profile := NFProfile{
		NFInstanceId: nfInstanceId,
		NFType:       nfType,
		NFStatus:     nfStatus,
	}
	body, err := json.Marshal(profile)
	if err != nil {
		t.Errorf("Error marshalling profile: %v", err)
	}
	w := httptest.NewRecorder()
	request, err := http.NewRequest("PUT", url+"/"+nfInstanceId, bytes.NewReader(body))
	if err != nil {
		t.Errorf("Error creating request: %v", err)
	}
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, request)
	var response NFProfile
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error unmarshalling response: %v", err)
	}
	// assert http response
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, url+"/"+nfInstanceId, w.Header().Get("Location"))
	assert.Equal(t, nfInstanceId, response.NFInstanceId)
	assert.Equal(t, nfType, response.NFType)
	assert.Equal(t, nfStatus, response.NFStatus)
	// http request NFProfileRetrieve
	w = httptest.NewRecorder()
	request, err = http.NewRequest("GET", url+"/"+nfInstanceId, nil)
	if err != nil {
		t.Errorf("Error creating request: %v", err)
	}
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, request)
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error unmarshalling response: %v", err)
	}
	// assert http response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, nfInstanceId, response.NFInstanceId)
	assert.Equal(t, nfType, response.NFType)
	assert.Equal(t, nfStatus, response.NFStatus)
}

func BenchmarkHandleNFProfileRetrieve(b *testing.B) {
	// initialize NRF Service
	NRFService = New()
	err := NRFService.Init()
	if err != nil {
		b.Error(err)
	}
	// start http test service
	server, router := startTestServer()
	defer server.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// construct network function request content
		rand.New(rand.NewSource(time.Now().UnixNano()))
		url := server.URL + "/nnrf-nfm/v1/nf-instances"
		nfInstanceId := uuid.New().String()
		nfType := NetworkFunctionType[rand.Intn(len(NetworkFunctionType))]
		nfStatus := NetworkFunctionStatus[rand.Intn(len(NetworkFunctionStatus))]
		// assemble network function http request
		profile := NFProfile{
			NFInstanceId: nfInstanceId,
			NFType:       nfType,
			NFStatus:     nfStatus,
		}
		body, err := json.Marshal(profile)
		if err != nil {
			b.Errorf("Error marshalling profile: %v", err)
		}
		w := httptest.NewRecorder()
		request, err := http.NewRequest("PUT", url+"/"+nfInstanceId, bytes.NewReader(body))
		if err != nil {
			b.Errorf("Error creating request: %v", err)
		}
		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, request)
		var response NFProfile
		err = json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			b.Errorf("Error unmarshalling response: %v", err)
		}
		// assert http response
		assert.Equal(b, http.StatusCreated, w.Code)
		assert.Equal(b, "application/json", w.Header().Get("Content-Type"))
		assert.Equal(b, url+"/"+nfInstanceId, w.Header().Get("Location"))
		assert.Equal(b, nfInstanceId, response.NFInstanceId)
		assert.Equal(b, nfType, response.NFType)
		assert.Equal(b, nfStatus, response.NFStatus)
		// http request NFProfileRetrieve
		w = httptest.NewRecorder()
		request, err = http.NewRequest("GET", url+"/"+nfInstanceId, nil)
		if err != nil {
			b.Errorf("Error creating request: %v", err)
		}
		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, request)
		err = json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			b.Errorf("Error unmarshalling response: %v", err)
		}
		// assert http response
		assert.Equal(b, http.StatusOK, w.Code)
		assert.Equal(b, "application/json", w.Header().Get("Content-Type"))
		assert.Equal(b, nfInstanceId, response.NFInstanceId)
		assert.Equal(b, nfType, response.NFType)
		assert.Equal(b, nfStatus, response.NFStatus)
	}
}

func BenchmarkHandleNFProfileRetrieveParallel(b *testing.B) {
	// initialize NRF Service
	NRFService = New()
	err := NRFService.Init()
	if err != nil {
		b.Error(err)
	}
	// start http test service
	server, router := startTestServer()
	defer server.Close()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// construct network function request content
			rand.New(rand.NewSource(time.Now().UnixNano()))
			url := server.URL + "/nnrf-nfm/v1/nf-instances"
			nfInstanceId := uuid.New().String()
			nfType := NetworkFunctionType[rand.Intn(len(NetworkFunctionType))]
			nfStatus := NetworkFunctionStatus[rand.Intn(len(NetworkFunctionStatus))]
			// assemble network function http request
			profile := NFProfile{
				NFInstanceId: nfInstanceId,
				NFType:       nfType,
				NFStatus:     nfStatus,
			}
			body, err := json.Marshal(profile)
			if err != nil {
				b.Errorf("Error marshalling profile: %v", err)
			}
			w := httptest.NewRecorder()
			request, err := http.NewRequest("PUT", url+"/"+nfInstanceId, bytes.NewReader(body))
			if err != nil {
				b.Errorf("Error creating request: %v", err)
			}
			request.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, request)
			var response NFProfile
			err = json.Unmarshal(w.Body.Bytes(), &response)
			if err != nil {
				b.Errorf("Error unmarshalling response: %v", err)
			}
			// assert http response
			assert.Equal(b, http.StatusCreated, w.Code)
			assert.Equal(b, "application/json", w.Header().Get("Content-Type"))
			assert.Equal(b, url+"/"+nfInstanceId, w.Header().Get("Location"))
			assert.Equal(b, nfInstanceId, response.NFInstanceId)
			assert.Equal(b, nfType, response.NFType)
			assert.Equal(b, nfStatus, response.NFStatus)
			// http request NFProfileRetrieve
			w = httptest.NewRecorder()
			request, err = http.NewRequest("GET", url+"/"+nfInstanceId, nil)
			if err != nil {
				b.Errorf("Error creating request: %v", err)
			}
		}
	})
}

func TestHandleNFRegisterSharedData(t *testing.T) {
	// initialize NRF Service
	NRFService = New()
	err := NRFService.Init()
	if err != nil {
		t.Error(err)
	}
	// start http test service
	server, router := startTestServer()
	defer server.Close()
	// construct network function request content
	url := server.URL + "/nnrf-nfm/v1/shared-data"
	sharedDataId := uuid.New().String()
	nfInstanceId := uuid.New().String()
	nfType := "AMF"
	nfStatus := "REGISTERED"
	// assemble network function http request
	profile := NFProfile{
		NFInstanceId: nfInstanceId,
		NFType:       nfType,
		NFStatus:     nfStatus,
	}
	sharedData := SharedData{
		SharedDataId:      sharedDataId,
		SharedProfileData: profile,
	}
	body, err := json.Marshal(sharedData)
	if err != nil {
		t.Errorf("Error marshalling shared data: %v", err)
	}
	// http request NFRegister
	w := httptest.NewRecorder()
	request, err := http.NewRequest("PUT", url+"/"+sharedDataId, bytes.NewReader(body))
	if err != nil {
		t.Errorf("Error creating request: %v", err)
	}
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, request)
	var response SharedData
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error unmarshalling response: %v", err)
	}
	// assert http response
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, url+"/"+sharedDataId, w.Header().Get("Location"))
	assert.Equal(t, sharedDataId, response.SharedDataId)
	assert.Equal(t, nfInstanceId, response.SharedProfileData.NFInstanceId)
	assert.Equal(t, nfType, response.SharedProfileData.NFType)
	assert.Equal(t, nfStatus, response.SharedProfileData.NFStatus)
}

func TestHandleNFSharedDataCompleteReplacement(t *testing.T) {
	// initialize NRF Service
	NRFService = New()
	err := NRFService.Init()
	if err != nil {
		t.Error(err)
	}
	// start http test service
	server, router := startTestServer()
	defer server.Close()
	// construct network function request content
	url := server.URL + "/nnrf-nfm/v1/shared-data"
	sharedDataId := uuid.New().String()
	nfInstanceId := uuid.New().String()
	nfType := "AMF"
	nfStatus := "REGISTERED"
	// assemble network function http request
	profile := NFProfile{
		NFInstanceId: nfInstanceId,
		NFType:       nfType,
		NFStatus:     nfStatus,
	}
	sharedData := SharedData{
		SharedDataId:      sharedDataId,
		SharedProfileData: profile,
	}
	body, err := json.Marshal(sharedData)
	if err != nil {
		t.Errorf("Error marshalling profile: %v", err)
	}
	// http request NFRegister
	w := httptest.NewRecorder()
	request, err := http.NewRequest("PUT", url+"/"+sharedDataId, bytes.NewReader(body))
	if err != nil {
		t.Errorf("Error creating request: %v", err)
	}
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, request)
	var response SharedData
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error unmarshalling response: %v", err)
	}
	// assert http response
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, url+"/"+sharedDataId, w.Header().Get("Location"))
	assert.Equal(t, sharedDataId, response.SharedDataId)
	assert.Equal(t, nfInstanceId, response.SharedProfileData.NFInstanceId)
	assert.Equal(t, nfType, response.SharedProfileData.NFType)
	assert.Equal(t, nfStatus, response.SharedProfileData.NFStatus)
	// construct network function request content
	nfStatusNew := "SUSPENDED"
	// assemble network function http request
	profileNew := NFProfile{
		NFInstanceId: nfInstanceId,
		NFType:       nfType,
		NFStatus:     nfStatusNew,
	}
	sharedDataNew := SharedData{
		SharedDataId:      sharedDataId,
		SharedProfileData: profileNew,
	}
	bodyNew, err := json.Marshal(sharedDataNew)
	if err != nil {
		t.Errorf("Error marshalling profile: %v", err)
	}
	// http request NFProfileCompleteReplacement
	w = httptest.NewRecorder()
	request, err = http.NewRequest("PUT", url+"/"+sharedDataId, bytes.NewReader(bodyNew))
	if err != nil {
		t.Errorf("Error creating request: %v", err)
	}
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, request)
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error unmarshalling response: %v", err)
	}
	// assert http response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, sharedDataId, response.SharedDataId)
	assert.Equal(t, nfInstanceId, response.SharedProfileData.NFInstanceId)
	assert.Equal(t, nfType, response.SharedProfileData.NFType)
	assert.Equal(t, nfStatusNew, response.SharedProfileData.NFStatus)
}

func TestHandleNFSharedDataRetrieve(t *testing.T) {
	// initialize NRF Service
	NRFService = New()
	err := NRFService.Init()
	if err != nil {
		t.Error(err)
	}
	// start http test service
	server, router := startTestServer()
	defer server.Close()
	// construct network function request content
	url := server.URL + "/nnrf-nfm/v1/shared-data"
	sharedDataId := uuid.New().String()
	nfInstanceId := uuid.New().String()
	nfType := "SMF"
	nfStatus := "REGISTERED"
	// assemble network function http request
	profile := NFProfile{
		NFInstanceId: nfInstanceId,
		NFType:       nfType,
		NFStatus:     nfStatus,
	}
	sharedData := SharedData{
		SharedDataId:      sharedDataId,
		SharedProfileData: profile,
	}
	body, err := json.Marshal(sharedData)
	if err != nil {
		t.Errorf("Error marshalling profile: %v", err)
	}
	w := httptest.NewRecorder()
	request, err := http.NewRequest("PUT", url+"/"+sharedDataId, bytes.NewReader(body))
	if err != nil {
		t.Errorf("Error creating request: %v", err)
	}
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, request)
	var response SharedData
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error unmarshalling response: %v", err)
	}
	// assert http response
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, url+"/"+sharedDataId, w.Header().Get("Location"))
	assert.Equal(t, sharedDataId, response.SharedDataId)
	assert.Equal(t, nfInstanceId, response.SharedProfileData.NFInstanceId)
	assert.Equal(t, nfType, response.SharedProfileData.NFType)
	assert.Equal(t, nfStatus, response.SharedProfileData.NFStatus)
	// http request NFProfileRetrieve
	w = httptest.NewRecorder()
	request, err = http.NewRequest("GET", url+"/"+sharedDataId, nil)
	if err != nil {
		t.Errorf("Error creating request: %v", err)
	}
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, request)
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error unmarshalling response: %v", err)
	}
	// assert http response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, sharedDataId, response.SharedDataId)
	assert.Equal(t, nfInstanceId, response.SharedProfileData.NFInstanceId)
	assert.Equal(t, nfType, response.SharedProfileData.NFType)
	assert.Equal(t, nfStatus, response.SharedProfileData.NFStatus)
}

func TestHandleNFListRetrieve(t *testing.T) {
	// initialize NRF Service
	NRFService = New()
	err := NRFService.Init()
	if err != nil {
		t.Error(err)
	}
	// start http test service
	server, router := startTestServer()
	defer server.Close()
	// construct network function request content
	url := server.URL + "/nnrf-nfm/v1/nf-instances"
	nfInstanceId := uuid.New().String()
	nfType := "SMF"
	nfStatus := "REGISTERED"
	// assemble network function http request
	profile := NFProfile{
		NFInstanceId: nfInstanceId,
		NFType:       nfType,
		NFStatus:     nfStatus,
	}
	body, err := json.Marshal(profile)
	if err != nil {
		t.Errorf("Error marshalling profile: %v", err)
	}
	w := httptest.NewRecorder()
	request, err := http.NewRequest("PUT", url+"/"+nfInstanceId, bytes.NewReader(body))
	if err != nil {
		t.Errorf("Error creating request: %v", err)
	}
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, request)
	var response NFProfile
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error unmarshalling response: %v", err)
	}
	// assert http response
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, url+"/"+nfInstanceId, w.Header().Get("Location"))
	assert.Equal(t, nfInstanceId, response.NFInstanceId)
	assert.Equal(t, nfType, response.NFType)
	assert.Equal(t, nfStatus, response.NFStatus)
	// http request NFListRetrieve
	w = httptest.NewRecorder()
	request, err = http.NewRequest("GET", url+"?nf-type=SMF&limit=1&page-number=1&page-size=1", nil)
	if err != nil {
		t.Errorf("Error creating request: %v", err)
	}
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, request)
	var responseNew UriList
	err = json.Unmarshal(w.Body.Bytes(), &responseNew)
	if err != nil {
		t.Errorf("Error unmarshalling response: %v", err)
	}
	// assert http response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/3gppHal+json", w.Header().Get("Content-Type"))
	assert.Equal(t, url+"/"+nfInstanceId, responseNew.Links[0])
	assert.Equal(t, 1, responseNew.TotalItemCount)
}
