package gol

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func isNil(value interface{}) bool {
	val := reflect.ValueOf(value)
	kind := val.Kind()
	if kind == reflect.Ptr {
		val = val.Elem()
		kind = val.Kind()
	}

	return !val.IsValid()
}

func getAddr(value reflect.Value) (string, error) {
	val := value
	kind := value.Kind()

	for kind == reflect.Array || kind == reflect.Slice || kind == reflect.Struct || kind == reflect.Interface || kind == reflect.Ptr {
		if kind == reflect.Ptr {
			elem := val.Elem()
			if !elem.IsValid() {
				break
			}
		}

		switch kind {
		case reflect.Array, reflect.Slice:
			if val.Len() < 1 {
				return "", errors.New("array or slice not value")
			}
			val = val.Index(0)
			kind = val.Kind()
		case reflect.Struct:
			if val.NumField() < 1 {
				return "", errors.New("struct not value")
			}
			val = val.Field(0)
			kind = val.Kind()
		case reflect.Ptr:
			val = val.Elem()
			kind = val.Kind()
		default:
			val = val.Elem()
			kind = val.Kind()
		}
	}

	addr := fmt.Sprintf("%x", val.Addr())

	return addr, nil
}

func getAddrFromInterface(value interface{}) (string, error) {
	val := reflect.ValueOf(value)

	addr, err := getAddr(val)
	if err != nil {
		return "", err
	}

	return addr, nil
}

func makeTagIndexMap(value reflect.Type, tagName string) (map[string][]int, error) {
	tagIndexMap := make(map[string][]int, 0)

	if value.Kind() != reflect.Struct {
		return tagIndexMap, errors.New("not struct")
	}

	tagIndexMap = makeTagIndexMapRe(tagIndexMap, []int{}, value, tagName)

	return tagIndexMap, nil
}

func makeTagIndexMapRe(tagIndexMap map[string][]int, indexList []int, value reflect.Type, tagName string) map[string][]int {
	if value.Kind() != reflect.Struct {
		return tagIndexMap
	}

	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		tag := field.Tag.Get(tagName)
		indexNextList := append(indexList, i)
		if tag == "" {
			tagIndexMap = makeTagIndexMapRe(tagIndexMap, indexNextList, field.Type, tagName)
		} else {
			tagIndexMap[tag] = indexNextList
		}
	}

	return tagIndexMap
}

func toCamelCase(value string) string {
	str := ""

	valueList := strings.Split(value, "")
	limit := len(valueList)
	var flag bool
	for i := 0; i < limit; i++ {
		val := valueList[i]
		if val == "_" {
			flag = true
			continue
		}

		if flag {
			val = strings.ToUpper(val)
		}
		flag = false

		str = fmt.Sprintf("%s%s", str, val)
	}

	return str
}

func toSnakeCase(value string) string {
	str := ""
	match := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	valueList := strings.Split(value, "")
	limit := len(valueList)
	for i := 0; i < limit; i++ {
		val := valueList[i]
		if strings.Contains(match, val) {
			val = strings.ToLower(val)
			if i != 0 {
				val = fmt.Sprintf("%s%s", "_", val)
			}
		}
		str = fmt.Sprintf("%s%s", str, val)
	}

	return str
}
