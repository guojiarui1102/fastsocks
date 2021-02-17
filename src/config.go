package src

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
)

var (
	configPath string
)

type Config struct {
	ListenAddr string `json:"listen"`
	RemoteAddr string `json:"remote"`
	Password   string `json:"password"`
}

func init() {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println("找不到home目录")
	}

	configFilename := ".fastsocks.json"
	if len(os.Args) == 2 {
		configPath = path.Join(home, configFilename)
	}
}

func (config *Config) SaveConfig() {
	configJson, _ := json.MarshalIndent(config, "", "      ")
	err := ioutil.WriteFile(configPath, configJson, 0644)
	if err != nil {
		fmt.Printf("保存配置到文件 %s 出错: %s", configPath, err)
	}
	log.Printf("保存配置到文件 %s 成功\n", configPath)
}

func (config *Config) ReadConfig() {
	if _, err := os.Stat(configPath); !os.IsNotExist(err) {
		log.Printf("从文件 %s 中读取配置\n", configPath)
		file, err := os.Open(configPath)
		if err != nil {
			log.Fatalf("打开配置文件 %s 出错: %s", configPath, err)
		}
		defer file.Close()

		err = json.NewDecoder(file).Decode(config)
		if err != nil {
			log.Fatalf("格式不合法的JSON配置文件: %s", file.Name())
		}
	}
}
