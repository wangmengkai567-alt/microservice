package config

import (
	"os"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`

	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
	} `yaml:"database"`

	JWT struct {
		Secret string `yaml:"secret"`
		Expire int    `yaml:"expire"`
	} `yaml:"jwt"`

	Services struct {
		UserService    string `yaml:"user_service"`
		ProductService string `yaml:"product_service"`
	} `yaml:"services"`
}

var Global Config

func Init() error {
	data, err := os.ReadFile("config/config.yaml")
	if err != nil {
		return err
	}
	
	if err := yaml.Unmarshal(data, &Global); err != nil {
		return err
	}

	// 环境变量覆盖配置文件
	if host := os.Getenv("DB_HOST"); host != "" {
		Global.Database.Host = host
	}
	if port := os.Getenv("DB_PORT"); port != "" {
		Global.Database.Port = 3306
	}
	if user := os.Getenv("DB_USER"); user != "" {
		Global.Database.User = user
	}
	if password := os.Getenv("DB_PASSWORD"); password != "" {
		Global.Database.Password = password
	}
	if dbname := os.Getenv("DB_NAME"); dbname != "" {
		Global.Database.DBName = dbname
	}
	if userService := os.Getenv("USER_SERVICE"); userService != "" {
		Global.Services.UserService = userService
	}
	if productService := os.Getenv("PRODUCT_SERVICE"); productService != "" {
		Global.Services.ProductService = productService
	}

	return nil
}
