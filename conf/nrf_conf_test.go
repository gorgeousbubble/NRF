package conf

import "testing"

func TestMarshalTo(t *testing.T) {
	conf := NRFConf{
		AcceptNFHeartBeatTimer: false,
		DefaultHeartBeatTimer:  60,
	}
	err := MarshalTo("./nrf_conf.yaml", conf)
	if err != nil {
		t.Error("Marshal to file failed:", err)
	}
}

func TestUnmarshalFromFile(t *testing.T) {
	var conf NRFConf
	err := UnmarshalFrom("./nrf_conf.yaml", &conf)
	if err != nil {
		t.Error("Unmarshal from file failed:", err)
	}
}
