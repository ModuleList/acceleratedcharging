#!/system/bin/sh
MODDIR=${0%/*}

until [ $(getprop init.svc.bootanim) = "stopped" ] ; do
    sleep 5
done


nohup "${MODDIR}/bin/charge-current" -service &
