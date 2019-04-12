package main

import (
	"fmt"
	"regexp"
)

func main() {
	test1()
	test2()
	test3()
	test4()
	test5()
	test6()
	test7()
	test8()
	//
	text := `My email is 952479221@qq.com!`
	re := regexp.MustCompile(`[0-9a-zA-Z]+@[0-9a-zA-Z]+\.[0-9a-zA-Z.]+`)
	as := re.FindAllString(text, -1)
	fmt.Println(as)
	text = `My email is 952479221@qq.com!
His email is Shilaimu@163.com!
Her email is juiSS@gmail.com.cn! hh.end@QQ `
	re = regexp.MustCompile(`([0-9a-zA-Z]+)@([0-9a-zA-Z]+)(\.[0-9a-zA-Z.]+)`)
	ass := re.FindAllStringSubmatch(text, -1)
	fmt.Println(ass)
}

// 示例：MatchString、QuoteMeta
func test1() {
	pat := `(((abc.)def.)ghi)`
	src := `abc-def-ghi abc+def+ghi`

	fmt.Println(regexp.MatchString(pat, src))
	// true <nil>

	fmt.Println(regexp.QuoteMeta(pat))
	// \(\(\(abc\.\)def\.\)ghi\)
}

// 示例：第一匹配和最长匹配
func test2() {
	b := []byte("abc1def1")
	pat := `abc1|abc1def1`
	reg1 := regexp.MustCompile(pat)      // 第一匹配
	reg2 := regexp.MustCompilePOSIX(pat) // 最长匹配
	fmt.Printf("%s\n", reg1.Find(b))     // abc1
	fmt.Printf("%s\n", reg2.Find(b))     // abc1def1

	b = []byte("abc1def1")
	pat = `(abc|abc1def)*1`
	reg1 = regexp.MustCompile(pat)      // 第一匹配
	reg2 = regexp.MustCompilePOSIX(pat) // 最长匹配
	fmt.Printf("%s\n", reg1.Find(b))    // abc1
	fmt.Printf("%s\n", reg2.Find(b))    // abc1def1
}

// 示例：正则信息
func test3() {
	pat := `(abc)(def)(ghi)`
	reg := regexp.MustCompile(pat)

	// 获取正则表达式字符串
	fmt.Println(reg.String()) // (abc)(def)(ghi)

	// 获取分组数量
	fmt.Println(reg.NumSubexp()) // 3

	fmt.Println()

	// 获取分组名称
	pat = `(?P<Name1>abc)(def)(?P<Name3>ghi)`
	reg = regexp.MustCompile(pat)

	for i := 0; i <= reg.NumSubexp(); i++ {
		fmt.Printf("%d: %q\n", i, reg.SubexpNames()[i])
	}
	// 0: ""
	// 1: "Name1"
	// 2: ""
	// 3: "Name3"

	fmt.Println()

	// 获取字面前缀
	pat = `(abc1)(abc2)(abc3)`
	reg = regexp.MustCompile(pat)
	fmt.Println(reg.LiteralPrefix()) // abc1abc2abc3 true

	pat = `(abc1)|(abc2)|(abc3)`
	reg = regexp.MustCompile(pat)
	fmt.Println(reg.LiteralPrefix()) // false

	pat = `abc1|abc2|abc3`
	reg = regexp.MustCompile(pat)
	fmt.Println(reg.LiteralPrefix()) // abc false
}

// 示例：Find、FindSubmatch
func test4() {
	pat := `(((abc.)def.)ghi)`
	reg := regexp.MustCompile(pat)

	src := []byte(`abc-def-ghi abc+def+ghi`)

	// 查找第一个匹配结果
	fmt.Printf("%s\n", reg.Find(src)) // abc-def-ghi

	fmt.Println()

	// 查找第一个匹配结果及其分组字符串
	first := reg.FindSubmatch(src)
	for i := 0; i < len(first); i++ {
		fmt.Printf("%d: %s\n", i, first[i])
	}
	// 0: abc-def-ghi
	// 1: abc-def-ghi
	// 2: abc-def-
	// 3: abc-
}

// 示例：FindIndex、FindSubmatchIndex
func test5() {
	pat := `(((abc.)def.)ghi)`
	reg := regexp.MustCompile(pat)

	src := []byte(`abc-def-ghi abc+def+ghi`)

	// 查找第一个匹配结果
	matched := reg.FindIndex(src)
	fmt.Printf("%v\n", matched) // [0 11]
	m := matched[0]
	n := matched[1]
	fmt.Printf("%s\n\n", src[m:n]) // abc-def-ghi

	// 查找第一个匹配结果及其分组字符串
	matched = reg.FindSubmatchIndex(src)
	fmt.Printf("%v\n", matched) // [0 11 0 11 0 8 0 4]
	for i := 0; i < len(matched)/2; i++ {
		m := matched[i*2]
		n := matched[i*2+1]
		fmt.Printf("%s\n", src[m:n])
	}
	// abc-def-ghi
	// abc-def-ghi
	// abc-def-
	// abc-
}

// 示例：FindAll、FindAllSubmatch
func test6() {
	pat := `(((abc.)def.)ghi)`
	reg := regexp.MustCompile(pat)

	s := []byte(`abc-def-ghi abc+def+ghi`)

	// 查找所有匹配结果
	for _, one := range reg.FindAll(s, -1) {
		fmt.Printf("%s\n", one)
	}
	// abc-def-ghi
	// abc+def+ghi

	// 查找所有匹配结果及其分组字符串
	all := reg.FindAllSubmatch(s, -1)
	for i := 0; i < len(all); i++ {
		fmt.Println()
		one := all[i]
		for i := 0; i < len(one); i++ {
			fmt.Printf("%d: %s\n", i, one[i])
		}
	}
	// 0: abc-def-ghi
	// 1: abc-def-ghi
	// 2: abc-def-
	// 3: abc-

	// 0: abc+def+ghi
	// 1: abc+def+ghi
	// 2: abc+def+
	// 3: abc+
}

// 示例：Expand
func test7() {
	pat := `(((abc.)def.)ghi)`
	reg := regexp.MustCompile(pat)

	src := []byte(`abc-def-ghi abc+def+ghi`)
	template := []byte(`$0 $1 $2 $3`)

	// 替换第一次匹配结果
	match := reg.FindSubmatchIndex(src)
	fmt.Printf("%v\n", match) // [0 11 0 11 0 8 0 4]
	dst := reg.Expand(nil, template, src, match)
	fmt.Printf("%s\n\n", dst)
	// abc-def-ghi abc-def-ghi abc-def- abc-

	// 替换所有匹配结果
	for _, match := range reg.FindAllSubmatchIndex(src, -1) {
		fmt.Printf("%v\n", match)
		dst := reg.Expand(nil, template, src, match)
		fmt.Printf("%s\n", dst)
	}
	// [0 11 0 11 0 8 0 4]
	// abc-def-ghi abc-def-ghi abc-def- abc-
	// [12 23 12 23 12 20 12 16]
	// abc+def+ghi abc+def+ghi abc+def+ abc+
}

func test8() {
	text := `Hello 世界！123 Go.`

	// 查找连续的小写字母
	reg := regexp.MustCompile(`[a-z]+`)
	fmt.Printf("%q\n", reg.FindAllString(text, -1))
	// ["ello" "o"]

	// 查找连续的非小写字母
	reg = regexp.MustCompile(`[^a-z]+`)
	fmt.Printf("%q\n", reg.FindAllString(text, -1))
	// ["H" " 世界！123 G" "."]

	// 查找连续的单词字母
	reg = regexp.MustCompile(`[\w]+`)
	fmt.Printf("%q\n", reg.FindAllString(text, -1))
	// ["Hello" "123" "Go"]

	// 查找连续的非单词字母、非空白字符
	reg = regexp.MustCompile(`[^\w\s]+`)
	fmt.Printf("%q\n", reg.FindAllString(text, -1))
	// ["世界！" "."]

	// 查找连续的大写字母
	reg = regexp.MustCompile(`[[:upper:]]+`)
	fmt.Printf("%q\n", reg.FindAllString(text, -1))
	// ["H" "G"]

	// 查找连续的非 ASCII 字符
	reg = regexp.MustCompile(`[[:^ascii:]]+`)
	fmt.Printf("%q\n", reg.FindAllString(text, -1))
	// ["世界！"]

	// 查找连续的标点符号
	reg = regexp.MustCompile(`[\pP]+`)
	fmt.Printf("%q\n", reg.FindAllString(text, -1))
	// ["！" "."]

	// 查找连续的非标点符号字符
	reg = regexp.MustCompile(`[\PP]+`)
	fmt.Printf("%q\n", reg.FindAllString(text, -1))
	// ["Hello 世界" "123 Go"]

	// 查找连续的汉字
	reg = regexp.MustCompile(`[\p{Han}]+`)
	fmt.Printf("%q\n", reg.FindAllString(text, -1))
	// ["世界"]

	// 查找连续的非汉字字符
	reg = regexp.MustCompile(`[\P{Han}]+`)
	fmt.Printf("%q\n", reg.FindAllString(text, -1))
	// ["Hello " "！123 Go."]

	// 查找 Hello 或 Go
	reg = regexp.MustCompile(`Hello|Go`)
	fmt.Printf("%q\n", reg.FindAllString(text, -1))
	// ["Hello" "Go"]

	// 查找行首以 H 开头，以空格结尾的字符串
	reg = regexp.MustCompile(`^H.*\s`)
	fmt.Printf("%q\n", reg.FindAllString(text, -1))
	// ["Hello 世界！123 "]

	// 查找行首以 H 开头，以空白结尾的字符串（非贪婪模式）
	reg = regexp.MustCompile(`(?U)^H.*\s`)
	fmt.Printf("%q\n", reg.FindAllString(text, -1))
	// ["Hello "]

	// 查找以 hello 开头（忽略大小写），以 Go 结尾的字符串
	reg = regexp.MustCompile(`(?i:^hello).*Go`)
	fmt.Printf("%q\n", reg.FindAllString(text, -1))
	// ["Hello 世界！123 Go"]

	// 查找 Go.
	reg = regexp.MustCompile(`\QGo.\E`)
	fmt.Printf("%q\n", reg.FindAllString(text, -1))
	// ["Go."]

	// 查找从行首开始，以空格结尾的字符串（非贪婪模式）
	reg = regexp.MustCompile(`(?U)^.* `)
	fmt.Printf("%q\n", reg.FindAllString(text, -1))
	// ["Hello "]

	// 查找以空格开头，到行尾结束，中间不包含空格字符串
	reg = regexp.MustCompile(` [^ ]*$`)
	fmt.Printf("%q\n", reg.FindAllString(text, -1))
	// [" Go."]

	// 查找“单词边界”之间的字符串
	reg = regexp.MustCompile(`(?U)\b.+\b`)
	fmt.Printf("%q\n", reg.FindAllString(text, -1))
	// ["Hello" " 世界！" "123" " " "Go"]

	// 查找连续 1 次到 4 次的非空格字符，并以 o 结尾的字符串
	reg = regexp.MustCompile(`[^ ]{1,4}o`)
	fmt.Printf("%q\n", reg.FindAllString(text, -1))
	// ["Hello" "Go"]

	// 查找 Hello 或 Go
	reg = regexp.MustCompile(`(?:Hell|G)o`)
	fmt.Printf("%q\n", reg.FindAllString(text, -1))
	// ["Hello" "Go"]

	// 查找 Hello 或 Go，替换为 Hellooo、Gooo
	reg = regexp.MustCompile(`(?:Hell|G)o`)
	fmt.Printf("%q\n", reg.ReplaceAllString(text, "${n}ooo"))
	// "Hellooo 世界！123 Gooo."

	// 交换 Hello 和 Go
	reg = regexp.MustCompile(`(Hello)(.*)(Go)`)
	fmt.Printf("%q\n", reg.ReplaceAllString(text, "$3$2$1"))
	// "Go 世界！123 Hello."

	// 特殊字符的查找
	reg = regexp.MustCompile(`[\f\t\n\r\v\123\x7F\x{10FFFF}\\\^\$\.\*\+\?\{\}\(\)\[\]\|]`)
	fmt.Printf("%q\n", reg.ReplaceAllString("\f\t\n\r\v\123\x7F\U0010FFFF\\^$.*+?{}()[]|", "-"))
	// "----------------------"
}
