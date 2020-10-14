package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config struct {
	RabbitMq RabbitMq `json:"rabbitMq"`
}

type RabbitMq struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Queue    Queue  `json:"queue"`
}

type Queue struct {
	Name string `json:"name"`
}

func InitConfig(path string) Config {
	var config Config
	jsonFile, _ := os.Open(path)
	byteValue, _ := ioutil.ReadAll(jsonFile)

	_ = json.Unmarshal(byteValue, &config)

	defer jsonFile.Close()
	return config
}
