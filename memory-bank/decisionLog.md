# 项目决策日志

## 2025-06-13: 配置文件路径适配

*   **决策**: 修改应用程序的配置文件存储位置，以更好地适配 Windows 和 macOS 操作系统。
*   **实现**:
    *   使用 Go 语言的 `os.UserConfigDir()` 函数获取操作系统推荐的用户配置目录。
    *   在该目录下创建一个名为 "ais" 的子目录。
    *   配置文件 `ais_config.json` 将存储在 `UserConfigDir/ais/` 路径下。
*   **理由**: 提高应用程序在不同操作系统上的用户体验和标准化。
*   **执行者**: NexusCore 协调，子任务由 Code 模式完成。
*   **相关文件**: [`internal/config/config.go`](internal/config/config.go:1)
*   **详细过程**: 记录在 `memory-bank/activeContext.md` (临时日志，内容已归档处理)。

## 2025-06-13: 添加 slog 日志支持 (阶段1: 配置更新)

*   **决策**: 为项目配置模块添加 `slog` 日志支持，允许通过配置启用/禁用debug日志。
*   **实现**:
    *   在 `internal/config/config.go` 的 `Config` 结构体中添加 `Debug bool` 和 `Logger *slog.Logger` 字段。
    *   修改 `LoadConfig` 函数，根据 `Debug` 字段的值初始化 `Logger`，并在 `Debug` 为 `true` 时设置日志级别为 `slog.LevelDebug`。
    *   在各 `Set<FieldName>` 方法中添加了对 `Logger.Debug` 的调用，以记录配置变更。
    *   添加 `SetDebug` 方法以允许动态更改日志级别并保存配置。
*   **理由**: 增强应用的可调试性，提供结构化的日志输出。
*   **执行者**: NexusCore 协调，子任务由 Code 模式完成。
*   **相关文件**: [`internal/config/config.go`](internal/config/config.go:1)
*   **详细过程**: 记录在 `memory-bank/activeContext.md` (临时日志，内容已归档处理)。

## 2025-06-13: 添加 slog 日志支持 (阶段2: 命令行集成)

*   **决策**: 在 `cmd/root.go` 中集成 `--debug` 命令行标志，并初始化全局 `slog` 日志记录器。
*   **实现**:
    *   在 `cmd/root.go` 中添加了 `debugMode`布尔变量。
    *   通过 `rootCmd.PersistentFlags().BoolVarP` 添加了 `--debug` (简写 `-d`) 标志，绑定到 `debugMode`。
    *   在 `rootCmd` 的 `PersistentPreRunE` 函数中：
        *   加载配置 (`config.LoadConfig()`)。
        *   如果命令行 `--debug` 标志被设置，则调用 `cfg.SetDebug(true)` 来更新配置和logger实例。
        *   根据命令行标志或配置文件中的 `Debug` 设置，记录相应的debug激活信息。
        *   调用 `slog.SetDefault(cfg.Logger)` 将配置好的logger设置为全局默认logger。
        *   记录 "日志系统初始化完成" Info级别日志。
*   **理由**: 允许用户通过命令行方便地控制debug日志的输出，并确保日志系统在程序早期正确初始化。
*   **执行者**: NexusCore 协调，子任务由 Code 模式完成。
*   **相关文件**: [`cmd/root.go`](cmd/root.go:1), [`internal/config/config.go`](internal/config/config.go:1)
*   **详细过程**: 记录在 `memory-bank/activeContext.md` (临时日志，内容已归档处理)。