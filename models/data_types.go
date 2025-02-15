package models

// AccessTokenReq Access Token Request
type AccessTokenReq struct {
	GrantType             string            `json:"grant_type" yaml:"grant_type"`                       // mandatory, 1, enum: client_credentials
	NfInstanceId          string            `json:"nfInstanceId" yaml:"nfInstanceId"`                   // mandatory, 1, string: uuid
	NFType                string            `json:"nfType" yaml:"nfType"`                               // conditional, 0..1, enum: eg. AMF, SMF
	TargetNfType          string            `json:"targetNfType" yaml:"targetNfType"`                   // conditional, 0..1, enum: eg. AMF, SMF
	Scope                 string            `json:"scope" yaml:"scope"`                                 // mandatory, 1, string:
	TargetNfInstanceId    string            `json:"targetNfInstanceId" yaml:"targetNfInstanceId"`       // conditional, 0..1, string: uuid
	RequesterPlmn         PlmnId            `json:"requesterPlmn" yaml:"requesterPlmn"`                 // conditional, 0..1,
	RequesterPlmnList     []PlmnId          `json:"requesterPlmnList" yaml:"requesterPlmnList"`         // conditional, 2..N,
	RequesterSnssaiList   []Snssai          `json:"requesterSnssaiList" yaml:"requesterSnssaiList"`     // optional, 1..N,
	RequesterFqdn         string            `json:"requesterFqdn" yaml:"requesterFqdn"`                 // optional, 1..N,
	RequesterSnpnList     []PlmnIdNid       `json:"requesterSnpnList" yaml:"requesterSnpnList"`         // optional, 1..N,
	TargetPlmn            PlmnId            `json:"targetPlmn" yaml:"targetPlmn"`                       // conditional, 0..1,
	TargetSnpn            PlmnIdNid         `json:"targetSnpn" yaml:"targetSnpn"`                       // conditional, 0..1,
	TargetSnssaiList      []Snssai          `json:"targetSnssaiList" yaml:"targetSnssaiList"`           // optional, 1..N,
	TargetNsiList         []string          `json:"targetNsiList" yaml:"targetNsiList"`                 // optional, 1..N,
	TargetNfSetId         string            `json:"targetNfSetId" yaml:"targetNfSetId"`                 // optional, 1..N,
	TargetNfServiceSetId  string            `json:"targetNfServiceSetId" yaml:"targetNfServiceSetId"`   // optional, 1..N,
	HnrfAccessTokenUri    string            `json:"hnrfAccessTokenUri" yaml:"hnrfAccessTokenUri"`       // conditional, 0..1,
	SourceNfInstanceId    string            `json:"sourceNfInstanceId" yaml:"sourceNfInstanceId"`       // conditional, 0..1,
	VendorId              string            `json:"vendorId" yaml:"vendorId"`                           // conditional, 0..1,
	AnalyticsIds          []string          `json:"analyticsIds" yaml:"analyticsIds"`                   // conditional, 0..N,
	RequesterInterIndList []MlModelInterInd `json:"requesterInterIndList" yaml:"requesterInterIndList"` // conditional, 0..N,
	SourceVendorId        string            `json:"sourceVendorId" yaml:"sourceVendorId"`               // conditional, 0..1,
}

// AccessTokenRsp Access Token Response
type AccessTokenRsp struct {
	AccessToken string `json:"access_token" yaml:"access_token"` // mandatory, 1,
	TokenType   string `json:"token_type" yaml:"token_type"`     // mandatory, 1,
	ExpiresIn   int64  `json:"expires_in" yaml:"expires_in"`     // conditional, 0..1,
	Scope       string `json:"scope" yaml:"scope"`               // conditional, 0..1,
}

// AccessTokenClaims Access Token Claims
type AccessTokenClaims struct {
	Iss                              string              `json:"iss" yaml:"iss"`                                                           // mandatory, 1, string: uuid
	Sub                              string              `json:"sub" yaml:"sub"`                                                           // mandatory, 1, string: uuid
	Aud                              string              `json:"aud" yaml:"aud"`                                                           // mandatory, 1, string: uuid
	Scope                            string              `json:"scope" yaml:"scope"`                                                       // mandatory, 1, string:
	Exp                              int64               `json:"exp" yaml:"exp"`                                                           // mandatory, 1,
	ConsumerPlmnId                   PlmnId              `json:"consumerPlmnId" yaml:"consumerPlmnId"`                                     // conditional, 0..1,
	ConsumerSnpnId                   PlmnIdNid           `json:"consumerSnpnId" yaml:"consumerSnpnId"`                                     // conditional, 0..1,
	ProducerPlmnId                   PlmnId              `json:"producerPlmnId" yaml:"producerPlmnId"`                                     // conditional, 0..1,
	ProducerSnpnId                   PlmnIdNid           `json:"producerSnpnId" yaml:"producerSnpnId"`                                     // conditional, 0..1,
	ProducerSnssaiList               []Snssai            `json:"producerSnssaiList" yaml:"producerSnssaiList"`                             // optional, 1..N,
	ProducerNsiList                  []string            `json:"producerNsiList" yaml:"producerNsiList"`                                   // optional, 1..N,
	ProducerNfSetId                  string              `json:"producerNfSetId" yaml:"producerNfSetId"`                                   // optional, 1..N,
	ProducerNfServiceSetId           string              `json:"producerNfServiceSetId" yaml:"producerNfServiceSetId"`                     // optional, 1..N,
	SourceNfInstanceId               string              `json:"sourceNfInstanceId" yaml:"sourceNfInstanceId"`                             // conditional, 0..1,
	AnalyticsIdList                  []string            `json:"analyticsIds" yaml:"analyticsIds"`                                         // conditional, 1..N,
	ApplicableResourceContentFilters []string            `json:"applicableResourceContentFilters" yaml:"applicableResourceContentFilters"` // conditional, 1..N,
	ConsumerRcfInfo                  map[string][]string `json:"consumerRcfInfo" yaml:"consumerRcfInfo"`                                   // optional, 1..N,
	Iat                              int64               `json:"iat" yaml:"iat"`                                                           // optional, 0..1,
}

// AccessTokenErr Access Token Error
type AccessTokenErr struct {
	Error            string `json:"error" yaml:"error"`                         // mandatory, 1,
	ErrorDescription string `json:"error_description" yaml:"error_description"` // optional, 1
	ErrorUri         string `json:"error_uri" yaml:"error_uri"`                 // optional, 1
}

type PlmnId struct {
	Mcc string `json:"mcc" yaml:"mcc"` // mandatory, 1, string: Mobile Country Code (MCC) (3 digits), '^\d{3}$'
	Mnc string `json:"mnc" yaml:"mnc"` // mandatory, 1, string: Mobile Network Code (MNC) (2 or 3 digits), '^\d{2,3}$'
}

type Snssai struct {
	Sst int32  `json:"sst" yaml:"sst"` // mandatory, 1, int32: SST (0-255)
	Sd  string `json:"sd" yaml:"sd"`   // mandatory, 1, string: SD (3-octet string)
}

type PlmnIdNid struct {
	Mcc string `json:"mcc" yaml:"mcc"` // mandatory, 1, string: Mobile Country Code (MCC) (3 digits), '^\d{3}$'
	Mnc string `json:"mnc" yaml:"mnc"` // mandatory, 1, string: Mobile Network Code (MNC) (2 or 3 digits), '^\d{2,3}$'
	Nid string `json:"nid" yaml:"nid"` // mandatory, 1, string, Network Identifier (NID), '^[A-Fa-f0-9]{11}$'
}

type MlModelInterInd struct {
	AnalyticsId string   `json:"analyticsId" yaml:"analyticsId"` // mandatory, 1, string: Analytics Id
	VendorList  []string `json:"vendorList" yaml:"vendorList"`   // mandatory, 1, []string: NWDAF vendors
}
