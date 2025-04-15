package helpers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
)

func JsonToMap(jsonStr string) map[string]interface{} {
	jsonMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonStr), &jsonMap)
	if err != nil {
		return nil
	}
	return jsonMap
}

func InterfaceToMap(interfaceObj interface{}) map[string]interface{} {
	marshal, _ := json.Marshal(interfaceObj)

	jsonMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(string(marshal)), &jsonMap)
	if err != nil {
		return nil
	}
	return jsonMap
}

func InterfaceToInt(interfaceObj interface{}) (integer int) {
	marshal, _ := json.Marshal(interfaceObj)

	err := json.Unmarshal([]byte(string(marshal)), &integer)
	if err != nil {
		return
	}
	return
}
func MsgForTag(tag string) string {
	switch tag {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email"
	case "min":
		return "Invalid item must be longer"
	case "max":
		return "Invalid item cannot be longer"
	case "len":
		return "Invalid length"
	}
	return "unknown error"
}

func StructToBodyParams(structure interface{}) map[string]interface{} {

	requestValue := InterfaceToMap(structure)
	var bodyParams = make(map[string]interface{})

	for key, val := range requestValue {
		if val == nil {
			continue
		}
		bodyParams[key] = val

	}

	return bodyParams
}

func ExtractValue(str, prefix string) string {
	startIndex := strings.Index(str, prefix)
	if startIndex == -1 {
		return ""
	}

	startIndex += len(prefix)
	endIndex := strings.Index(str[startIndex:], " ")
	if endIndex == -1 {
		endIndex = len(str)
	} else {
		endIndex += startIndex
	}

	return str[startIndex:endIndex]
}


func ValidateClientGrant(grantTypes []string, requestedGrant string) error {
	if !slices.Contains(grantTypes, requestedGrant) {
		return fmt.Errorf("grant_type '%s' not allowed for this client", requestedGrant)
	}
	return nil
}
func GenerateAuthCode(length int) (string, error) {
    bytes := make([]byte, length)
    _, err := rand.Read(bytes)
    if err != nil {
        return "", err
    }
    return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(bytes), nil
}