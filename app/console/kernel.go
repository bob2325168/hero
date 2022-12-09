package console

import (
	"github.com/bob2325168/gohero/app/console/command/demo"
	"github.com/bob2325168/gohero/framework"
	"github.com/bob2325168/gohero/framework/cobra"
	"github.com/bob2325168/gohero/framework/command"
)

// RunCommand 初始化根command并运行
func RunCommand(container framework.Container) error {
	//根command
	var rootCmd = &cobra.Command{
		// 定义根命令的关键词
		Use: "hero",
		// 简短介绍
		Short: "hero 命令",
		// 根命令详细介绍
		Long: "gohero 框架提供的命令行工具，使用这个命令行工具可以很方便的执行框架自带命令，也可以很方便的编写业务命令",
		// 执行根命令函数
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.InitDefaultHelpFlag()
			return cmd.Help()
		},
		// 不需要出现cobra 默认的completion子命令
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}
	// 给根command设置服务容器
	rootCmd.SetContainer(container)
	// 绑定框架命令
	command.AddKernelCommands(rootCmd)
	AddAppCommand(rootCmd)
	//执行rootCmd命令
	return rootCmd.Execute()
}

// AddAppCommand 绑定业务命令
func AddAppCommand(cmd *cobra.Command) {
	// 正常调用
	cmd.AddCommand(demo.InitFoo())
	// 每秒调用一次Foo命令
	//cmd.AddCronCommand("* * * * * *", demo.FooCommand)
	// 启动分布式任务调度，每个节点每5秒调用一次foo命令，抢占到了调度任务的节点将抢占锁持续挂载2秒才释放
	// cmd.AddDistributedCronCommand("foo_func_for_test", "*/5 * * * * *", demo.FooCommand, 2*time.Second)
}
