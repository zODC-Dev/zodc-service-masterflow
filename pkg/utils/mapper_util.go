package utils

import (
	"encoding/json"
	"fmt"
	"log/slog"
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

func PtrStringToPtrInt32(value *string) (*int32, error) {
	slog.Info("value", "value", value)
	if value == nil {
		return nil, fmt.Errorf("value is empty")
	}
	valueInt, _ := strconv.Atoi(*value)
	result := int32(valueInt)
	return &result, nil
}

func StringToFloat64(value string) *float64 {
	if value == "" {
		return nil
	}
	valueFloat, _ := strconv.ParseFloat(value, 64)
	return &valueFloat
}
