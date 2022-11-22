package utils

import (
	"encoding/json"
	"math/rand"
	"strings"
	"time"
	"unsafe"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var (
	chars    = "0123456789abcdefghijklmnopqrstuvwxyz"
	charsLen = len(chars)
	rng      = rand.NewSource(time.Now().UnixNano())
	mask     = int64(1<<6 - 1)
)

func RandomString(n int) string {
	buf := make([]byte, n)
	for idx, cache, remain := n, rng.Int63(), 10; idx > 0; {
		if remain == 0 {
			cache, remain = rng.Int63(), 10
		}
		buf[idx-1] = chars[int(cache&mask)%charsLen]
		cache >>= 6
		remain--
		idx--
	}
	return *(*string)(unsafe.Pointer(&buf))
}

// JsonToMap Convert json string to map
func JsonToMap(jsonStr *string) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(*jsonStr), &m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

// MapToJson Convert map json string
func MapToJson(m *map[string]interface{}) (string, error) {
	jsonByte, err := json.Marshal(m)
	if err != nil {
		return "", err
	}

	return string(jsonByte), nil
}

func MaskJsonStr(jsonStr *string, fieldNames []string) string {
	if !json.Valid([]byte(*jsonStr)) {
		return *jsonStr
	}

	jm, err := JsonToMap(jsonStr)
	if err != nil {
		return *jsonStr
	}

	for _, key := range fieldNames {
		MaskField(jm, key)
	}

	maskedJson, err := MapToJson(&jm)
	if err != nil {
		return *jsonStr
	}

	return maskedJson
}

// MaskStruct mask struct object then return string
func MaskStruct(src interface{}, fieldNames []string) string {
	jsonBytes, err := json.Marshal(src)
	if err != nil {
		return ""
	}

	jsonStr := string(jsonBytes)
	jm, err := JsonToMap(&jsonStr)
	if err != nil {
		return jsonStr
	}

	for _, key := range fieldNames {
		MaskField(jm, key)
	}

	maskedJson, err := MapToJson(&jm)
	if err != nil {
		return jsonStr
	}

	return maskedJson
}

func MaskField(jm map[string]interface{}, field string) {
	for k, v := range jm {
		switch vv := v.(type) {
		case bool, string, float64, int, []interface{}:
			if k == field || strings.Contains(k, field) {
				jm[k] = nil
			}
		case map[string]interface{}:
			MaskField(vv, field)
		case nil:
		default:
		}
	}
}

func MaskHttpHeader(jm map[string][]string, fieldNames []string) map[string][]string {
	for _, key := range fieldNames {
		for k, _ := range jm {
			if k == key || strings.Contains(k, key) {
				jm[k][0] = "***"
			}
		}
	}

	return jm
}
