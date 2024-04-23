package util

import (
	"bytes"
	"encoding/json"
)

func MarshalUnescape(v interface{}) (string, error) {
	bf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(bf)
	jsonEncoder.SetEscapeHTML(false)
	err := jsonEncoder.Encode(v)
	if err != nil {
		return "", err
	}
	return bf.String(), nil
}

func MarshalIndentUnescape(v interface{}, prefix, indent string) (string, error) {
	bf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(bf)
	jsonEncoder.SetEscapeHTML(false)
	err := jsonEncoder.Encode(v)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = json.Indent(&buf, bf.Bytes(), prefix, indent)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
