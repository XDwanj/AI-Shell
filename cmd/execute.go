package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"log/slog"

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
	slog.Debug("开始执行 runExecute", "args", args)
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("加载配置失败", "error", err)
		return fmt.Errorf("加载配置失败: %v", err)
	}
	slog.Debug("配置加载成功", "config", cfg)

	// 获取系统信息
	sysInfo, err := system.GetSystemInfo()
	if err != nil {
		slog.Error("获取系统信息失败", "error", err)
		return fmt.Errorf("获取系统信息失败: %v", err)
	}
	slog.Debug("系统信息获取成功", "sysInfo", sysInfo)

	// 获取show-data标志
	showData, err := cmd.Flags().GetBool("show-data")
	if err != nil {
		slog.Error("获取show-data标志失败", "error", err)
		return fmt.Errorf("获取show-data标志失败: %v", err)
	}
	slog.Debug("获取show-data标志", "showData", showData)

	// 创建OpenAI客户端
	client := openai.NewClient(cfg)
	slog.Debug("OpenAI客户端创建成功")

	// 准备系统提示
	systemPrompt := "你是一个命令行命令翻译机，负责将用户输入翻译为命令行命令，你需要以json方式回复，以下是示例\n" +
		"{\"command\": [\"ls\"],\"msg\": \"执行此命令将列出当前目录中的文件和子目录。\",\"code\": 0}\n" +
		"command是可执行命令，可以有多种翻译结果，每一项都是完整的命令，不要把一条命令拆分为开，用户选择其中一条执行，最多为10个，" +
		"msg是展示给用户的提示信息，code为翻译结果，0为成功翻译，1为不能翻译、缺少信息或其他异常情况。"
	slog.Debug("系统提示准备完成", "systemPrompt", systemPrompt)

	// 构建用户提示（包含系统信息）
	userPrompt := sysInfo + "\n" + args[0]
	slog.Debug("用户提示构建完成", "userPrompt", userPrompt)

	var resp *openai.Response
	var reqResp *openai.RequestResponse

	// 根据是否需要显示数据选择不同的方法
	if showData {
		slog.Debug("使用 SendRequestWithData 发送请求")
		// 发送请求到OpenAI并获取请求和响应数据
		reqResp, err = client.SendRequestWithData(systemPrompt, userPrompt)
		if err != nil {
			slog.Error("SendRequestWithData 发送请求失败", "error", err)
			return fmt.Errorf("发送请求失败: %v", err)
		}
		resp = reqResp.Response
		slog.Debug("SendRequestWithData 响应接收成功", "response", resp)
	} else {
		slog.Debug("使用 SendRequest 发送请求")
		// 只发送请求到OpenAI
		resp, err = client.SendRequest(systemPrompt, userPrompt)
		if err != nil {
			slog.Error("SendRequest 发送请求失败", "error", err)
			return fmt.Errorf("发送请求失败: %v", err)
		}
		slog.Debug("SendRequest 响应接收成功", "response", resp)
	}

	if resp == nil || len(resp.Choices) == 0 {
		slog.Error("未收到有效响应或响应中没有Choices")
		return fmt.Errorf("未收到有效响应")
	}
	slog.Debug("OpenAI响应有效", "choicesCount", len(resp.Choices))

	// 解析响应
	var aiResp AIResponse
	content := resp.Choices[0].Message.Content
	slog.Debug("获取到响应内容", "content", content)

	// 如果响应是Markdown json 格式，去除Markdown标记
	if len(content) > 7 && content[:7] == "```json" {
		slog.Debug("去除Markdown json前缀", "originalContent", content)
		content = content[7:]
		slog.Debug("去除Markdown json前缀后", "newContent", content)
	}
	// 如果响应是Markdown js 格式，去除Markdown标记
	if len(content) > 5 && content[:5] == "```js" {
		slog.Debug("去除Markdown js前缀", "originalContent", content)
		content = content[5:]
		slog.Debug("去除Markdown js前缀后", "newContent", content)
	}
	if len(content) > 3 && content[len(content)-3:] == "```" {
		slog.Debug("去除Markdown后缀", "originalContent", content)
		content = content[:len(content)-3]
		slog.Debug("去除Markdown后缀后", "newContent", content)
	}

	slog.Debug("准备解析JSON内容", "contentToParse", content)
	if err := json.Unmarshal([]byte(content), &aiResp); err != nil {
		slog.Error("解析响应JSON失败", "error", err, "content", content)
		return fmt.Errorf("解析响应失败: %v", err)
	}
	slog.Debug("响应JSON解析成功", "aiResponse", aiResp)

	// 如果showData为true，显示发送和接收的数据
	if showData {
		slog.Debug("开始显示发送和接收的数据")
		fmt.Println("=== 请求和响应数据 ===")

		// 显示发送的数据
		fmt.Println("发送数据:")
		requestData, err := json.MarshalIndent(reqResp.Request, "", "  ")
		if err != nil {
			slog.Error("格式化请求数据失败", "error", err)
			return fmt.Errorf("格式化请求数据失败: %v", err)
		}
		fmt.Println(string(requestData))
		slog.Debug("请求数据已显示")

		fmt.Println("\n响应数据:")
		responseData, err := json.MarshalIndent(reqResp.Response, "", "  ")
		if err != nil {
			slog.Error("格式化响应数据失败", "error", err)
			return fmt.Errorf("格式化响应数据失败: %v", err)
		}
		fmt.Println(string(responseData))
		slog.Debug("响应数据已显示")

		fmt.Println("\n解析后的AI响应:")
		aiRespData, err := json.MarshalIndent(aiResp, "", "  ")
		if err != nil {
			slog.Error("格式化AI响应数据失败", "error", err)
			return fmt.Errorf("格式化AI响应数据失败: %v", err)
		}
		fmt.Println(string(aiRespData))
		slog.Debug("解析后的AI响应数据已显示")
		fmt.Println("======================")
	}

	// 输出提示信息
	fmt.Println(aiResp.Msg)
	fmt.Println("---------------------")
	fmt.Println("可用的命令选项:")

	slog.Debug("输出提示信息", "message", aiResp.Msg)
	// 检查翻译结果
	if aiResp.Code != 0 {
		slog.Error("命令翻译失败", "aiResponseCode", aiResp.Code, "aiResponseMessage", aiResp.Msg)
		return fmt.Errorf("命令翻译失败: %s (code: %d)", aiResp.Msg, aiResp.Code)
	}
	slog.Debug("命令翻译成功")

	// 显示可用的命令选项
	slog.Debug("显示可用命令选项", "commands", aiResp.Command)
	for i, cmd := range aiResp.Command {
		fmt.Printf("%d: %s\n", i+1, cmd)
	}
	fmt.Println("0: 退出")

	// 获取用户选择
	fmt.Print("请选择要执行的命令: ")
	var choice string
	fmt.Scanln(&choice)
	slog.Debug("用户选择", "choice", choice)

	// 解析用户选择
	num, err := strconv.Atoi(choice)
	if err != nil {
		slog.Error("无效的用户选择，无法转换为数字", "choice", choice, "error", err)
		return fmt.Errorf("无效的选择")
	}
	slog.Debug("用户选择解析为数字", "number", num)

	if num == 0 {
		slog.Debug("用户选择退出程序")
		fmt.Println("退出程序。")
		return nil
	}

	if num < 1 || num > len(aiResp.Command) {
		slog.Error("用户选择的数字超出范围", "number", num, "commandCount", len(aiResp.Command))
		return fmt.Errorf("无效的选择")
	}

	// 执行选中的命令
	selectedCmd := aiResp.Command[num-1]
	slog.Debug("选中的命令", "selectedCmd", selectedCmd)
	fmt.Printf("执行命令: %s\n", selectedCmd)
	fmt.Println("---------------------")

	// 设置环境变量
	env := os.Environ()
	env = append(env, "TERM=xterm-256color")

	slog.Debug("设置环境变量", "envCount", len(env))
	// 创建命令
	command := exec.Command("bash", "-i", "-c", selectedCmd)
	command.Env = env
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin
	slog.Debug("命令已创建", "commandPath", command.Path, "commandArgs", command.Args)

	// 执行命令
	slog.Debug("开始执行命令")
	if err := command.Run(); err != nil {
		slog.Error("命令执行失败", "error", err, "command", selectedCmd)
		return fmt.Errorf("命令执行失败: %v", err)
	}
	slog.Debug("命令执行成功")

	return nil
}
