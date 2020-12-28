package main

import (
	"errors"
	"fmt"
	"github.com/13283339616/http"
)

type InstantItem struct {
	InstanceId string
	HostName   string
	App        string
	IpAddr     string
	Status     string
}
type Instance struct {
	Name     string
	Instance []InstantItem
}

type Applications struct {
	VersionsDelta string `json:"apps__hashcode"`
	AppsHashcode  string
	Application   []Instance
}
type ListItem struct {
	Applications Applications `json:"applications"`
}

//获取eureka实例列表
func GetList(url string) {
	headerMap := make(map[string]string, 2)
	headerMap["Content-Type"] = "application/json"
	headerMap["Accept"] = "application/json"
	var listItem ListItem
	list, err := http.Curl(url, "GET", "", headerMap, listItem)
	if err != nil {
		panic(err)
	}
	fmt.Println(list)
	listItem = list.(ListItem)
	fmt.Print(listItem.Applications.VersionsDelta)

}
func main() {
	GetList("http://39.97.194.236:8081/eureka/apps")
}

//获取eureka服务端的地址
func eurekaCurl(url, method string, data interface{}, headerMap map[string]string, act *interface{}) (*interface{}, error) {

	return nil, errors.New("test")
}
