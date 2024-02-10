#!/system/bin/sh

if [ ! $KSU ];then
    ui_print "- Magisk ver: $MAGISK_VER"
    if [[ $($MAGISK_VER | grep "kitsune") ]] || [[ $($MAGISK_VER | grep "delta") ]]; then
        ui_print "*********************************************************"
        ui_print "Magisk delta and magisk kitsune are not supported"
        echo "">remove
        abort "*********************************************************"
    fi
    
    ui_print "- Magisk version: $MAGISK_VER_CODE"
    if [ "$MAGISK_VER_CODE" -lt 26301 ]; then
        ui_print "*********************************************************"
        ui_print "! Please use Magisk alpha 26301+"
        abort "*********************************************************"
    fi
elif [ $KSU ];then
    ui_print "- KernelSU version: $KSU_KERNEL_VER_CODE (kernel) + $KSU_VER_CODE (ksud)"
    if ! [ "$KSU_KERNEL_VER_CODE" ] || [ "$KSU_KERNEL_VER_CODE" -lt 11413 ]; then
        ui_print "*********************************************************"
        ui_print "! KernelSU; version is too old!"
        ui_print "! Please update KernelSU to latest version"
        abort "*********************************************************"
    fi
else
    ui_print "! Unknown Module Manager"
    ui_print "$(set)"
    abort
fi
function mk_thermal_folder(){
    if [ ! -d "/data/vendor/thermal" ];then
        mkdir -p /data/vendor/thermal/config
    fi
    chattr -R -i -a /data/vendor/thermal
    chmod 0771 /data/vendor/thermal
    chmod 0771 /data/vendor/thermal/config
    chmod -R 644 /data/vendor/thermal/config
    chown -R root:system /data/vendor/thermal
}


#避免冻结电量和性能
function enable_miui_powerkeeper(){
    if test "$( pm list package | grep -w 'com.miui.powerkeeper' | wc -l)" -gt "0" ;then
        pm enable com.miui.powerkeeper >/dev/null 2>&1
        pm unsuspend com.miui.powerkeeper >/dev/null 2>&1
        pm unhide com.miui.powerkeeper >/dev/null 2>&1
        pm install-existing --user 0 com.miui.powerkeeper >/dev/null 2>&1
    fi
}

#重新启用电量与性能
function call_cloud_conf_release(){
    pm enable com.miui.powerkeeper/com.miui.powerkeeper.cloudcontrol.CloudUpdateReceiver >/dev/null 2>&1
    pm enable com.miui.powerkeeper/com.miui.powerkeeper.cloudcontrol.CloudUpdateJobService >/dev/null 2>&1
    pm enable com.miui.powerkeeper/com.miui.powerkeeper.ui.CloudInfoActivity >/dev/null 2>&1
    pm enable com.miui.powerkeeper/com.miui.powerkeeper.statemachine.PowerStateMachineService >/dev/null 2>&1
    am broadcast --user 0 -a update_profile com.miui.powerkeeper/com.miui.powerkeeper.cloudcontrol.CloudUpdateReceiver >/dev/null 2>&1
}

enable_miui_powerkeeper
call_cloud_conf_release
mk_thermal_folder

set_perm_recursive $MODPATH 0 0 0777 0777
set_perm $MODPATH/bin/charge-current 0 0 0777 0777
