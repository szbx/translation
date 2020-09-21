/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"crypto/sha256"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "fanyi",
	Short: "简单命令行翻译工具",
	Long: `简单命令行翻译工具
	使用有道翻译
	- 没有指定参数时翻译剪切版内容
	- 指定参数翻译命令行后的文本内容`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		defer func() {
			err := recover()
			if err != nil {
				fmt.Printf("出错了:%v\n", err)
			}
		}()
		var content string
		var err error
		if len(args) > 0 {
			content = strings.Join(args, " ")
		} else {
			// 读取剪切板中的内容到字符串
			content, err = clipboard.ReadAll()
			if err != nil {
				panic(err)
			}
		}

		// fmt.Println("content", content)
		ret, _ := transform(content)
		// fmt.Printf("%#v", ret)
		for _, s := range ret.Translation {
			fmt.Printf("%v\n", s)
		}
	},
}

const apiURL = "https://openapi.youdao.com/api"

func transform(q string) (r resp, err error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	} //如果需要测试自签名的证书 这里需要设置跳过证书检测 否则编译报错
	client := &http.Client{Transport: tr}
	data := make(url.Values)
	appKey := "17f90495a0c492e6"
	salt := "understand-reading"
	utime := time.Now().Unix()
	curtime := fmt.Sprintf("%v", utime)
	data["q"] = []string{q}
	data["from"] = []string{"en"}
	data["to"] = []string{"zh-CHS"}
	data["appKey"] = []string{appKey}
	data["salt"] = []string{salt}
	data["sign"] = []string{getSign(q, appKey, salt, curtime)}
	data["signType"] = []string{"v3"}
	data["curtime"] = []string{curtime}
	resp, err := client.PostForm(apiURL, data)

	if err != nil {
		fmt.Println("error:", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	json.Unmarshal(body, &r)

	return
}

type resp struct {
	ErrorCode   string   `json:"errorCode"`
	Query       string   `json:"query"`
	Translation []string `json:"translation"`
}

var secretKey = "bVsWxOXsUnLHtZjxNEE8tCjpA3pYgIyF"

func getSign(q, appKey, salt, curtime string) string {
	// sign=sha256(应用ID+input+salt+curtime+应用密钥)
	// 其中，input的计算方式为：input=q前10个字符 + q长度 + q后10个字符（当q长度大于20）或 input=q字符串（当q长度小于等于20）
	s := appKey + getInputStr(q) + salt + curtime + secretKey

	return getSha256(s)
}

func getInputStr(q string) string {
	inputSlice := []rune(q)
	if len(inputSlice) >= 20 {
		input10 := string(inputSlice[:10])
		inputx10 := string(inputSlice[len(inputSlice)-10:])
		input := input10 + strconv.Itoa(len(inputSlice)) + inputx10
		return input
	}
	return q
}

func getSha256(s string) string {
	sum := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", sum)
}

// ==============================================================================================================

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.fanyi.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".fanyi" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".fanyi")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
