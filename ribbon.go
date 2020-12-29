package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/13283339616/balance"
	"github.com/13283339616/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

//management.port": "8081"
type Meta struct {
	ManagementPort string `json:"management.port"`
}

type InstantItem struct {
	InstanceId string `json:"instanceId"`
	HostName   string `json:"hostName"`
	App        string `json:"app"`
	IpAddr     string `json:"ipAddr"`
	Status     string `json:"status"`
	Metadata   Meta   `json:"metadata"`
}
type Instance struct {
	Name     string        `json:"name"`
	Instance []InstantItem `json:"instance"`
}

type Applications struct {
	VersionsDelta string     `json:"versions__delta"`
	AppsHashcode  string     `json:"apps__hashcode"`
	Application   []Instance `json:"application"`
}
type ListItem struct {
	Applications Applications `json:"applications"`
}

var instanceMap = make(map[string][]*balance.Instance, 32)
var shutdown = make(chan os.Signal, 1)

func RibbonInit(url string) {
	InitOrigin(url)
	go InitMap(url)
}

func InitOrigin(url string) {
	list := GetList("http://39.97.194.236:8081/eureka/apps")
	application := list.Applications.Application
	for _, v := range application {
		for _, item := range v.Instance {
			port, _ := strconv.Atoi(item.Metadata.ManagementPort)
			if _, ok := instanceMap[v.Name]; ok {
				items := instanceMap[v.Name]
				newItem := balance.NewInstance(item.IpAddr, int64(port), int64(1))
				items = append(items, newItem)
				instanceMap[v.Name] = items

			} else {
				items := make([]*balance.Instance, 0, 0)
				newItem := balance.NewInstance(item.IpAddr, int64(port), int64(1))
				items = append(items, newItem)
				instanceMap[v.Name] = items
			}
		}
	}
}

//每一分钟拉去一次
func InitMap(url string) {
	for {
		time.Sleep(time.Second * 180)
		InitOrigin(url)
		fmt.Println(syscall.Getegid())
		fmt.Println(syscall.Getppid())
	}
}

//获取eureka实例列表
func GetList(url string) ListItem {
	headerMap := make(map[string]string, 2)
	headerMap["Content-Type"] = "application/json"
	headerMap["Accept"] = "application/json"
	listItem := new(ListItem)
	content, err := http.Curl(url, "GET", "", headerMap)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal([]byte(content), listItem)
	if err != nil {
		panic(err)
	}
	return *listItem
}
func getUrlByName(name string) (string, error) {
	if itemMap, ok := instanceMap[name]; ok {
		_, err := balance.DoBalance("shuffle2", itemMap)
		if err != nil {
			panic(err)
		}
		result := itemMap[0].GetResult()
		resArr := strings.Split(result, ";")
		return resArr[0], nil

	} else {
		return "nil", errors.New("查询的地址不存在")
	}
}
func main() {
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	RibbonInit("")
	fmt.Println(getUrlByName("GSVIP-MANAGE-WECHAT"))
	<-shutdown
	fmt.Println("退出")
}

//获取eureka服务端的地址
func eurekaCurl(url, method string, data interface{}, headerMap map[string]string, act *interface{}) (*interface{}, error) {

	return nil, errors.New("test")
}
