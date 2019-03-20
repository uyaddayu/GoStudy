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
	fmt.Println(strings.EqualFold(str1, str2)) // 试了下不区分大小写。。
	fmt.Println(str1 == str2)                  // 这种够用了

	var str string = "this is an example of a string"
	fmt.Println(strings.HasPrefix(str, "th"))  // 判断前缀是否是某字符串
	fmt.Println(strings.HasSuffix(str, "ing")) // 判断后缀是否是某字符串
	fmt.Println(strings.Contains(str, "go"))   // 判断是否包含指定字符串
	fmt.Println(strings.Index(str, "a"))       // 得到指定字符串原字符串的第一个位置下标，返回-1代表不存在
	fmt.Println(strings.LastIndex(str, "a"))   // 得到指定字符串原字符串的最后一个位置下标，返回-1代表不存在
	// 返回str中每个单词的首字母都改为标题格式的字符串拷贝
	title := strings.Title(str)
	fmt.Println(title) // This Is An Example Of A String
	// 返回将所有字母都转为对应的标题版本的拷贝
	to_title := strings.ToTitle(str)
	fmt.Println(to_title) // THIS IS AN EXAMPLE OF A STRING
	// 返回count个str串联的字符串
	s_repeat := strings.Repeat(str, 3)
	fmt.Println(s_repeat) // this is an example of a stringthis is an example of a stringthis is an example of a string
	// 返回将str中前n个不重叠old子串都替换为new的新字符串，如果n<0会替换所有old子串
	s_replace := strings.Replace(str, "this", "<->", -1) // 此处n为-1
	fmt.Println(s_replace)                               // <-> is an example of a string
	// 返回将s前后端所有cutset包含的utf-8码值都去掉的字符串
	s, cutset := "#abc!!!", "#!"
	s_new := strings.Trim(s, cutset)
	fmt.Println(s, s_new) // #abc!!! abc
	// 返回将字符串按照空白（unicode.IsSpace确定，可以是一到多个连续的空白字符）分割的多个字符串
	s = "hello world! go language"
	s_fields := strings.Fields(s)
	for k, v := range s_fields {
		fmt.Println(k, v)
	}
	// 0 hello
	// 1 world!
	// 2 go
	// 3 language

	// 用去掉s中出现的sep的方式进行分割，会分割到结尾，并返回生成的所有片段组成的切片
	s_split := strings.Split(s, " ")
	fmt.Println(s_split) // [hello world! go language]

	// 将一系列字符串连接为一个字符串，之间用sep来分隔
	s_join := strings.Join([]string{"a", "b", "c"}, "/")
	fmt.Println(s_join) // a/b/c

	// 将s的每一个unicode码值r都替换为mapping(r)，返回这些新码值组成的字符串拷贝。如果mapping返回一个负值，将会丢弃该码值而不会被替换
	map_func := func(r rune) rune {
		switch {
		case r > 'A' && r < 'Z':
			return r + 32
		case r > 'a' && r < 'z':
			return r - 32
		}
		return r
	}
	s = "Hello World!"
	s_map := strings.Map(map_func, s)
	fmt.Println(s_map) // hELLO wORLD!
}
