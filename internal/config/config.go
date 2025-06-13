package config

import (
	"encoding/json"
	"fmt"
	"log/slog"
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
	Debug       bool    `json:"debug"`
}

const (
	DefaultURL         = "https://api.openai.com/v1/chat/completions"
	DefaultModel       = "gpt-4o-mini"
	DefaultMaxTokens   = 1000
	DefaultTemperature = 0.7
	DefaultDebug       = false // 默认不启用调试模式
)

var (
	configDir  string
	configFile string
)

func init() {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "获取用户配置目录失败: %v\n", err)
		os.Exit(1)
	}

	configDir = filepath.Join(userConfigDir, "ais")
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
		config := &Config{
			URL:         DefaultURL,
			Model:       DefaultModel,
			MaxTokens:   DefaultMaxTokens,
			Temperature: DefaultTemperature,
			Debug:       DefaultDebug,
		}
		return config, nil
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

	// 如果配置文件中存在 debug 字段且为 true，则更新日志级别
	if config.Debug {
		slog.SetLogLoggerLevel(slog.LevelDebug) // 设置全局日志级别为 Debug
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
	slog.Debug("设置配置项", "字段", "URL", "值", url)
	c.URL = url
	return c.SaveConfig()
}

// SetAPIKey 设置API密钥
func (c *Config) SetAPIKey(key string) error {
	slog.Debug("设置配置项", "字段", "APIKey", "值", key)
	c.APIKey = key
	return c.SaveConfig()
}

// SetModel 设置模型名称
func (c *Config) SetModel(model string) error {
	slog.Debug("设置配置项", "字段", "Model", "值", model)
	c.Model = model
	return c.SaveConfig()
}

// SetMaxTokens 设置最大Token数
func (c *Config) SetMaxTokens(maxTokens int) error {
	slog.Debug("设置配置项", "字段", "MaxTokens", "值", maxTokens)
	c.MaxTokens = maxTokens
	return c.SaveConfig()
}

// SetTemperature 设置温度参数
func (c *Config) SetTemperature(temperature float64) error {
	slog.Debug("设置配置项", "字段", "Temperature", "值", temperature)
	c.Temperature = temperature
	return c.SaveConfig()
}

// SetDebug 设置调试模式
func (c *Config) SetDebug(debug bool) error {
	slog.SetLogLoggerLevel(slog.LevelDebug) // 设置全局日志级别为 Debug
	slog.Debug("设置配置项", "字段", "Debug", "值", debug)
	c.Debug = debug
	return c.SaveConfig()
}
