package util

import (
	"github.com/google/uuid"
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
