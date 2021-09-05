package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Database Database `json:"database"`
	Http     Http     `json:"http"`
	Provider Provider `json:"provider"`
	Grpc     Grpc     `json:"grpc"`
	Sports   []Sport  `json:"sports"`
	Log      Log      `json:"log"`
}

type Database struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Name     string `json:"name"`
}

type Http struct {
	Adress string `json:"adress"`
}

type Provider struct {
	Adress string `json:"adress"`
}

type Grpc struct {
	Adress string `json:"adress"`
}

type Sport struct {
	Sport string `json:"sport"`
	N     string `json:"duration"`
}

type Log struct {
	EnableConsole     bool   `json:"enable_console"`
	ConsoleJSONFormat bool   `json:"console_JSON_format"`
	ConsoleLevel      string `json:"console_level"`
	EnableFile        bool   `json:"enable_file"`
	FileJSONFormat    bool   `json:"file_json_format"`
	FileLevel         string `json:"file_level"`
	FileLocation      string `json:"file_location"`
}

func LoadConfiguration(file string) (config Config, err error) {
	configFile, err := os.Open(file)
	if err != nil {
		return
	}

	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	if err != nil {
		return
	}

	return
}
