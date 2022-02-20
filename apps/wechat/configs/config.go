package configs

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	HttpPort string `yaml:"http_port"`
	RpcPort  string `yaml:"rpc_port"`

	Wechat struct {
		Appid  string `yaml:"appid"`
		Secret string `yaml:"secret"`
	} `yaml:"wechat"`

	Log struct {
		Dir          string `yaml:"dir"`
		Suffix       string `yaml:"suffix"`
		CutDays      int    `yaml:"cut_days"`
		FileSaveDays int    `yaml:"file_save_days"`
	} `yaml:"log"`
}

func MakeFromYaml(filePath string) (*Config, error) {

	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var c = new(Config)

	if err := yaml.Unmarshal(file, c); err != nil {
		return nil, err
	}

	return c, nil
}
