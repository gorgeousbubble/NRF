package util

import (
	"errors"
	"github.com/google/uuid"
	. "nrf/conf"
	"strings"
)

func CheckNFInstanceId(nfInstanceId string) (b bool, err error) {
	b, err = true, nil
	// parse NFInstanceId
	_, err = uuid.Parse(nfInstanceId)
	if err != nil {
		b = false
		return b, err
	}
	return b, err
}

func MarshalNFInstanceId(nfInstanceId *string) (err error) {
	err = nil
	// marshal NFInstanceId
	*nfInstanceId = strings.ToLower(*nfInstanceId)
	return err
}

func CheckNFType(nfType string) (b bool, err error) {
	b, err = true, nil
	// check NFType
	switch nfType {
	case "NRF", "UDM", "AMF", "SMF", "AUSF", "NEF", "PCF", "SMSF", "NSSF", "UDR", "LMF", "GMLC", "5G_EIR", "SEPP":
	case "UPF", "N3IWF", "AF", "UDSF", "BSF", "CHF", "NWDAF", "PCSCF", "CBCF", "UCMF", "HSS", "SOR_AF", "SPAF", "MME":
	case "SCSAS", "SCEF", "SCP", "NSSAAF", "ICSCF", "SCSCF", "DRA", "IMS_AS", "AANF", "5G_DDNMF", "NSACF", "MFAF":
	case "EASDF", "DCCF", "MB_SMF", "TSCTSF", "ADRF", "GBA_BSF", "CEF", "MB_UPF", "NSWOF", "PKMF", "MNPF", "SMS_GMSC":
	case "SMS_IWMSC", "MBSF", "MBSTF", "PANF", "IP_SM_GW", "SMS_ROUTER", "DCSF", "MRF", "MRFP", "MF", "SLPKMF", "RH":
	default:
		b, err = false, errors.New("NFType is invalid")
		return b, err
	}
	return b, err
}

func CheckNFStatus(nfStatus string) (b bool, err error) {
	b, err = true, nil
	// check NFStatus
	switch nfStatus {
	case "REGISTERED":
	case "SUSPENDED":
	case "UNDISCOVERABLE":
	case "CANARY_RELEASE":
	default:
		b, err = false, errors.New("NFStatus is invalid")
		return b, err
	}
	return b, err
}

func CheckHeartBeatTimer(heartBeatTimer int) (b bool, err error) {
	b, err = true, nil
	// check HeartBeatTimer
	if heartBeatTimer < 1 {
		b, err = false, errors.New("HeartBeatTimer should be greater than 1")
		return b, err
	} else if heartBeatTimer > 3600 {
		b, err = false, errors.New("HeartBeatTimer should be less than 3600")
		return b, err
	}
	return b, err
}

func MarshalHeartBeatTimer(heartBeatTimer *int) (err error) {
	err = nil
	// marshal HeartBeatTimer
	if !NRFConfigure.AcceptNFHeartBeatTimer {
		*heartBeatTimer = NRFConfigure.DefaultHeartBeatTimer
	}
	return err
}
