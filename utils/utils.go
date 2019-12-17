package utils

import (
	"encoding/json"
	logging "github.com/ipfs/go-log"
)

var coreLogSystems = []string{
	"hon.bitswap",
	"hon.engine",
}

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

func GetLogSystems()[]string{
	return logging.GetSubsystems()
}

func SetCoreLogLevel(level string) error {
	for _, s := range coreLogSystems{
		err := logging.SetLogLevel(s, level)
		if err != nil {
			return err
		}
	}
	return nil
}