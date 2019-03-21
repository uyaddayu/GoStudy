package main

import (
	"fmt"
	"github.com/djimenez/iconv-go"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"mahonia"
	"strings"
)

// 字符编码转换
func main() {
	str := "你好，世界"
	gb_out := make([]byte, 2*len(str)) // 给足够的长度
	utf8_out := make([]byte, 2*len(str))
	// 第一种方式
	iconv.Convert([]byte(str), gb_out, "utf-8", "gb2312")
	iconv.Convert(gb_out, utf8_out, "gb2312",
		"utf-8")
	fmt.Println("utf-8:", string(utf8_out))
	fmt.Println("gb2312:", string(gb_out))
	// 第二种方式
	fmt.Println("utf-8:", gbkToUtf8(string(gb_out)))
	fmt.Println("gb2312:", utf8ToGbk(str))
	// 第三种方式
	fmt.Println("utf-8:", gbkToUtf8ByTransform(string(gb_out)))
	fmt.Println("gb2312:", utf8ToGbkByTransform(str))
}

// gbk转utf8，可扩展为其他编码转utf8
func gbkToUtf8(str string) string {
	var dec mahonia.Decoder
	dec = mahonia.NewDecoder("gbk")

	if ret, ok := dec.ConvertStringOK(str); ok {
		return ret
	}
	return ""
}

// 将utf8转成gbk，可扩展为其他编码
func utf8ToGbk(str string) string {
	var enc mahonia.Encoder
	enc = mahonia.NewEncoder("gbk")

	if ret, ok := enc.ConvertStringOK(str); ok {
		return ret
	}
	return ""
}

func utf8ToGbkByTransform(str string) string {
	utf8_r := strings.NewReader(str)
	gb_r := transform.NewReader(utf8_r, simplifiedchinese.GBK.NewEncoder())
	gb_str, err := ioutil.ReadAll(gb_r)
	if err != nil {
		return ""
	}
	return string(gb_str)
}
func gbkToUtf8ByTransform(str string) string {
	gb_r := strings.NewReader(str)
	utf8_r := transform.NewReader(gb_r, simplifiedchinese.GBK.NewDecoder())
	utf8_str, err := ioutil.ReadAll(utf8_r)
	if err != nil {
		return ""
	}
	return string(utf8_str)
}
