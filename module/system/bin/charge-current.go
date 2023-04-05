package main

import (
    "fmt"
    "os"
    "os/exec"
    "io/ioutil"
    "strconv"
    "strings"
    "log"
    "time"
    "sync"
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
    os.Chmod("/data/vendor/thermal", 0771)
    os.Chown("/data/vendor/thermal", 0, 100)
    shell("chcon -R u:object_r:vendor_data_file:s0 /data/vendor/thermal", true)
    runlog("info","已重启温控相关进程");
}

func writeFile(File string, Text string) {
    file, err := os.Create(File)
    if err != nil {
        panic(err)
    }
    defer file.Close()
    file.WriteString(string(Text))
}


func modify_temperature_control_config(resourceDir string) { //删除温控
    if rmthermal == "true" {
        runlog("info","写入云控配置");
        generalthermal, err := ioutil.ReadFile(resourceDir + "general.conf")
        if err != nil {
            log.Fatal("[err]无法读取文件文件: ", err)
            os.Exit(1)
        }
        normalthermal, err := ioutil.ReadFile(resourceDir + "normal.conf")
        if err != nil {
            log.Fatal("[err]无法读取温控文件: ", err)
            os.Exit(1)
        }
        allFileNames := []string{
            "thermal-class0.conf",
            "thermal-hp-normal.conf",
            "thermal-india-class0.conf",
            "thermal-india-nolimits.conf",
            "thermal-india-normal.conf",
            "thermal-india-per-class0.conf",
            "thermal-india-per-normal.conf",
            "thermal-india-phone.conf",
            "thermal-nolimits.conf",
            "thermal-normal.conf",
            "thermal-per-class0.conf",
            "thermal-per-normal.conf",
            "thermal-phone.conf",
        }
        var wg sync.WaitGroup
        wg.Add(len(allFileNames))
    
        for _, FileNames := range allFileNames {
            go func(FileNames string) {
                defer wg.Done()
                file, err := os.Create("/data/vendor/thermal/config/" + FileNames)
                if err != nil {
                    log.Fatal("[err]无法写入温控文件: ", err)
                    os.Exit(1)
                }
                defer file.Close()
                if strings.Contains(FileNames, "normal") {
                    file.WriteString(string(normalthermal))
                } else {
                    file.WriteString(string(generalthermal))
                }
                os.Chmod("/data/vendor/thermal/config/" + FileNames, 444)
    	        shell("dos2unix /data/vendor/thermal/config/" + FileNames, true)
            }(FileNames)
        }
    
        wg.Wait()
        	restart()
        }
    writeFile("/data/adb/modules/acceleratedcharging/" + file, speed) //写入充电电流到模块缓存文件
    shell("mount /data/adb/modules/acceleratedcharging/" + file + " /sys/class/power_supply/battery/" + file,true); //通过mount命令挂载充电电流速度
    runlog("info","已修改充电最大电流设置和动态温控");
    shell("sed -i 's/\\[.*\\]/[已修改快充]/g' /data/adb/modules/acceleratedcharging/module.prop", true)
}

func recovery_temperature_control_config() {
    shell("rm -rf /data/vendor/thermal/config",true)
    shell("mkdir -p /data/vendor/thermal/config",true)
    restart()
    shell("umount /sys/class/power_supply/battery/" + file,true)
    runlog("info","已恢复温控和最大电流设置");
    shell("sed -i 's/\\[.*\\]/[未充电或满电]/g' /data/adb/modules/acceleratedcharging/module.prop", true)
}



func sleeps(times int) { //硬核休眠
    sum := 1
    for sum <= times {
        sum = sum + 1
        time.Sleep(time.Second);
    }
}

func runlog(level string, text string) {
    log.Println("[" + level + "]:" +text)
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
    shell("rm -rf /data/adb/modules/acceleratedcharging/charge-current.lo",true);
    file, err := os.Create("charge-current.log")
    if err != nil {
        log.Fatal("[err]无法创建日志文件: ", err)
    }
    defer file.Close()
    log.SetOutput(file)
    log.SetFlags(log.Ldate | log.Ltime)
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
                        runlog("info","检测到充电状态");
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
                    runlog("info","已充满电");
                    recovery_temperature_control_config(); //恢复
                    start = 0
                    stop = 1
                }
            }
        }
        sleeps(timesleep); //休眠
    }
}
