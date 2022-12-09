package command

import "github.com/bob2325168/gohero/framework/cobra"

func AddKernelCommands(root *cobra.Command) {

	//启动APP服务
	root.AddCommand(initAppCommand())
	// env 命令
	root.AddCommand(initEnvCommand())
	// cron 命令
	root.AddCommand(initCronCommand())
	// config 命令
	root.AddCommand(initConfigCommand())
	// build 命令
	root.AddCommand(initBuildCommand())
	// go build
	root.AddCommand(goCommand)
	// npm build
	//root.AddCommand(npmCommand)
	// dev
	root.AddCommand(initDevCommand())
	// cmd
	root.AddCommand(initCmdCommand())
	// provider
	root.AddCommand(initProviderCommand())
	// middleware
	root.AddCommand(initMiddlewareCommand())
	// new
	root.AddCommand(initNewCommand())
	// swagger
	//root.AddCommand(initSwaggerCommand())
}
