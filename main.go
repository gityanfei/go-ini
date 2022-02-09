package main

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
)

type connMysql struct {
	User     string `ini:"user"`
	Password string `ini:"password"`
	Host     string `ini:"host"`
	Port     uint32 `ini:"port"`
}
type connRedis struct {
	Password string `ini:"password"`
	Host     string `ini:"host"`
	Port     uint32 `ini:"port"`
	Database uint32 `ini:"database"`
}

type config struct {
	connMysql `ini:"mysql"`
	connRedis `ini:"redis"`
}

// config 用于从文本文件中解析配置
func readConfigFromIni(filename string, a interface{}) error {
	//定义取出文件中的配置项后存储的变量
	var configKey string
	var configValue string
	//定义存储节的变量
	var sectionName string
	//定义获取反射的动态类型和动态值
	c := reflect.TypeOf(a)
	d := reflect.ValueOf(a)
	// 1.读文件
	b, error := ioutil.ReadFile(filename)
	if error != nil {
		fmt.Printf("open file %s fail: %v", filename, error)
	}
	//1.1按行切割文件到line切片中
	line := strings.Split(string(b), "\r\n")
	//1.2变量中的每一行
	for k, v := range line {
		//去除每一行中首尾的空格
		noSpaceLine := strings.TrimSpace(v)
		//2.1如果去除前后空格后为空行或开头为# ;则表示是注释跳过处理
		if len(noSpaceLine) == 0 || strings.HasPrefix(noSpaceLine, "#") || strings.HasPrefix(noSpaceLine, ";") {
			continue
		}
		//2.2判断当前行是否为节
		if strings.HasPrefix(noSpaceLine, "[") && strings.HasSuffix(noSpaceLine, "]") {
			//将节的名称取出
			sectionName = strings.TrimSpace(noSpaceLine[1 : len(noSpaceLine)-1])
			//如果得到的节名称是空的，则报语法错误
			if len(sectionName) == 0 {
				return fmt.Errorf("sytax is error at:%v", k+1)
			}
		} else { //2.3不为节就是为键值对
			//2.3.1开头是=表示语法出错,除开节没有=号则是语法错误
			if strings.HasPrefix(noSpaceLine, "=") || !strings.Contains(noSpaceLine, "=") {
				return fmt.Errorf("sytax is error at:%v", k+1)
			}
			//2.3.2对根据等号切割字符串，将key configValue两边的空白字符去除
			configKey = strings.TrimSpace(strings.Split(noSpaceLine, "=")[0])
			configValue = strings.TrimSpace(strings.Split(noSpaceLine, "=")[1])
			//遍历结构体中第一层字段
			for i := 0; i < c.Elem().NumField(); i++ {
				//遍历结构体中第二层字段
				for m := 0; m < c.Elem().Field(i).Type.NumField(); m++ {
					//如果节的名称和第一层中元素TAG相同并且和第一层元素中的结构体中的元素也相同则设置值
					if c.Elem().Field(i).Tag.Get("ini") == sectionName && c.Elem().Field(i).Type.Field(m).Tag.Get("ini") == configKey {
						if !d.Elem().Field(i).Field(m).CanSet() {
							fmt.Println("can not be set")
						}
						f := d.Elem().Field(i).Field(m)
						//通过反射更改结构体值时，结构体变量有多种，需要判断后赋值
						if f.Kind() == reflect.String {
							f.SetString(configValue)
						} else if f.Kind() == reflect.Uint32 {
							after, err := strconv.ParseUint(configValue, 10, 32)
							if err != nil {
								return fmt.Errorf("sytax is error at:%v,%v", k+1, err)
							}
							f.SetUint(after)
						}
					}
					// break
				}
			}
		}
	}
	return nil
}
func main() {
	aa := config{}
	err := readConfigFromIni("D:/Go/src/code.oldboyedu.com/studygo/day07/01config/config.ini", &aa)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%#v\n", aa)

}
