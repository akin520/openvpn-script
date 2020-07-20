package main

//common_name="akin" trusted_ip="127.0.0.1" trusted_port="1999" proto="tcp-server" ifconfig_pool_remote_ip="10.0.0.1" bytes_received="1000" bytes_sent="1000" ./vpn-disconnect

import (
	"fmt"
	"github.com/astaxie/beego/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"os"
	"path/filepath"
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
	dataSourceName := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8", username, password, host, port, database)

	engine, err := xorm.NewEngine("mysql", dataSourceName)
	if err != nil {
		fmt.Println(err)
	}
	//写入数据库
	Common_name := os.Getenv("common_name")
	Trusted_ip := os.Getenv("trusted_ip")
	Trusted_port := os.Getenv("trusted_port")
	//Proto := os.Getenv("proto")
	Ifconfig_pool_remote_ip := os.Getenv("ifconfig_pool_remote_ip")
	Bytes_received := os.Getenv("bytes_received")
	Bytes_sent := os.Getenv("bytes_sent")

	loc, _ := time.LoadLocation("Asia/Shanghai")
	t := time.Now()
	t = t.In(loc)

	upsql := "update login_log SET end_time=?,bytes_received=?,bytes_sent=? WHERE trusted_ip=? and trusted_port=? and remote_ip=? and username=?"
	res, err := engine.Table("login_log").Exec(upsql, t, Bytes_received, Bytes_sent, Trusted_ip, Trusted_port, Ifconfig_pool_remote_ip, Common_name)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
	//res, err := engine.Table("login_log").Where("trusted_ip = ?", Trusted_ip).And("trusted_port= ?", Trusted_port).And("remote_ip=?", Ifconfig_pool_remote_ip).And("username=?", Common_name).Update(map[string]interface{}{"end_time": t, "bytes_received": Bytes_received, "bytes_sent": Bytes_sent})
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(res)

	//更新过期用户数据
	sql := "UPDATE user SET active=0 WHERE expired_time<now();"
	res1, err := engine.Table("user").Exec(sql)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res1)

	os.Exit(0)
}
