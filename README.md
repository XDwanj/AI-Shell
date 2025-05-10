# AI-Shell

AI-Shell 是一个基于 OpenAI 的自然语言处理能力开发的命令行工具，旨在简化 Linux 命令的使用。通过 AI-Shell，用户可以通过自然语言直接在终端中执行复杂的 Linux 命令，而无需记住繁琐的命令语法。

## 特性

- **自然语言命令**：输入自然语言描述，即可生成并执行相应的 Linux 命令
- **智能命令补全**：基于上下文理解，自动生成最合适的命令
- **可扩展性**：支持自定义命令和语义映射，适应多种应用场景
- **开箱即用**：简单易用的命令行工具，适合从初学者到高级用户

## 安装

### 从源码安装

```bash
# 克隆仓库
git clone https://github.com/XDwanj/AI-Shell.git

# 进入项目目录
cd AI-Shell

# 安装
go install
```

## 使用说明

### 基本用法

```bash
# 直接使用自然语言描述
ais "检查磁盘使用情况"

# 使用 -s 标志显示API调用数据
ais -s "在当前目录下查找所有的 .txt 文件"
```

### 配置命令

AI-Shell 提供了一系列配置命令来管理设置：

```bash
# 查看当前配置
ais config view

# 设置 API URL
ais config set url https://api.openai.com/v1/chat/completions

# 设置 API 密钥
ais config set key sk-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

# 设置模型
ais config set model gpt-4

# 设置最大令牌数
ais config set max-tokens 1000

# 设置温度参数（控制随机性，范围 0-1）
ais config set temperature 0.7
```

### 配置文件

配置文件位于 `~/.config/ais_config.json`，包含以下配置项：

```json
{
  "url": "https://api.openai.com/v1/chat/completions",
  "api_key": "your-api-key",
  "model": "gpt-4",
  "max_tokens": 1000,
  "temperature": 0.7
}
```

## 使用示例

1. 查找文件：

   ```bash
   ais "找出所有大于100MB的文件"
   ```

2. 系统管理：

   ```bash
   ais "显示当前系统资源使用情况"
   ```

3. 网络诊断：

   ```bash
   ais "检查与 google.com 的网络连接"
   ```

4. 文本处理：

   ```bash
   ais "统计当前目录下所有 .go 文件的代码行数"
   ```

## 项目灵感

感谢 [by123456by/AI-Shell](https://github.com/by123456by/AI-Shell) 带来的灵感

## 贡献

欢迎提交问题、建议或 PR 来改善 AI-Shell。我们希望与开源社区共同发展，让 AI-Shell 更加强大和易用。

## 许可证

[MIT License](LICENSE)
