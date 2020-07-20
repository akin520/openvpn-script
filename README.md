# openvpn-script
openvpn认证脚本

### 安装
    make

再将编译好的文件放到/etc/openvpn【根据自己实际情况放】

##  gin-dingding ##
钉钉报警使用，需要替换源码中的webHook和keyword

## server-go.conf ##
openvpn配置文件

    /usr/local/openvpn/sbin/openvpn --config /etc/openvpn/server-go.conf

## checkpwd ##
认证脚本

## vpn-connect  ##
连接时执行脚本

## vpn-disconnect ##
断开时执行脚本

## web管理 ##
[https://github.com/akin520/openvpn-beego](https://github.com/akin520/openvpn-beego "https://github.com/akin520/openvpn-beego")