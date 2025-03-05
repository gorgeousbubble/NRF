package conf

var NRFConfigure NRFConf

func LoadConf() (err error) {
	return UnmarshalFrom("./conf/nrf_conf.yaml", &NRFConfigure)
}
