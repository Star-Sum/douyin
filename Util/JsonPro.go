package Util

import (
	"douyin/Log"
	"encoding/json"
)

// JsonUnmarshal 将json输出为结构体
func JsonUnmarshal(data []byte, targetElem any) (any, error) {
	err := json.Unmarshal(data, &targetElem)
	if err != nil {
		Log.ErrorLogWithoutPanic("JSON parsing failed!", err)
		return nil, err
	}
	Log.NormalLog("JSON parsing successful!", err)
	return targetElem, err
}

// JsonMarshal 将结构体输出为json
func JsonMarshal(targetElem any) ([]byte, error) {
	jsonFile, err := json.Marshal(&targetElem)
	if err != nil {
		Log.ErrorLogWithoutPanic("Failed to output JSON!", err)
		return []byte{}, err
	}
	Log.NormalLog("Successfully output JSON!", err)
	return jsonFile, err
}
