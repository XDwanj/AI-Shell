package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"AI-Shell/internal/config"
	"AI-Shell/internal/openai"
	"AI-Shell/internal/system"

	"github.com/spf13/cobra"
)

// AIResponse 表示AI返回的命令选项
type AIResponse struct {
	Command []string `json:"command"`
	Msg     string   `json:"msg"`
	Code    int      `json:"code"`
}

var executeCmd = &cobra.Command{
	Use:   "exec [description]",
	Short: "执行自然语言命令",
	Long:  `将自然语言描述转换为命令行命令并执行。`,
	Args:  cobra.MinimumNArgs(1),
	RunE:  runExecute,
}

func init() {
	rootCmd.AddCommand(executeCmd)
}

func runExecute(cmd *cobra.Command, args []string) error {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("加载配置失败: %v", err)
	}

	// 获取系统信息
	sysInfo, err := system.GetSystemInfo()
	if err != nil {
		return fmt.Errorf("获取系统信息失败: %v", err)
	}

	// 获取show-data标志
	showData, err := cmd.Flags().GetBool("show-data")
	if err != nil {
		return fmt.Errorf("获取show-data标志失败: %v", err)
	}

	// 创建OpenAI客户端
	client := openai.NewClient(cfg)

	// 准备系统提示
	systemPrompt := "你是一个命令行命令翻译机，负责将用户输入翻译为命令行命令，你需要以json方式回复，以下是示例\n" +
		"{\"command\": [\"ls\"],\"msg\": \"执行此命令将列出当前目录中的文件和子目录。\",\"code\": 0}\n" +
		"command是可执行命令，可以有多种翻译结果，每一项都是完整的命令，不要把一条命令拆分为开，用户选择其中一条执行，最多为10个，" +
		"msg是展示给用户的提示信息，code为翻译结果，0为成功翻译，1为不能翻译、缺少信息或其他异常情况。"

	// 构建用户提示（包含系统信息）
	userPrompt := sysInfo + "\n" + args[0]

	// 发送请求到OpenAI
	resp, err := client.SendRequest(systemPrompt, userPrompt)
	if err != nil {
		return fmt.Errorf("发送请求失败: %v", err)
	}

	if len(resp.Choices) == 0 {
		return fmt.Errorf("未收到有效响应")
	}

	// 解析响应
	var aiResp AIResponse
	content := resp.Choices[0].Message.Content

	// 如果响应是Markdown json 格式，去除Markdown标记
	if len(content) > 7 && content[:7] == "```json" {
		content = content[7:]
	}
	// 如果响应是Markdown js 格式，去除Markdown标记
	if len(content) > 5 && content[:5] == "```js" {
		content = content[5:]
	}
	if len(content) > 3 && content[len(content)-3:] == "```" {
		content = content[:len(content)-3]
	}

	if err := json.Unmarshal([]byte(content), &aiResp); err != nil {
		return fmt.Errorf("解析响应失败: %v", err)
	}

	// 如果showData为true，显示发送的数据
	if showData {
		fmt.Print("响应数据: ")
		data, err := json.MarshalIndent(aiResp, "", "  ")
		if err != nil {
			return fmt.Errorf("格式化响应数据失败: %v", err)
		}
		fmt.Println(string(data))
	}

	// 输出提示信息
	fmt.Println(aiResp.Msg)
	fmt.Println("---------------------")
	fmt.Println("可用的命令选项:")

	// 检查翻译结果
	if aiResp.Code != 0 {
		return fmt.Errorf("命令翻译失败")
	}

	// 显示可用的命令选项
	for i, cmd := range aiResp.Command {
		fmt.Printf("%d: %s\n", i+1, cmd)
	}
	fmt.Println("0: 退出")

	// 获取用户选择
	fmt.Print("请选择要执行的命令: ")
	var choice string
	fmt.Scanln(&choice)

	// 解析用户选择
	num, err := strconv.Atoi(choice)
	if err != nil {
		return fmt.Errorf("无效的选择")
	}

	if num == 0 {
		fmt.Println("退出程序。")
		return nil
	}

	if num < 1 || num > len(aiResp.Command) {
		return fmt.Errorf("无效的选择")
	}

	// 执行选中的命令
	selectedCmd := aiResp.Command[num-1]
	fmt.Printf("执行命令: %s\n", selectedCmd)
	fmt.Println("---------------------")

	// 设置环境变量
	env := os.Environ()
	env = append(env, "TERM=xterm-256color")

	// 创建命令
	command := exec.Command("bash", "-i", "-c", selectedCmd)
	command.Env = env
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin

	// 执行命令
	if err := command.Run(); err != nil {
		return fmt.Errorf("命令执行失败: %v", err)
	}

	return nil
}
