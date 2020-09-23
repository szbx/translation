package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
	"os"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

const apiURL = "https://openapi.youdao.com/api"

var (
	secretKey = "bVsWxOXsUnLHtZjxNEE8tCjpA3pYgIyF"
	appKey    = "17f90495a0c492e6"
)

var rootCmd = &cobra.Command{
	Use:   "fanyi",
	Short: "简单命令行翻译工具",
	Long: `简单命令行翻译工具
	使用有道翻译
	- 没有指定参数时翻译剪切版内容
	- 指定参数翻译命令行后的文本内容`,
	Run: func(cmd *cobra.Command, args []string) {
		defer handleError()
		queryStr, err := getQueryStr(args)
		if err != nil {
			panic(err)
		}
		if queryStr == "" {
			panic(errors.New("请输入翻译内容"))
		}
		resp, err := translate(queryStr)
		if err != nil {
			panic(err)
		}
		resp.print()
	},
}

func handleError() {
	err := recover()
	if err != nil {
		fmt.Printf("出错了:%v\n", err)
	}
}

func getQueryStr(args []string) (q string, err error) {
	if len(args) > 0 {
		q = strings.Join(args, " ")
	} else {
		q, err = clipboard.ReadAll()
	}
	return
}

func translate(q string) (r translateResp, err error) {

	body, err := httpPostForm(apiURL, getTranslateRequestData(q, appKey))
	if err != nil {
		return
	}

	json.Unmarshal(body, &r)
	return
}

type translateResp struct {
	ErrorCode   string   `json:"errorCode"`
	Query       string   `json:"query"`
	Translation []string `json:"translation"`
}

func (tr *translateResp) print() {
	if tr.ErrorCode != "0" {
		panic(errors.New("翻译失败,错误码:" + tr.ErrorCode))
	}
	for _, s := range tr.Translation {
		fmt.Printf("%v\n", s)
	}
}

// ---------------------------------------------------------------------------------

var cfgFile string

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
	/*
		if err := viper.ReadInConfig(); err == nil {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		}
	*/
}
