#!/system/bin/sh
if ! command -v go >/dev/null 2>&1; then
  echo "未安装golang"
  exit
fi
if ! command -v zip >/dev/null 2>&1; then
  echo "未安装zip"
  exit
fi
cd module
go build -ldflags="-s -w" -o ./system/bin/charge-current ./system/bin/charge-current.go
if [ -f ./system/bin/charge-current ];then
    rm -rf ./system/bin/charge-current.go
else
    echo "编译失败"
    exit
fi
zip -r acceleratedcharging.zip *
echo "打包编译成功"
echo "模块在module文件夹下"