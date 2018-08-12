package main

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const isHaveNumber = false                         // 用户名是否含有数字
const requestIntervalTime = 200 * time.Millisecond // 请求间间隔
const requestURL = "https://gitee.com/check"
const cookie = "aliyungf_tc=AQAAAEuNgn5fDwkA5b8QtzELCwpJY/OL; user_locale=zh-CN; oschina_new_user=false; tz=Asia%2FShanghai; Hm_lvt_24f17767262929947cc3631f99bfd274=1533520007; visit-gitee-11=1; relative_time=true; OUTFOX_SEARCH_USER_ID_NCOO=334023554.32660013; gitee-session-n=BAh7CkkiD3Nlc3Npb25faWQGOgZFVEkiJTJkZWU2NzBiMWZhNmMwMjYyZTc2ZDgxYjYxZmFmMzJlBjsAVEkiF21vYnlsZXR0ZV9vdmVycmlkZQY7AEY6CG5pbEkiEF9jc3JmX3Rva2VuBjsARkkiMWw4SGVLSGtRWW9BT0tlVWRnRGs2SFpXRjRTNHU3dkNMaW5LZmx4S1I0eHc9BjsARkkiDGJyb3dzZXIGOwBUVEkiDGNhcHRjaGEGOwBGSSILSkhGTFpVBjsAVA%3D%3D--99e9ad42347684a8a948053197d91da9fa746755; ___rl__test__cookies=1533601066748; Hm_lpvt_24f17767262929947cc3631f99bfd274=1533601091"
const userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3510.0 Safari/537.36"
const xCSRFToken = "l8HeKHkQYoAOKeUdgDk6HZWF4S4u7vCLinKflxKR4xw="
const xRequestedWith = "XMLHttpRequest"
const fileJSON = "./data.json"

// 储存集合(1296)
var assemble []string

// 过滤后的集合(936)(33696)
var assembleFilter []string

// 存放json数据
var dataJSON map[string]interface{}

func initJSONFile() {
	_, err := os.Open(fileJSON)
	// 文件不存在则，创建
	if err != nil {
		outputFile, outputError := os.OpenFile(fileJSON, os.O_WRONLY|os.O_CREATE, 0666)
		if outputError != nil {
			fmt.Printf("An error occurred with file opening or creation\n")
			return
		}
		defer outputFile.Close()
		outputWriter := bufio.NewWriter(outputFile)
		outputString := "{}"

		outputWriter.WriteString(outputString)
		outputWriter.Flush()
	}
}

// 读取json文件
func readJSONFile() []byte {
	fp, err := os.OpenFile(fileJSON, os.O_RDONLY, 0755)
	defer fp.Close()
	if err != nil {
		log.Fatal(err)
	}
	data := make([]byte, 100000000)
	n, err := fp.Read(data)
	if err != nil {
		log.Fatal(err)
	}
	return data[:n]
}

// 写入json文件
func writeJSONFile() {
	// return
	file, _ := os.OpenFile(fileJSON, os.O_CREATE|os.O_WRONLY, 0666)
	defer file.Close()
	enc := json.NewEncoder(file)
	err := enc.Encode(dataJSON)
	if err != nil {
		log.Println("Error in encoding json")
	}
}

func post(username string) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", requestURL, strings.NewReader("do=user_username&val="+username))

	if err != nil {
		fmt.Println(err)
	}

	req.Header.Add("Accept", "*/*")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7,vi;q=0.6")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("Cookie", cookie)
	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("X-CSRF-Token", xCSRFToken)
	req.Header.Add("X-Requested-With", xRequestedWith)
	resp, err := client.Do(req)

	var reader io.ReadCloser
	if resp.Header.Get("Content-Encoding") == "gzip" {
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return
		}
	} else {
		reader = resp.Body
	}

	defer resp.Body.Close()

	body, err1 := ioutil.ReadAll(reader)
	if err1 != nil {
		fmt.Println(err1)
	}

	if string(body) == "1" {
		dataJSON[username] = "ok"
		fmt.Printf("ok")
	} else {

		if string(body) == "地址已存在" {
			dataJSON[username] = "already"
			fmt.Printf("already")
		} else {
			dataJSON[username] = "no"
			fmt.Printf("no")
		}
	}
	writeJSONFile()
}

func makeAssemble() (data *[]string) {
	if isHaveNumber {
		data = &[]string{
			"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
			"a", "b", "c", "d", "e", "f",
			"g", "h", "i", "j", "k", "l",
			"m", "n", "o", "p", "q", "r",
			"s", "t", "u", "v", "w", "x",
			"y", "z"}
	} else {
		data = &[]string{
			"a", "b", "c", "d", "e", "f",
			"g", "h", "i", "j", "k", "l",
			"m", "n", "o", "p", "q", "r",
			"s", "t", "u", "v", "w", "x",
			"y", "z"}
	}
	return
}

// 生成可能的集合
func calcAssemble(data *[]string) {
	for _, v0 := range *data {
		for _, v1 := range *data {
			for _, v2 := range *data {
				assemble = append(assemble, v0+v1+v2)
			}
		}
	}
}

func filter(a []string) {
	for _, v := range a {
		// 判断首字是否为数字
		_, err := strconv.Atoi(v[0:1])
		// 不为数字则加入
		if err != nil {
			assembleFilter = append(assembleFilter, v)
		}
	}
}

func cycleRequest(index int) {
	if index < len(assembleFilter) {
		post(assembleFilter[index])
		fmt.Printf("-")
		fmt.Printf(assembleFilter[index])
		fmt.Printf("|")
		time.Sleep(requestIntervalTime)
		index++
		cycleRequest(index)
	} else {
		fmt.Println("扫描完毕")
	}
}

func main() {
	initJSONFile()

	data := readJSONFile()

	err := json.Unmarshal(data, &dataJSON)
	if err != nil {
		fmt.Println("json.Unmarshal err")
	}

	calcAssemble(makeAssemble())
	// 如果允许用户名包含数字，则筛掉数字开头的用户名(gitee的用户名只能字母开头)
	if isHaveNumber {
		filter(assemble)
	} else {
		assembleFilter = assemble
	}
	fmt.Printf("需请求的个数: %d\n", len(assembleFilter))
	cycleRequest(0)
}
