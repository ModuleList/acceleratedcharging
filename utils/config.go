package utils

import (
    "github.com/spf13/viper"
    "github.com/ModuleList/acceleratedcharging/log"
)


func GetConfig(key string) interface{} {
    viper.SetConfigName("config")
    viper.AddConfigPath(".")

    err := viper.ReadInConfig()
    if err != nil {
        log.Info("读取温配置文件失败")
        log.Error(err)
    }
    return viper.Get(key)
}
func GetThermalFile() []string {
    var ret []string
    file := GetConfig("temp.file")
    for _, val := range file.([]interface{}) {
        ret = append(ret, val.(string))
    }
    return ret
}