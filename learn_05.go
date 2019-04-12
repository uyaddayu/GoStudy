package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"
)

// 看看加密效率
func main() {
	str := `{"app":"ad","cid":"F7BDF3CB4985413CB42E229DD1E96918","dhtY1E286Axcpw":"7zdW6vO19p1B","fission":"OD7FSu30ay","itemid":"5A3DFA122A454584A7AC5D3320C9C61A","nowt":1553512175,"uid":"1F7C503FF0114773919B0D5703E7FBFF"}`
	z1, err := Encrypt(str, []byte("auukiis0"))
	fmt.Println(err)
	t1 := time.Now()
	b := []byte("auukiis0")
	for i := 0; i < 100000; i += 1 {
		str2, err := Decrypt(z1, b)
		fmt.Println(str2, err)
	}
	t2 := time.Now()
	sub1 := t2.Sub(t1).Seconds()
	fmt.Println("2-->")
	// 第二种写法
	z2 := AESEncodeStr(str, "auukiis0auukiis0auukiis0auukiis0")
	t1 = time.Now()
	for i := 0; i < 100000; i += 1 {
		str2 := AESDecodeStr(z2, "auukiis0auukiis0auukiis0auukiis0")
		fmt.Println(str2)
	}
	t2 = time.Now()
	sub2 := t2.Sub(t1).Seconds()
	fmt.Println(z1, "耗时(s)", sub1)
	fmt.Println(z2, "耗时(s)", sub2)
}

func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}

func ZeroUnPadding(origData []byte) []byte {
	return bytes.TrimFunc(origData,
		func(r rune) bool {
			return r == rune(0)
		})
}
func Encrypt(text string, key []byte) (string, error) {
	src := []byte(text)
	block, err := des.NewCipher(key)
	if err != nil {
		return "", err
	}
	bs := block.BlockSize()
	src = ZeroPadding(src, bs)
	if len(src)%bs != 0 {
		return "", errors.New("Need a multiple of the blocksize")
	}
	out := make([]byte, len(src))
	dst := out
	for len(src) > 0 {
		block.Encrypt(dst, src[:bs])
		src = src[bs:]
		dst = dst[bs:]
	}
	return hex.EncodeToString(out), nil
}

func Decrypt(decrypted string, key []byte) (string, error) {
	src, err := hex.DecodeString(decrypted)
	if err != nil {
		return "", err
	}
	block, err := des.NewCipher(key)
	if err != nil {
		return "", err
	}
	out := make([]byte, len(src))
	dst := out
	bs := block.BlockSize()
	if len(src)%bs != 0 {
		return "", errors.New("crypto/cipher: input not full blocks")
	}
	for len(src) > 0 {
		block.Decrypt(dst, src[:bs])
		src = src[bs:]
		dst = dst[bs:]
	}
	out = ZeroUnPadding(out)
	return string(out), nil
}

var ivspec = []byte("0000000000000000")

func AESEncodeStr(src, key string) string {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		fmt.Println("key error1", err)
		return ""
	}
	if src == "" {
		fmt.Println("plain content empty")
		return ""
	}
	ecb := cipher.NewCBCEncrypter(block, ivspec)
	content := []byte(src)
	content = PKCS5Padding(content, block.BlockSize())
	crypted := make([]byte, len(content))
	ecb.CryptBlocks(crypted, content)
	return hex.EncodeToString(crypted)
}

func AESDecodeStr(crypt, key string) string {
	crypted, err := hex.DecodeString(strings.ToLower(crypt))
	if err != nil || len(crypted) == 0 {
		fmt.Println("plain content empty")
	}
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		fmt.Println("key error1", err)
	}
	ecb := cipher.NewCBCDecrypter(block, ivspec)
	decrypted := make([]byte, len(crypted))
	ecb.CryptBlocks(decrypted, crypted)

	return string(PKCS5Trimming(decrypted))
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5Trimming(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]
	return encrypt[:len(encrypt)-int(padding)]
}
