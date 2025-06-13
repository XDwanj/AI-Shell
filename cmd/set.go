package cmd

import (
	"fmt"
	"strconv"

	"AI-Shell/internal/config"

	"github.com/spf13/cobra"
)

var (
	configCmd = &cobra.Command{
		Use:   "config",
		Short: "配置管理命令",
		Long:  `管理 AI-Shell 的配置，包括API URL、密钥、模型等设置。`,
	}

	setCmd = &cobra.Command{
		Use:   "set",
		Short: "设置配置项",
		Long:  `设置各种配置项的值，包括API URL、密钥、模型等。`,
	}

	viewCmd = &cobra.Command{
		Use:   "view",
		Short: "查看当前配置",
		Long:  `显示所有当前配置项的值。`,
		RunE:  runView,
	}

	// 设置各个配置项的子命令
	setURLCmd = &cobra.Command{
		Use:   "url [API_URL]",
		Short: "设置API URL",
		Long:  `设置OpenAI API的URL地址。`,
		Args:  cobra.ExactArgs(1),
		RunE:  runSetURL,
	}

	setKeyCmd = &cobra.Command{
		Use:   "key [API_KEY]",
		Short: "设置API密钥",
		Long:  `设置OpenAI API的访问密钥。`,
		Args:  cobra.ExactArgs(1),
		RunE:  runSetKey,
	}

	setModelCmd = &cobra.Command{
		Use:   "model [MODEL_NAME]",
		Short: "设置模型",
		Long:  `设置使用的AI模型名称。`,
		Args:  cobra.ExactArgs(1),
		RunE:  runSetModel,
	}

	setMaxTokensCmd = &cobra.Command{
		Use:   "max-tokens [NUMBER]",
		Short: "设置最大令牌数",
		Long:  `设置API请求的最大令牌数。`,
		Args:  cobra.ExactArgs(1),
		RunE:  runSetMaxTokens,
	}

	setTemperatureCmd = &cobra.Command{
		Use:   "temperature [NUMBER]",
		Short: "设置温度参数",
		Long:  `设置生成文本的随机性，范围从0到1。`,
		Args:  cobra.ExactArgs(1),
		RunE:  runSetTemperature,
	}

	setDebugCmd = &cobra.Command{
		Use:   "debug [true|false]",
		Short: "设置调试模式",
		Long:  `启用或禁用调试模式，调试模式下会输出更多的日志信息。`,
		Args:  cobra.ExactArgs(1),
		RunE:  runSetDebug,
	}
)

func init() {
	// 添加配置命令到根命令
	rootCmd.AddCommand(configCmd)

	// 添加子命令到config命令
	configCmd.AddCommand(setCmd, viewCmd)

	// 添加设置子命令
	setCmd.AddCommand(setURLCmd)
	setCmd.AddCommand(setKeyCmd)
	setCmd.AddCommand(setModelCmd)
	setCmd.AddCommand(setMaxTokensCmd)
	setCmd.AddCommand(setTemperatureCmd)
	setCmd.AddCommand(setDebugCmd)
}

func runView(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("加载配置失败: %v", err)
	}

	fmt.Printf("当前配置:\n")
	fmt.Printf("URL = %s\n", cfg.URL)
	fmt.Printf("API_KEY = %s\n", cfg.APIKey)
	fmt.Printf("MODEL = %s\n", cfg.Model)
	fmt.Printf("MAX_TOKENS = %d\n", cfg.MaxTokens)
	fmt.Printf("TEMPERATURE = %.1f\n", cfg.Temperature)

	return nil
}

func runSetURL(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("加载配置失败: %v", err)
	}

	if err := cfg.SetURL(args[0]); err != nil {
		return fmt.Errorf("设置URL失败: %v", err)
	}

	fmt.Printf("已设置 URL = %s\n", args[0])
	return nil
}

func runSetKey(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("加载配置失败: %v", err)
	}

	if err := cfg.SetAPIKey(args[0]); err != nil {
		return fmt.Errorf("设置API密钥失败: %v", err)
	}

	fmt.Printf("已设置 API_KEY = %s\n", args[0])
	return nil
}

func runSetModel(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("加载配置失败: %v", err)
	}

	if err := cfg.SetModel(args[0]); err != nil {
		return fmt.Errorf("设置模型失败: %v", err)
	}

	fmt.Printf("已设置 MODEL = %s\n", args[0])
	return nil
}

func runSetMaxTokens(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("加载配置失败: %v", err)
	}

	maxTokens, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("无效的最大令牌数值: %v", err)
	}

	if err := cfg.SetMaxTokens(maxTokens); err != nil {
		return fmt.Errorf("设置最大令牌数失败: %v", err)
	}

	fmt.Printf("已设置 MAX_TOKENS = %d\n", maxTokens)
	return nil
}

func runSetTemperature(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("加载配置失败: %v", err)
	}

	temp, err := strconv.ParseFloat(args[0], 64)
	if err != nil {
		return fmt.Errorf("无效的温度参数值: %v", err)
	}

	if temp < 0 || temp > 1 {
		return fmt.Errorf("温度参数必须在0到1之间")
	}

	if err := cfg.SetTemperature(temp); err != nil {
		return fmt.Errorf("设置温度参数失败: %v", err)
	}

	fmt.Printf("已设置 TEMPERATURE = %.1f\n", temp)
	return nil
}

func runSetDebug(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("加载配置失败: %v", err)
	}
	debugMode, err := strconv.ParseBool(args[0])
	if err != nil {
		return fmt.Errorf("无效的Debug模式值: %v", err)
	}

	if err := cfg.SetDebug(debugMode); err != nil {
		return fmt.Errorf("设置Debug模式失败: %v", err)
	}

	fmt.Printf("已设置 Debug 模式 = %v\n", debugMode)
	return nil
}
