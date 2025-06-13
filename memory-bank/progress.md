# 项目进展记录

## 2025-06-13

*   **已完成**: 适配了应用程序配置文件在 Windows 和 macOS 上的存储位置。
    *   **详情**: 修改了 [`internal/config/config.go`](internal/config/config.go:1) 中的 `init()` 函数，使用 `os.UserConfigDir()` 来确定配置文件的基础路径，并在其下创建 `ais` 子目录存放 `ais_config.json`。
    *   **状态**: 已由 Code 模式子任务完成，并通过用户确认。
    *   **决策记录**: 参见 [`memory-bank/decisionLog.md`](memory-bank/decisionLog.md)

*   **进行中**: 添加 slog 日志支持。
    *   **阶段1 (已完成)**: 更新了 [`internal/config/config.go`](internal/config/config.go:1) 以集成 `slog`。
        *   添加了 `Debug` 配置项和 `Logger` 实例到 `Config` 结构体。
        *   修改了配置加载和设置逻辑以支持和使用日志记录器。
        *   **状态**: 已由 Code 模式子任务完成。详细过程记录在 `memory-bank/activeContext.md` (后已清理)。
    *   **阶段2 (已完成)**: 更新了 [`cmd/root.go`](cmd/root.go:1) 以添加 `--debug` 命令行标志并初始化全局日志。
        *   添加了 `--debug` / `-d` 标志。
        *   在 `PersistentPreRunE` 中实现了日志初始化逻辑，包括根据命令行或配置文件设置debug级别，并将logger设为全局默认。
        *   **状态**: 已由 Code 模式子任务完成。详细过程记录在 `memory-bank/activeContext.md` (后已清理)。
*   **整体状态**: `slog` 日志支持已基本完成。