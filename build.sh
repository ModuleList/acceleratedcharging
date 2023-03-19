
if ! command -v go >/dev/null 2>&1; then
  echo "未安装golang"
  exit
fi
if ! command -v zip >/dev/null 2>&1; then
  echo "未安装zip"
  exit
fi

go build -ldflags="-s -w" -o ./module/system/bin/charge-current ./module/system/bin/charge-current.go
zip -q -r acceleratedcharging.zip ./module/*