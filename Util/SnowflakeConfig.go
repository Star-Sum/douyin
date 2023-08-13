package Util

import (
	"douyin/Log"
	"gopkg.in/yaml.v3"
	"os"
)

type SnowflakeConfig struct {
	DataCenterId int64 `yaml:"dataCenterId"`
	MachineId    int64 `yaml:"machineId"`
}

var Snowflake SnowflakeConfig

func init() {

	// 读取配置文件
	data, err := os.ReadFile("./Config/config.yaml")
	if err != nil {
		Log.ErrorLogWithPanic("Failed to read snowflakeConfig configuration file!", err)
	}
	Log.NormalLog("Successfully read snowflakeConfig configuration file!", err)
	// 解析配置文件内容
	err = yaml.Unmarshal(data, &Snowflake)
	if err != nil {
		Log.ErrorLogWithPanic("Failed to parse snowflakeConfig", err)
	}
	Log.NormalLog("Successfully obtained snowflakeConfig settings!", err)
}
