package utils

import (
	"encoding/json"
)

type Loggable interface{
	Loggable() map[string]interface{}
}

func Loggable2json(loginfo Loggable) string {
	jsonString, err := json.MarshalIndent(loginfo.Loggable(), "", "\t")
	if err != nil{
		//log.Error(err)
		return ""
	}
	return string(jsonString)
}
