package command

import (
	"fmt"
	"github.com/bob2325168/gohero/framework/cobra"
	"github.com/bob2325168/gohero/framework/contract"
)

// DemoCommand show current envionment
var DemoCommand = &cobra.Command{
	Use:   "demo",
	Short: "demo for framework",
	Run: func(c *cobra.Command, args []string) {
		container := c.GetContainer()
		appService := container.MustMake(contract.AppKey).(contract.App)
		fmt.Println("app base folder:", appService.BaseFolder())
	},
}
