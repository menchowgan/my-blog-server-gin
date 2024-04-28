HOST_LIST="
  47.116.163.91
"

APP=gmc-blog-server
APP_PATH=/root/gmc

build () {
  rm -rf ./$APP
  GOOS=linux GOARCH=amd64 go build -o $APP
}

upload () {
  for host in $HOST_LIST;do
    ssh root@$host "cd ${APP_PATH};\
            killall -0 ${APP} &>/dev/null && killall ${APP} \
            rm -rf ${APP};
          "
    
    scp ./$APP root@$host:$APP_PATH/

    ssh root@$host "cd ${APP_PATH} && nohup ./${APP} &>/dev/null"
  done
}

build
upload