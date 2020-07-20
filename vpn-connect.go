package main

//common_name="akin" trusted_ip="127.0.0.1" trusted_port="1999" proto="tcp-server" ifconfig_pool_remote_ip="10.0.0.1" bytes_received="1000" bytes_sent="1000" ./vpn-disconnect

import (
	"fmt"
	"github.com/astaxie/beego/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Login_log struct {
	Id             int64
	Username       string `xorm:"null"`
	Login_time     time.Time
	Trusted_ip     string `xorm:"null"`
	Trusted_port   string `xorm:"null"`
	Protocol       string `xorm:"null"`
	Remote_ip      string `xorm:"null"`
	End_time       time.Time
	Bytes_received string `xorm:"null"`
	Bytes_sent     string `xorm:"null"`
}

func Match(arr []string, str string) bool {
	for _, v := range arr {
		if v == str {
			return true

		}

	}
	return false

}

func main() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(dir)
	conf, err := config.NewConfig("ini", dir+"/app.conf")
	if err != nil {
		fmt.Println("new config failed, err:", err)
		return
	}
	host := conf.String("db_host")
	port := conf.String("db_port")
	username := conf.String("db_user")
	password := conf.String("db_pass")
	database := conf.String("db_name")
	lists := conf.String("list")
	URL := conf.String("url")
	dataSourceName := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8", username, password, host, port, database)

	engine, err := xorm.NewEngine("mysql", dataSourceName)
	if err != nil {
		fmt.Println(err)
	}
	//写入数据库
	Common_name := os.Getenv("common_name")
	Trusted_ip := os.Getenv("trusted_ip")
	Trusted_port := os.Getenv("trusted_port")
	Proto := os.Getenv("proto")
	Ifconfig_pool_remote_ip := os.Getenv("ifconfig_pool_remote_ip")

	var login_log = new(Login_log)
	login_log.Username = Common_name
	login_log.Login_time = time.Now()
	login_log.Trusted_ip = Trusted_ip
	login_log.Trusted_port = Trusted_port
	login_log.Protocol = Proto
	login_log.Remote_ip = Ifconfig_pool_remote_ip
	rows, err := engine.Table("login_log").Insert(login_log)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(rows)

	//登录报警API
	LL := strings.Split(lists, ",")
	if !Match(LL, Common_name) {
		url := URL + Trusted_ip + "%20" + Common_name + "%20vpn"
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(resp)
	} else {
		fmt.Println(Common_name)
	}

	//更新过期用户数据
	sql := "UPDATE user SET active=0 WHERE expired_time<now();"
	res, err := engine.Table("user").Exec(sql)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)

	os.Exit(0)
}
