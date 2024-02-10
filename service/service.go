package service

import (
   "io/ioutil"
   "strings"
   "strconv"
   "github.com/ModuleList/acceleratedcharging/log"
   "github.com/ModuleList/acceleratedcharging/utils"
)

func Start() {
    status := false
    log.Info("开始监听")
    if utils.GetConfig("temp.beta").(bool) {
        log.Info("已开启 实验性动态修改温控(支持非小米机型)")
    }
    for true {
        batterydata := utils.Shell("dumpsys battery")
        bytes, err := ioutil.ReadFile("/sys/class/power_supply/battery/temp")
        if err != nil {
            log.Error(err)
        }
        temperature, _ := strconv.Atoi(string(bytes))
    
        if strings.Contains(batterydata, "status: 2") || strings.Contains(batterydata, "AC powered: true") {
            if utils.GetConfig("debug").(bool) {
                log.Debug("监听到已进入充电状态")
            }
            if ! strings.Contains(batterydata, "level: 100") && ! status {
                status = true
                utils.Modify()
                log.Info("已修改快充")

            } else if temperature > utils.GetConfig("temp.max").(int) && ! strings.Contains(batterydata, "level: 100") && status {
                status = false
                utils.Recovery()
                log.Info("温度超过限制")
            } else if strings.Contains(batterydata, "level: 100") && status {
                status = false
                utils.Recovery()
                log.Info("电池已充满")
            } else {
                status = false
                utils.Recovery()
                log.Info("未知状态 已初始化")
            }
        } else {
            if utils.GetConfig("debug").(bool) {
                log.Debug("未在充电")
            }
        }
        if utils.GetConfig("debug").(bool) {
            log.Debug("已完成一轮检测 休眠" + strconv.Itoa(utils.GetConfig("sleep").(int)) +"秒")
        }
        utils.Sleep(utils.GetConfig("sleep").(int))

    }
}