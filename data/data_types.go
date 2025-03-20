package data

type NFProfile struct {
	NFInstanceId   string      `json:"nfInstanceId" yaml:"nfInstanceId" binding:"required"`
	NFType         string      `json:"nfType" yaml:"nfType" binding:"required"`
	NFStatus       string      `json:"nfStatus" yaml:"nfStatus" binding:"required"`
	NFInstanceName string      `json:"nfInstanceName" yaml:"nfInstanceName" binding:"omitempty"`
	HeartBeatTimer int         `json:"heartBeatTimer" yaml:"heartBeatTimer" binding:"omitempty"`
	NFServices     []NFService `json:"nfServices" yaml:"nfServices" binding:"omitempty"`
}

type NFService struct {
	ServiceInstanceId string `json:"serviceInstanceId" yaml:"serviceInstanceId" binding:"required"`
	SupportedFeatures string `json:"supportedFeatures" yaml:"supportedFeatures" binding:"omitempty"`
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

type SharedData struct {
	SharedDataId      string    `json:"sharedDataId" yaml:"sharedDataId" binding:"required"`
	SharedProfileData NFProfile `json:"sharedProfileData" yaml:"sharedProfileData" binding:"omitempty"`
	SharedServiceData NFService `json:"sharedServiceData" yaml:"sharedServiceData" binding:"omitempty"`
}
