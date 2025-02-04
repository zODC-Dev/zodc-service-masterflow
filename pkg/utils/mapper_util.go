package utils

import (
	"encoding/json"
)

func Mapper(source interface{}, dest interface{}) error {
	data, err := json.Marshal(source)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, dest); err != nil {
		_, ok := err.(*json.UnmarshalTypeError)
		if ok {
			return nil
		}

		return err
	}

	return nil
}

func MapToString(value any) (string, error) {
	bytes, err := json.Marshal(value)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
