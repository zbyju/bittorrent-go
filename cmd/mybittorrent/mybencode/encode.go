package mybencode

import (
	"fmt"
	"sort"
	"strconv"
)

func EncodeBencode(input interface{}) (string, error) {
	res, err := Encode(input)

	if err != nil {
		return "", err
	}

	return res, nil
}

func Encode(input interface{}) (res string, err error) {
	switch v := input.(type) {
	case int:
		return EncodeInt(v)
	case string:
		return EncodeString(v)
	case []interface{}:
		return EncodeList(v)
	case map[string]interface{}:
		return EncodeDict(v)
	default:
		return "", fmt.Errorf("unsupported type: %T", v)
	}
}

func EncodeString(in string) (string, error) {
	return strconv.Itoa(len(in)) + ":" + in, nil
}

func EncodeInt(in int) (string, error) {
	return "i" + strconv.Itoa(in) + "e", nil
}

func EncodeList(in []interface{}) (string, error) {
	res := "l"

	for x := range in {
		str, err := Encode(x)

		if err != nil {
			return "", err
		}

		res += str
	}

	res += "e"
	return res, nil
}

func EncodeDict(in map[string]interface{}) (string, error) {
	res := "d"

	var sortedKeys []string = make([]string, 0)
	for key := range in {
		sortedKeys = append(sortedKeys, key)
	}
	sort.Strings(sortedKeys)

	for _, k := range sortedKeys {
		v := in[k]
		ks, err := Encode(k)

		if err != nil {
			return "", err
		}

		res += ks

		vs, err := Encode(v)

		if err != nil {
			return "", err
		}

		res += vs
	}

	res += "e"
	return res, nil
}
