package cmd

import (
	"crypto/sha256"
	"crypto/tls"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func httpPostForm(url string, data url.Values) (body []byte, err error) {
	tr := &http.Transport{
		//如果需要测试自签名的证书 这里需要设置跳过证书检测 否则编译报错
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.PostForm(url, data)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	return
}

func getSign(q, appKey, salt, curtime string) string {
	// sign=sha256(应用ID+input+salt+curtime+应用密钥)
	// 其中，input的计算方式为：input=q前10个字符 + q长度 + q后10个字符（当q长度大于20）或 input=q字符串（当q长度小于等于20）
	s := appKey + getInputStr(q) + salt + curtime + secretKey

	return getSha256(s)
}

func getInputStr(q string) string {
	inputSlice := []rune(q)
	if len(inputSlice) >= 20 {
		inputStart10 := string(inputSlice[:10])
		inputEnd10 := string(inputSlice[len(inputSlice)-10:])
		input := inputStart10 + strconv.Itoa(len(inputSlice)) + inputEnd10
		return input
	}
	return q
}

func getSha256(s string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(s)))
}

func getUUID() string {
	return uuid.New().String()
}

func getTranslateRequestData(q, appKey string) url.Values {
	salt := getUUID()
	utime := time.Now().Unix()
	curtime := fmt.Sprintf("%v", utime)

	data := make(url.Values)

	data["q"] = []string{q}
	data["from"] = []string{"en"}
	data["to"] = []string{"zh-CHS"}
	data["appKey"] = []string{appKey}
	data["salt"] = []string{salt}
	data["sign"] = []string{getSign(q, appKey, salt, curtime)}
	data["signType"] = []string{"v3"}
	data["curtime"] = []string{curtime}

	return data
}
