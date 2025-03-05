package data

type NFProfile struct {
	NFInstanceId   string `json:"nfInstanceId" yaml:"nfInstanceId"`
	NFType         string `json:"nfType" yaml:"nfType"`
	NFStatus       string `json:"nfStatus" yaml:"nfStatus"`
	NFInstanceName string `json:"nfInstanceName" yaml:"nfInstanceName"`
	HeartBeatTimer int    `json:"heartBeatTimer" yaml:"heartBeatTimer"`
}
