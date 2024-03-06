package utils

import "github.com/bytedance/sonic"

func ToJson(val interface{}) ([]byte, error) {
	return sonic.Marshal(val)
}

func FromJson(d []byte, val interface{}) error {
	return sonic.Unmarshal(d, val)
}
