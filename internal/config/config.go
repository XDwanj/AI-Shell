package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config 存储应用程序的配置信息
type Config struct {
	URL         string  `json:"url"`
	APIKey      string  `json:"api_key"`
	Model       string  `json:"model"`
	MaxTokens   int     `json:"max_tokens"`
	Temperature float64 `json:"temperature"`
}

const (
	DefaultURL         = "https://api.openai.com/v1/chat/completions"
	DefaultModel       = "gpt-4o-mini"
	DefaultMaxTokens   = 1000
	DefaultTemperature = 0.7
)

var (
	configDir  string
	configFile string
)

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "获取用户主目录失败: %v\n", err)
		os.Exit(1)
	}

	configDir = filepath.Join(homeDir, ".config")
	configFile = filepath.Join(configDir, "ais_config.json")
}

// LoadConfig 加载配置文件，如果文件不存在则创建默认配置
func LoadConfig() (*Config, error) {
	// 确保配置目录存在
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("创建配置目录失败: %v", err)
	}

	// 如果配置文件不存在，创建默认配置
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return &Config{
			URL:         DefaultURL,
			Model:       DefaultModel,
			MaxTokens:   DefaultMaxTokens,
			Temperature: DefaultTemperature,
		}, nil
	}

	// 读取配置文件
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	return &config, nil
}

// SaveConfig 保存配置到文件
func (c *Config) SaveConfig() error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失败: %v", err)
	}

	if err := os.WriteFile(configFile, data, 0644); err != nil {
		return fmt.Errorf("保存配置文件失败: %v", err)
	}

	return nil
}

// SetURL 设置API URL
func (c *Config) SetURL(url string) error {
	c.URL = url
	return c.SaveConfig()
}

// SetAPIKey 设置API密钥
func (c *Config) SetAPIKey(key string) error {
	c.APIKey = key
	return c.SaveConfig()
}

// SetModel 设置模型名称
func (c *Config) SetModel(model string) error {
	c.Model = model
	return c.SaveConfig()
}

// SetMaxTokens 设置最大Token数
func (c *Config) SetMaxTokens(maxTokens int) error {
	c.MaxTokens = maxTokens
	return c.SaveConfig()
}

// SetTemperature 设置温度参数
func (c *Config) SetTemperature(temperature float64) error {
	c.Temperature = temperature
	return c.SaveConfig()
}
