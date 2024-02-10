package utils

import (
    "github.com/ModuleList/acceleratedcharging/log"
    "fmt"
    "os"
    "io"
    "os/exec"
    "strings"
    "time"
    "sync"
    "io/ioutil"
    "crypto/sha256"
)
var cmd *exec.Cmd
var env = os.Getenv("su")
var allFileNames = GetThermalFile()
func Verify(file string, verify string) {
    f, err := os.Open(file)
    if err != nil {
        log.Error(err)
    }
    defer f.Close()

    h := sha256.New()
    if _, err := io.Copy(h, f); err != nil {
        log.Error(err)
    }

    if fmt.Sprintf("%x", h.Sum(nil)) != verify {
        log.Info("签名校验失败")
        os.Exit(1)
    }
}
func IsFileEmpty(filename string) (bool) {
    info, err := os.Stat(filename)
    if err != nil {
        return false
    }
    return info.Size() == 0
}

func fileExists(filename string) bool {
    _, err := os.Stat(filename)
    return !os.IsNotExist(err)
}

func IfRoot() bool {
    if os.Geteuid() == 0 {
        return true
    }
    return false
}

func Shell(command string) string {
    cmd := exec.Command(env, "-c", command)
    output, err := cmd.CombinedOutput()
    if err != nil {
        return err.Error() // 返回错误信息
    }

    return strings.TrimSpace(string(output))
}

func Restart() {
    processes := []string{"mi_thermald", "thermal-engine", "thermald"}
    for _, p := range processes {
        Shell("stop " + p)
        Shell("start " + p)
    }
    Shell("chmod -R 0771 /data/vendor/thermal")
    Shell("chown -R root:system /data/vendor/thermal")
    Shell("chcon -R u:object_r:vendor_data_file:s0 /data/vendor/thermal")
    log.Info("已重启温控相关进程")
}

func Sleep(times int) { //硬核休眠
    sum := 1
    for sum <= times {
        sum = sum + 1
        time.Sleep(time.Second);
    }
}
func Modify() { //删除温控
    var File string;
    speed := GetConfig("temp.speed").(string)
    if fileExists("/sys/class/power_supply/battery/constant_charge_current") {
        File = "/sys/class/power_supply/battery/constant_charge_current"
    } else {
        File = "/sys/class/power_supply/battery/constant_charge_current_max"
    }
    err := ioutil.WriteFile("/data/adb/modules/acceleratedcharging/speed", []byte(speed+"\n"), 644)
    if err != nil {
        log.Info("写入最大充电电流文件失败")
        log.Error(err)
    }
    Shell("mount /data/adb/modules/acceleratedcharging/speed" + " " + File)
    log.Info("已修改最大充电电流" + speed + "到 " + File)
    if GetConfig("temp.dynamic").(bool) {
        log.Info("写入云控配置")
        var wg sync.WaitGroup
        wg.Add(len(allFileNames))
        normal, err := ioutil.ReadFile(GetConfig("configfile").(string))
        if err != nil {
            log.Error(err)
        }
        for _, FileNames := range allFileNames {
            go func(FileNames string) {
                defer wg.Done()
                err := ioutil.WriteFile("/data/vendor/thermal/config/" + FileNames, normal, 440)
                if err != nil {
                    log.Error(err)
                }
                Shell("dos2unix /data/vendor/thermal/config/" + FileNames)

            }(FileNames)
        }
    
        wg.Wait()
        Restart()
    }
    log.Info("已设置和动态温控")
}
func Recovery() {
    var File string;
    if fileExists("/sys/class/power_supply/battery/constant_charge_current") {
        File = "/sys/class/power_supply/battery/constant_charge_current"
    } else {
        File = "/sys/class/power_supply/battery/constant_charge_current_max"
    }
    Shell("rm -rf /data/vendor/thermal/config")
    Shell("mkdir -p /data/vendor/thermal/config")
    Shell("umount " + File)
    Restart()
    log.Info("已恢复温控和最大电流设置")
}