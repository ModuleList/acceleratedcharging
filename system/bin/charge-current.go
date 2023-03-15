package main

import (
    "fmt"
    "os"
    "os/exec"
    "strconv"
    "strings"
    "time"
)

var (
    temperature     string
    batterylevel     string
    temperaturewall    string
    speed       string
    rmthermal       string
    timesleep       string
    file    string
    thermalfile    string
    times       int
    start     int
    stop     int
)

func shell(command string, su bool ) string { //调用shell执行命令(root权限)
    var output []byte
    var err error
    var cmd *exec.Cmd
    if su == true {
        cmd = exec.Command("su", "-c", command);
    } else {
        cmd = exec.Command("bash", "-c", command);
    }
    if output, err = cmd.CombinedOutput(); err == nil {
    }
    returntext := strings.Trim(string(output), "\n")
    return returntext
}

func restart() {
    processes := []string{"mi_thermald", "thermal-engine"}
    for _, p := range processes {
        shell("stop " + p, true)
        shell("start " + p, true)
    }
    shell("chmod -R 0771 /data/vendor/thermal", true)
    shell("chown -R root:system /data/vendor/thermal", true)
    shell("chcon -R u:object_r:vendor_data_file:s0 /data/vendor/thermal", true)
    runlog("info","已重启温控相关进程");
}

func modify_temperature_control_config(resourceDir string) string{ //删除温控
    if rmthermal == "true" {
        runlog("info","写入云控配置");
    	allFileNames := []string{
    		"thermal-4k.conf",               //所有要处理的文件名
    		"thermal-arvr.conf",
    		"thermal-camera.conf",
    		"thermal-class0.conf",
    		"thermal-hp-mgame.conf",
    		"thermal-hp-normal.conf",
    		"thermal-huanji.conf",
    		"thermal-india-4k.conf",
    		"thermal-india-camera.conf",
    		"thermal-india-class0.conf",
    		"thermal-india-huanji.conf",
    		"thermal-india-mgame.conf",
    		"thermal-india-navigation.conf",
    		"thermal-india-nolimits.conf",
    		"thermal-india-normal.conf",
    		"thermal-india-per-camera.conf",
    		"thermal-india-per-class0.conf",
    		"thermal-india-per-navigation.conf",
    		"thermal-india-per-normal.conf",
    		"thermal-india-per-video.conf",
    		"thermal-india-phone.conf",
    		"thermal-india-tgame.conf",
    		"thermal-india-video.conf",
    		"thermal-india-videochat.conf",
    		"thermal-india-yuanshen.conf",
    		"thermal-mgame.conf",
    		"thermal-navigation.conf",
    		"thermal-nolimits.conf",
    		"thermal-normal.conf",
    		"thermal-per-camera.conf",
    		"thermal-per-class0.conf",
    		"thermal-per-navigation.conf",
    		"thermal-per-normal.conf",
    		"thermal-per-video.conf",
    		"thermal-phone.conf",
    		"thermal-tgame.conf",
    		"thermal-video.conf",
    		"thermal-videochat.conf",
    		"thermal-yuanshen.conf",
    	}
    	for _, p := range allFileNames {
    	    if strings.Contains(p, "normal") {
    	        shell("cp -Rf " + resourceDir + p + " /data/vendor/thermal/config/" + p,true)
    	        shell("chmod 444 /data/vendor/thermal/config/" + p,true)
    	        shell("dos2unix /data/vendor/thermal/config/" + p,true)
    	    } else {
    	        shell("cp -Rf " + resourceDir + "general.conf" + " /data/vendor/thermal/config/" + p,true)
    	        shell("chmod 444 /data/vendor/thermal/config/" + p,true)
    	        shell("dos2unix /data/vendor/thermal/config/" + p,true)
    	    }
    	}
    	restart()
    }
    shell("echo '" + speed + "' > /data/adb/modules/acceleratedcharging/" + file,true) //写入充电电流到模块缓存文件
    shell("mount /data/adb/modules/acceleratedcharging/" + file + " /sys/class/power_supply/battery/" + file,true); //通过mount命令挂载充电电流速度
    runlog("info","修改充电最大电流");
    return "true"
}

func recovery_temperature_control_config() {
    shell("rm -rf /data/vendor/thermal/config",true)
    shell("mkdir -p /data/vendor/thermal/config",true)
    restart()
    runlog("info","恢复云控配置");
    shell("umount /sys/class/power_supply/battery/" + file,true)
    runlog("info","恢复充电最大电流设置");
}



func sleeps(times int) { //硬核休眠
    sum := 1
    for sum <= times {
        sum = sum + 1
        time.Sleep(time.Second);
    }
}

func runlog(level string, text string) {
    shell("echo \"[$(TZ=Asia/Shanghai date \"+%Y-%m-%d-%H:%M:%S\")]"+ level + ": " + text + "\" >>/data/adb/modules/acceleratedcharging/charge-current.log",true);
    fmt.Println(text)
}
func main() {
    //读取命令行参数
    args := os.Args
    if args == nil || len(args) < 6{
        fmt.Println("未传入参数");
        return
    }
    speed = args[1]
    temperaturewall = args[2]
    timesleep = args[3]
    rmthermal = args[4]
    file = args[5]
    thermalfile = args[6]
    timesleep, err := strconv.Atoi(timesleep); //将string类型转为int类型
    if err != nil {
        fmt.Println(err)
        return
    }
    //初始化变量
    start = 0
    stop = 0
    shell("rm -rf /data/adb/modules/acceleratedcharging/charge-current.log && touch /data/adb/modules/acceleratedcharging/charge-current.log",true);
    runlog("info","初始化完成✓");
    for true { //循环
        var batterydata = shell("dumpsys battery",true)
        temperature = shell("cat /sys/class/power_supply/battery/temp",true);
        var dl = strings.Contains(batterydata, "status: 2")
        if dl { //判断是否在充电
            if temperature > temperaturewall {
                runlog("info","温度超过限制");
                recovery_temperature_control_config();
                start = 0
                stop = 1
            } else {
                if strings.Contains(batterydata, "level: 100") {
                    runlog("info","已充满电");
                    recovery_temperature_control_config(); //恢复
                    start = 1
                    stop = 0
                } else {
                    if start == 0 {
                        modify_temperature_control_config(thermalfile); //删除温控 修改充电速度
                        start = 1
                        stop = 0
                    }
                }
            }
        } else {
            if strings.Contains(batterydata, "level: 100") {
            } else {
                if stop == 0 {
                    recovery_temperature_control_config(); //恢复
                    start = 0
                    stop = 1
                }
            }
        }
        runlog("info","end");
        sleeps(timesleep); //休眠
    }
}
