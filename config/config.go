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

/* 基础组件配置 */
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
	TopK                int     `mapstructure:"TOPK"`
	SimilarityThreshold float32 `mapstructure:"SIMILARITY_THRESHOLD"`
}

type SqlConfig struct {
	Path string `mapstructure:"PATH"`
}

/* 记忆层配置 */
// MemoryContextConfig 定义记忆上下文的配置
type ContextMemoryConfig struct {
	SummaryGap int `mapstructure:"SUMMARY_GAP"`
}

// LongMemoryConfig 定义长记忆的配置
type LongMemoryConfig struct {
	LongGap int `mapstructure:"LONG_GAP"`
}

// ShortMemoryConfig 定义短记忆的配置

type ShortMemoryConfig struct {
	ShortWindow int `mapstructure:"SHORT_WINDOW"`
}

// Config 管理所有配置
type Config struct {
	ChatConfig          *LLMConfig           `mapstructure:"LLM"`
	EmbeddingConfig     *EmbeddingConfig     `mapstructure:"EMBEDDING"`
	VectorConfig        *VectorConfig        `mapstructure:"VECTOR_DB"`
	SqlConfig           *SqlConfig           `mapstructure:"SQL_DB"`
	MemoryContextConfig *ContextMemoryConfig `mapstructure:"CONTEXT_MEMORY"`
	LongMemoryConfig    *LongMemoryConfig    `mapstructure:"LONG_MEMORY"`
	ShortMemoryConfig   *ShortMemoryConfig   `mapstructure:"SHORT_MEMORY"`
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

// GetSqlConfig 获取 Sql 配置
func (c *Config) GetSqlConfig() *SqlConfig {
	return c.SqlConfig
}

// GetMemoryContextConfig 获取 MemoryContext 配置
func (c *Config) GetMemoryContextConfig() *ContextMemoryConfig {
	return c.MemoryContextConfig
}

// GetLongMemoryConfig 获取 LongMemory 配置
func (c *Config) GetLongMemoryConfig() *LongMemoryConfig {
	return c.LongMemoryConfig
}

// GetShortMemoryConfig 获取 ShortMemory 配置
func (c *Config) GetShortMemoryConfig() *ShortMemoryConfig {
	return c.ShortMemoryConfig
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
		sb.WriteString(fmt.Sprintf("    MaxTopK: %d\n", c.VectorConfig.TopK))
		sb.WriteString(fmt.Sprintf("    SimilarityThreshold: %.2f\n", c.VectorConfig.SimilarityThreshold))
	} else {
		sb.WriteString("  Vector Database Configuration: nil\n")
	}

	if c.SqlConfig != nil {
		sb.WriteString("  SQL Database Configuration:\n")
		sb.WriteString(fmt.Sprintf("    Path: %s\n", c.SqlConfig.Path))
	} else {
		sb.WriteString("  SQL Database Configuration: nil\n")
	}

	if c.MemoryContextConfig != nil {
		sb.WriteString("  Memory Context Configuration:\n")
		sb.WriteString(fmt.Sprintf("    SummaryGap: %d\n", c.MemoryContextConfig.SummaryGap))

	} else {
		sb.WriteString("  Memory Context Configuration: nil\n")
	}

	if c.LongMemoryConfig != nil {
		sb.WriteString("  Long Memory Configuration:\n")
		sb.WriteString(fmt.Sprintf("    LongGap: %d\n", c.LongMemoryConfig.LongGap))
	} else {
		sb.WriteString("  Long Memory Configuration: nil\n")
	}

	if c.ShortMemoryConfig != nil {
		sb.WriteString("  Short Memory Configuration:\n")
		sb.WriteString(fmt.Sprintf("    ShortWindow: %d\n", c.ShortMemoryConfig.ShortWindow))

	} else {
		sb.WriteString("  Short Memory Configuration: nil\n")
	}

	return sb.String()
}
