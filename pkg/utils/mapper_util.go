package utils

import (
	"encoding/json"
	"fmt"
	"strconv"
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

func StringToInt32(value string) (int32, error) {
	if value == "" {
		return 0, fmt.Errorf("value is empty")
	}
	valueInt, _ := strconv.Atoi(value)
	return int32(valueInt), nil
}
