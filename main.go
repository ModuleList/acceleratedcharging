package main

import (
    "flag"
    "github.com/ModuleList/acceleratedcharging/log"
    "github.com/ModuleList/acceleratedcharging/utils"
    s "github.com/ModuleList/acceleratedcharging/service"
)

var service bool
var command string
var signature string

func main() {
    log.Init()
    utils.Verify("module.prop", signature)
    if utils.Shell("getprop ro.product.manufacturer") == "Xiaomi" {
        empty := utils.IsFileEmpty("/vendor/bin/mi_thermald")
        if ! empty {
            log.Info("/vendor/bin/mi_thermald文件不存在 请开启实验性移除温控功能")
        }
    }
    flag.BoolVar(&service, "service", false, "Whether to enable background running")
    flag.StringVar(&command, "command", "", "Execute modification of fast charging[modify]/recovery fast charging[recovery] settings")
    flag.Parse()
    

    if ! utils.IfRoot() {
        log.Info("请使用root用户运行")
        return
    }
    if service {
        s.Start()
    }
    if command == "modify" {
        utils.Modify()
    } else if command == "recovery" {
        utils.Recovery()
    }
}
