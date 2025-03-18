package conf

import (
	"gopkg.in/yaml.v3"
	"io"
	"os"
)

type NRFConf struct {
	SBITLSSettings         string `json:"sbiTLSSettings" yaml:"sbiTLSSettings"`
	AcceptNFHeartBeatTimer bool   `json:"acceptNFHeartBeatTimer" yaml:"acceptNFHeartBeatTimer"`
	DefaultHeartBeatTimer  int    `json:"defaultHeartBeatTimer" yaml:"defaultHeartBeatTimer"`
}

func MarshalTo(file string, t interface{}) (err error) {
	return marshalTo(file, t)
}

func UnmarshalFrom(file string, t interface{}) (err error) {
	return unmarshalFrom(file, t)
}

func marshalTo(in string, t interface{}) (err error) {
	// try to open file...
	file, err := os.Open(in)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err = file.Close()
	}(file)
	// marshal
	out, err := yaml.Marshal(t)
	if err != nil {
		return err
	}
	// write to the file...
	err = os.WriteFile(in, out, 0644)
	if err != nil {
		return err
	}
	return err
}

func unmarshalFrom(in string, t interface{}) (err error) {
	// try to open file...
	file, err := os.Open(in)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err = file.Close()
	}(file)
	// read the file...
	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	// unmarshal
	err = yaml.Unmarshal(data, t)
	if err != nil {
		return err
	}
	return err
}
