package data

type NFProfile struct {
	NFInstanceId   string      `json:"nfInstanceId" yaml:"nfInstanceId"`
	NFType         string      `json:"nfType" yaml:"nfType"`
	NFStatus       string      `json:"nfStatus" yaml:"nfStatus"`
	NFInstanceName string      `json:"nfInstanceName" yaml:"nfInstanceName"`
	HeartBeatTimer int         `json:"heartBeatTimer" yaml:"heartBeatTimer"`
	NFServices     []NFService `json:"nfServices" yaml:"nfServices"`
}

type NFService struct {
	ServiceInstanceId string `json:"serviceInstanceId" yaml:"serviceInstanceId"`
	SupportedFeatures string `json:"supportedFeatures" yaml:"supportedFeatures"`
}

type NFProfileRegistrationError struct {
	ProblemDetails   ProblemDetails   `json:"problemDetails" yaml:"problemDetails"`
	SharedDataIdList SharedDataIdList `json:"sharedDataIdList" yaml:"sharedDataIdList"`
}

type ProblemDetails struct {
	Type     string `json:"type" yaml:"type"`
	Title    string `json:"title" yaml:"title"`
	Status   int    `json:"status" yaml:"status"`
	Detail   string `json:"detail" yaml:"detail"`
	Instance string `json:"instance" yaml:"instance"`
	Cause    string `json:"cause" yaml:"cause"`
}

type SharedDataIdList struct {
	SharedDataIds []string `json:"sharedDataIds" yaml:"sharedDataIds"`
}
