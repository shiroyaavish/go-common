package utils

import (
	"encoding/json"
	"io"
)

// UnmarshalJSON unmarshalls JSON data from an io.ReadCloser into the provided interface{} value.
func UnmarshalJSON(reader io.ReadCloser, v interface{}) error {
	dataBytes, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	return json.Unmarshal(dataBytes, v)
}

func ConvertJSON(src interface{}, dst interface{}) error {
	dataBytes, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(dataBytes, dst)
}
