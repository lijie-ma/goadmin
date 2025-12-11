package setting

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func encoding(value any) (string, error) {
	str, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return string(str), nil
}

func decoding(value string, resultPtr any) error {
	if reflect.ValueOf(resultPtr).Kind() != reflect.Ptr {
		return fmt.Errorf("resultPtr must be a pointer")
	}
	return json.Unmarshal([]byte(value), resultPtr)
}
