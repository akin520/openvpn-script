package main

//username=anks password=anks ./getenv

import (
    "fmt"
    "path/filepath"
    "github.com/astaxie/beego/config"
    "time"
    "os"
    "crypto/md5"
    "encoding/hex"
    _ "github.com/go-sql-driver/mysql"
    "github.com/go-xorm/xorm"
)

type User struct {
    Id           int64
    Name         string `xorm:"unique"`
    Password     string
    Expired_time time.Time `xorm:"index"`
    Active       int
}

func main(){
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
    if err != nil{
        fmt.Println(err)
    }
    var USERNAME string
    var PASSWORD string
    var STATUS int
    USERNAME = os.Getenv("username")
    PASSWORD = os.Getenv("password")

    //password md5
    srcData := []byte(PASSWORD)
    hash := md5.New()
    hash.Write(srcData)
    cipherText2 := hash.Sum(nil)
    hexText := make([]byte, 32)
    hex.Encode(hexText, cipherText2)

    //fmt.Println(USERNAME,PASSWORD,string(hexText))

    var users []User
    err = engine.Table("user").Where("name = ?", USERNAME).And("password = ?", string(hexText)).And("active=1").And("expired_time>now()").Find(&users)
    if err != nil{
        fmt.Println(err)
    }
    
    if len(users) == 1 {
        STATUS=0
    }else{
        STATUS=1
    }
    //fmt.Println(STATUS)
    
    os.Exit(STATUS)
}
