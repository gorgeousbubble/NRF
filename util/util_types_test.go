package util

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestCheckNFInstanceId(t *testing.T) {
	nfInstanceId := uuid.New().String()
	b, err := CheckNFInstanceId(nfInstanceId)
	if b != true || err != nil {
		t.Fatal("Error Check NFInstanceId:", err)
	}
}

func BenchmarkCheckNFInstanceId(b *testing.B) {
	for i := 0; i < b.N; i++ {
		nfInstanceId := uuid.New().String()
		r, err := CheckNFInstanceId(nfInstanceId)
		if r != true || err != nil {
			b.Fatal("Error Check NFInstanceId:", err)
		}
	}
}

func TestMarshalNFInstanceId(t *testing.T) {
	nfInstanceId := uuid.New().String()
	marshalInstanceId := strings.ToLower(nfInstanceId)
	nfInstanceId = strings.ToUpper(nfInstanceId)
	err := MarshalNFInstanceId(&nfInstanceId)
	if err != nil {
		t.Fatal("Error Marshal NFInstanceId:", err)
	}
	assert.Equal(t, nfInstanceId, marshalInstanceId)
}

func BenchmarkMarshalNFInstanceId(b *testing.B) {
	for i := 0; i < b.N; i++ {
		nfInstanceId := uuid.New().String()
		marshalInstanceId := strings.ToLower(nfInstanceId)
		nfInstanceId = strings.ToUpper(nfInstanceId)
		err := MarshalNFInstanceId(&nfInstanceId)
		if err != nil {
			b.Fatal("Error Marshal NFInstanceId:", err)
		}
		assert.Equal(b, nfInstanceId, marshalInstanceId)
	}
}

func TestCheckNFType(t *testing.T) {
	nfType := "AMF"
	b, err := CheckNFType(nfType)
	if b != true || err != nil {
		t.Fatal("Error Check NFType:", err)
	}
}

func BenchmarkCheckNFType(b *testing.B) {
	for i := 0; i < b.N; i++ {
		nfType := "AMF"
		r, err := CheckNFType(nfType)
		if r != true || err != nil {
			b.Fatal("Error Check NFType:", err)
		}
	}
}

func TestCheckNFStatus(t *testing.T) {
	nfStatus := "REGISTERED"
	b, err := CheckNFStatus(nfStatus)
	if b != true || err != nil {
		t.Fatal("Error Check NFStatus:", err)
	}
}

func BenchmarkCheckNFStatus(b *testing.B) {
	for i := 0; i < b.N; i++ {
		nfStatus := "REGISTERED"
		r, err := CheckNFStatus(nfStatus)
		if r != true || err != nil {
			b.Fatal("Error Check NFStatus:", err)
		}
	}
}

func TestCheckHeartBeatTimer(t *testing.T) {
	heartBeatTimer := 60
	b, err := CheckHeartBeatTimer(heartBeatTimer)
	if b != true || err != nil {
		t.Fatal("Error Check HeartBeatTimer:", err)
	}
}

func BenchmarkCheckHeartBeatTimer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		heartBeatTimer := 60
		r, err := CheckHeartBeatTimer(heartBeatTimer)
		if r != true || err != nil {
			b.Fatal("Error Check HeartBeatTimer:", err)
		}
	}
}
