package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func HttpGet(url string) (result string, err error) {
	resp, err1 := http.Get(url)
	if err1 != nil {
		err = err1
		return
	}
	defer resp.Body.Close()

	buf := make([]byte, 4*1024)
	for {
		n, _ := resp.Body.Read(buf)
		if n == 0 {
			break
		}
		result += string(buf[:n])
	}
	return
}

func Page(i int, ch chan string) {
	url := "http://jwzx.node1.84cdn.com/kebiao/kb_stu.php?xh=2019" + strconv.Itoa(i)

	//开始获取页面信息
	result, err := HttpGet(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	//开始获取页面细节（名字）
	reg := regexp.MustCompile(`学生课表>>2019[0-9][0-9][0-9][0-9][0-9][0-9](?s:(.*?))</li>`)
	if reg == nil {
		fmt.Println("MustComple err")
		return
	}

	name := reg.FindAllStringSubmatch(result, -1)

	for _, text := range name {
		ch <- text[1]
	}
}

var myres = make(map[int]string, 100)

func Do() {
	var str string
	ch := make(chan string)
	for i := 215037; i <= 215060; i++ {
		//获取页面函数
		go Page(i, ch)
	}
	for j := 0; j < 500; j++ {
		str += <-ch
	}
	close(ch)

	m1 := make(map[string]int)
	for _, v := range strings.Fields(str) {
		if _, ok := m1[v]; ok {
			m1[v] = m1[v] + 1
		} else {
			m1[v] = 1
		}
	}
	//fmt.Println(m1)
	for key, value := range m1 {
		fmt.Printf("%q:%d\n", key, value)
	}
}

func main() {
	Do() //工作函数
}
