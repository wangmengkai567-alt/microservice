package config

import (
	"os"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Port int    `yaml:"port"`
		Mode string `yaml:"mode"`
	} `yaml:"server"`

	JWT struct {
		Secret string `yaml:"secret"`
	} `yaml:"jwt"`

	Services struct {
		UserService    string `yaml:"user_service"`
		ProductService string `yaml:"product_service"`
		OrderService   string `yaml:"order_service"`
	} `yaml:"services"`

	CORS struct {
		AllowOrigins []string `yaml:"allow_origins"`
		AllowMethods []string `yaml:"allow_methods"`
		AllowHeaders []string `yaml:"allow_headers"`
	} `yaml:"cors"`
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
	if userService := os.Getenv("USER_SERVICE"); userService != "" {
		Global.Services.UserService = userService
	}
	if productService := os.Getenv("PRODUCT_SERVICE"); productService != "" {
		Global.Services.ProductService = productService
	}
	if orderService := os.Getenv("ORDER_SERVICE"); orderService != "" {
		Global.Services.OrderService = orderService
	}

	return nil
}
