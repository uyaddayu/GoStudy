package main

import (
	"fmt"
	"strings"
)

// 简单的字符串
func main() {
	str1 := "Go"
	str2 := "go"
	fmt.Println(str1, str2)
	// 将小写转成大写
	str2 = strings.ToUpper(str2)
	// 将大写转成小写
	str1 = strings.ToLower(str1)
	fmt.Println(str1, str2)
	// 判断是否相等
	fmt.Println(strings.EqualFold(str1, str2))
	fmt.Println(str1 == str2)
}
