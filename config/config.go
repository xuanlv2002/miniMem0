package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

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

// VectorConfig 定义向量数据库的配置结构
type VectorConfig struct {
	Path                string  `mapstructure:"PATH"`
	Collection          string  `mapstructure:"COLLECTION_NAME"`
	MaxTopK             int     `mapstructure:"MAX_TOPK"`
	SimilarityThreshold float32 `mapstructure:"MIN_SIMILARITY_THRESHOLD"`
}

// Config 定义整个配置结构
type Config struct {
	ChatConfig      *LLMConfig       `mapstructure:"LLM"`
	EmbeddingConfig *EmbeddingConfig `mapstructure:"EMBEDDING"`
	VectorConfig    *VectorConfig    `mapstructure:"VECTOR_DB"`
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

// GetEmbeddingConfig 获取 Embedding 配置
func (c *Config) GetVectorConfig() *VectorConfig {
	return c.VectorConfig
}

func (c *Config) String() string {
	var sb strings.Builder

	sb.WriteString("Config:\n")

	if c.ChatConfig != nil {
		sb.WriteString("  LLM Configuration:\n")
		sb.WriteString(fmt.Sprintf("    Model: %s\n", c.ChatConfig.Model))
		sb.WriteString(fmt.Sprintf("    BaseURL: %s\n", c.ChatConfig.BaseURL))
		sb.WriteString("    APIKey: [REDACTED]\n")
		sb.WriteString(fmt.Sprintf("    Temperature: %.2f\n", c.ChatConfig.Temperature))
	} else {
		sb.WriteString("  LLM Configuration: nil\n")
	}

	if c.EmbeddingConfig != nil {
		sb.WriteString("  Embedding Configuration:\n")
		sb.WriteString(fmt.Sprintf("    Model: %s\n", c.EmbeddingConfig.Model))
		sb.WriteString(fmt.Sprintf("    BaseURL: %s\n", c.EmbeddingConfig.BaseURL))
		sb.WriteString("    APIKey: [REDACTED]\n")
		sb.WriteString(fmt.Sprintf("    Dimensions: %d\n", c.EmbeddingConfig.Dimensions))
	} else {
		sb.WriteString("  Embedding Configuration: nil\n")
	}

	if c.VectorConfig != nil {
		sb.WriteString("  Vector Database Configuration:\n")
		sb.WriteString(fmt.Sprintf("    Path: %s\n", c.VectorConfig.Path))
		sb.WriteString(fmt.Sprintf("    Collection: %s\n", c.VectorConfig.Collection))
		sb.WriteString(fmt.Sprintf("    MaxTopK: %d\n", c.VectorConfig.MaxTopK))
		sb.WriteString(fmt.Sprintf("    SimilarityThreshold: %.2f\n", c.VectorConfig.SimilarityThreshold))
	} else {
		sb.WriteString("  Vector Database Configuration: nil\n")
	}

	return sb.String()
}
