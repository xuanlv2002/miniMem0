package config

import (
	"errors"
	"os"

	"github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// LLMConfig 定义 LLM 服务的配置结构
type LLMConfig struct {
	Model       string  `mapstructure:"MODEL"`
	BaseURL     string  `mapstructure:"BASE_URL"`
	APIKey      string  `mapstructure:"API_KEY"`
	Temperature float32 `mapstructure:"TEMPERATURE"`
}

// LLMConfig 定义 LLM 服务的配置结构
type EmbeddingConfig struct {
	Model      openai.EmbeddingModel `mapstructure:"MODEL"`
	BaseURL    string                `mapstructure:"BASE_URL"`
	APIKey     string                `mapstructure:"API_KEY"`
	Dimensions int                   `mapstructure:"DIMENSIONS"`
}

// Config 定义整个配置结构
type Config struct {
	ChatConfig      *LLMConfig       `mapstructure:"LLM"`
	EmbeddingConfig *EmbeddingConfig `mapstructure:"EMBEDDING"`
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

// LoadConfig 从指定的配置文件中加载配置
func LoadConfig(configFilePath string) (*Config, error) {
	// 实现从配置文件加载配置的逻辑
	// 返回加载后的配置结构和可能的错误
	if configFilePath == "" {
		return nil, errors.New("no config file")
	}

	if !fileExists(configFilePath) {
		return nil, errors.New("no config file exists")
	}
	// 加载配置文件
	viper.SetConfigFile(configFilePath)
	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	logrus.Infof("load target config file success: %s", configFilePath)
	// 解析配置文件
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// 返回一个空的配置结构,用于用户自定义配置
func NewConfig() (*Config, error) {
	return &Config{}, nil
}

// GetLLMConfig 获取 LLM 配置
func (c *Config) GetChatConfig() *LLMConfig {
	return c.ChatConfig
}

// GetEmbeddingConfig 获取 Embedding 配置
func (c *Config) GetEmbeddingConfig() *EmbeddingConfig {
	return c.EmbeddingConfig
}
