package utils

import (
	"fmt"
	"reflect"
)

func MapData(dest interface{}, src interface{}) error {
	destValue := reflect.ValueOf(dest)
	srcValue := reflect.ValueOf(src)

	if destValue.Kind() != reflect.Ptr {
		return fmt.Errorf("destination must be a pointer")
	}

	destValue = destValue.Elem()

	if srcValue.Kind() == reflect.Ptr {
		srcValue = srcValue.Elem()
	}

	if destValue.Kind() != reflect.Struct || srcValue.Kind() != reflect.Struct {
		return fmt.Errorf("both source and destination must be structs")
	}

	destType := destValue.Type()

	for i := 0; i < destValue.NumField(); i++ {
		destField := destValue.Field(i)
		destFieldName := destType.Field(i).Name

		srcField := srcValue.FieldByName(destFieldName)

		if srcField.IsValid() && destField.CanSet() {
			if destField.Type() == srcField.Type() {
				destField.Set(srcField)
			} else {
				switch {
				case destField.Kind() == reflect.String && srcField.Kind() == reflect.Int:
					destField.SetString(fmt.Sprintf("%d", srcField.Int()))
				case destField.Kind() == reflect.Int && srcField.Kind() == reflect.String:
					continue
				default:
					continue
				}
			}
		}
	}

	return nil
}
