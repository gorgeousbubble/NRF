package app

import "sync"

type NRF struct {
	profiles map[string][]NFProfile
	mutex    sync.Mutex
}

type NFProfile struct {
	NFInstanceId string `json:"nfInstanceId" yaml:"nfInstanceId"`
	NFType       string `json:"nfType" yaml:"nfType"`
	NFStatus     string `json:"nfStatus" yaml:"nfStatus"`
}
