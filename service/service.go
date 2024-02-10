package service

import (
   "strings"
   "strconv"
   "github.com/ModuleList/acceleratedcharging/log"
   "github.com/ModuleList/acceleratedcharging/utils"
)

func Start() {
    status := false
    recovery := false
    log.Info("=====开始监听=====")
    for true {
        batterydata := utils.Shell("dumpsys battery")
        temperature, _ := strconv.Atoi(utils.Shell("cat /sys/class/power_supply/battery/temp"))
    
        if strings.Contains(batterydata, "status: 2") || strings.Contains(batterydata, "AC powered: true") && ! strings.Contains(batterydata, "level: 100") {
            if utils.GetConfig("debug").(bool) {
                log.Debug("监听到已进入充电状态")
                log.Debug(strconv.FormatBool(status))
            }
            if temperature < utils.GetConfig("temp.max").(int) {
                recovery = false
            }
            
            if ! status && ! recovery {
                status = true
                recovery = false
                utils.Modify()
                log.Info("已修改快充设置")
            } else {
                if temperature >= utils.GetConfig("temp.max").(int) && status && ! recovery {
                    status = false
                    recovery = true
                    utils.Recovery()
                    log.Info("温度超过限制 已恢复快充设置")
                }
            }
        } else {
            if status {
                utils.Recovery()
                log.Info("已拔出充电器或已满电")
            }
            if utils.GetConfig("debug").(bool) {
                log.Debug("未在充电")
            }
        }
        if utils.GetConfig("debug").(bool) {
            log.Debug("=====已完成一轮检测 休眠" + strconv.Itoa(utils.GetConfig("sleep").(int)) +"秒=====")
        }
        utils.Sleep(utils.GetConfig("sleep").(int))

    }
}