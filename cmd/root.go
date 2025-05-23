package cmd

import (
	"github.com/spf13/cobra"
)

var (
	showData bool
	rootCmd  = &cobra.Command{
		Use:   "ais [description]",
		Short: "AI-Shell - 基于 OpenAI 的命令行工具",
		Long: `AI-Shell 是一个基于 OpenAI 的自然语言处理能力开发的命令行工具，
旨在简化 Linux 命令的使用。通过 AI-Shell，用户可以通过自然语言直接
在终端中执行复杂的 Linux 命令，而无需记住繁琐的命令语法。`,
		Args: cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// 如果没有参数，显示帮助信息
			if len(args) == 0 {
				return cmd.Help()
			}
			// 否则执行exec命令
			return runExecute(cmd, args)
		},
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&showData, "show-data", "s", false, "显示发送到API的数据")
}
