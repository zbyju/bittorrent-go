package mybencode

import (
	"errors"
	"strconv"
	"unicode"
)

// Example:
// - 5:hello -> hello
// - 10:hello12345 -> hello12345
func DecodeBencode(bencodedString string) (interface{}, error) {
	res, err := Decode(bencodedString, 0)

	if err != nil {
		return "", err
	}

	return res.result, nil
}

type decodeResult struct {
	result interface{}
	index  int
}

func Decode(bencodedString string, index int) (decodeResult, error) {
	var res decodeResult
	var err error

	if unicode.IsDigit(rune(bencodedString[index])) {
		res, err = DecodeString(bencodedString, index)
	} else if bencodedString[index] == 'i' {
		res, err = DecodeInt(bencodedString, index)
	} else if bencodedString[index] == 'l' {
		res, err = DecodeList(bencodedString, index)
	} else if bencodedString[index] == 'd' {
		res, err = DecodeDict(bencodedString, index)
	} else {
		return decodeResult{}, errors.New("Input string needs to be in format of String, Int, Array, Dict. Found: " + string(bencodedString[index]))
	}

	if err != nil {
		return decodeResult{}, nil
	}

	return res, nil
}

func decodeRawString(bencodedString string, index int) (string, int, error) {
	var firstColonIndex int

	for i := index; i < len(bencodedString); i++ {
		if bencodedString[i] == ':' {
			firstColonIndex = i
			break
		}
	}

	lengthStr := bencodedString[index:firstColonIndex]

	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		return "", 0, err
	}

	res := bencodedString[firstColonIndex+1 : firstColonIndex+1+length]

	return res, firstColonIndex + length + 1, nil
}

func DecodeString(bencodedString string, index int) (decodeResult, error) {
	res, index, err := decodeRawString(bencodedString, index)

	if err != nil {
		return decodeResult{}, err
	}

	return decodeResult{res, index}, nil
}

func DecodeInt(bencodedString string, index int) (decodeResult, error) {
	var firstEIndex int

	for i := index; i < len(bencodedString); i++ {
		if bencodedString[i] == 'e' {
			firstEIndex = i
			break
		}
	}

	numStr := bencodedString[index+1 : firstEIndex]

	num, err := strconv.Atoi(numStr)
	if err != nil {
		return decodeResult{}, err
	}

	return decodeResult{num, firstEIndex + 1}, nil
}

func DecodeList(bencodedString string, index int) (decodeResult, error) {
	list := []interface{}{}
	index += 1

	for bencodedString[index] != 'e' {
		res, err := Decode(bencodedString, index)

		if err != nil {
			return decodeResult{}, err
		}

		list = append(list, res.result)
		index = res.index
	}

	return decodeResult{list, index + 1}, nil
}

func DecodeDict(bencodedString string, index int) (decodeResult, error) {
	dict := make(map[string]interface{})
	index += 1

	for bencodedString[index] != 'e' {
		key, i, err := decodeRawString(bencodedString, index)

		if err != nil {
			return decodeResult{}, err
		}

		index = i

		val, err := Decode(bencodedString, index)

		if err != nil {
			return decodeResult{}, err
		}

		dict[key] = val.result
		index = val.index
	}

	return decodeResult{dict, index + 1}, nil
}
