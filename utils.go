package main

import "encoding/json"

func PayloadToJSONString(payload Payload) string {
	bytes, _ := json.Marshal(payload)

	return string(bytes)
}
