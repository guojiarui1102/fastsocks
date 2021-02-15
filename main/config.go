package main

import (
	"fmt"
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
		fmt.Errorf("找不到home目录")
	}

	configFilename := ".fastsocks.json"
	if len(os.Args) == 2 {
		configPath = path.Join(home, configFilename)
	}
}
