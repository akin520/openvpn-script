package main

//https://blog.csdn.net/e421083458/article/details/91994788
import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Msg struct {
	Errcode int
	Errmsg  string
}

func main() {
	router := gin.Default()

	router.GET("/ding/:msg", func(c *gin.Context) {
		msg := c.Param("msg")
		log.Println("msg:", msg)
		status, code := SendDing(msg)
		if status == 200 && code == 0 {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		} else {
			c.JSON(http.StatusOK, gin.H{"status": "error"})
		}
	})
	router.Run(":8084")
}

func SendDing(msg string) (int, int) {
	//请求地址模板
	webHook := `https://oapi.dingtalk.com/robot/send?access_token=241a69afec747a40121db7d15d72eb1e260d54618fffe75063c2`
	keyword := `[vpn]: `
	content := `{"msgtype": "text",
		"text": {"content": "` + keyword + msg + `"}
	}`
	log.Println("json:", content)
	//创建一个请求
	req, err := http.NewRequest("POST", webHook, strings.NewReader(content))
	if err != nil {
		log.Println("NewRquest:", err)
	}

	client := &http.Client{}
	//设置请求头
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	//发送请求
	resp, err := client.Do(req)
	//log.Println(resp)
	//关闭请求
	defer resp.Body.Close()

	if err != nil {
		log.Println("resp:", err)
	}
	//json转换
	message := Msg{}
	body, err := ioutil.ReadAll(resp.Body)
	log.Println("request context:", string(body))
	err = json.Unmarshal(body, &message)
	if err != nil {
		log.Println("msg:", err)
	}
	return resp.StatusCode, message.Errcode
}
