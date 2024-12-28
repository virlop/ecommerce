package strs

import "encoding/json"

// conversion a json
func ToJson(obj interface{}) string {
	jsonData, _ := json.Marshal(obj)
	return string(jsonData)
}
