package cobra

import (
	"github.com/bob2325168/gohero/framework"
)

// SetContainer 设置服务容器
func (c *Command) SetContainer(container framework.Container) {
	c.container = container
}

// GetContainer 获取服务容器
func (c *Command) GetContainer() framework.Container {
	return c.Root().container
}

//type CronSpec struct {
//	Type        string
//	Cmd         *Command
//	Spec        string
//	ServiceName string
//}
//
//func (c *Command) SetParentNull() {
//	c.parent = nil
//}
//
//// AddCronCommand 创建一个Cron任务
//func (c *Command) AddCronCommand(spec string, cmd *Command) {
//	// cron结构是挂载在根command上的
//	root := c.Root()
//	if root.Cron == nil {
//		root.Cron = cron.New(cron.WithParser(cron.NewParser(cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)))
//		root.CronSpecs = []CronSpec{}
//	}
//	// 增加说明信息
//	root.CronSpecs = append(root.CronSpecs, CronSpec{
//		Type: "normal-cron",
//		Cmd:  cmd,
//		Spec: spec,
//	})
//	// 创建一个root Command
//	var cronCmd Command
//	ctx := root.Context()
//	cronCmd = *cmd
//	cronCmd.args = []string{}
//	cronCmd.SetParentNull()
//	cronCmd.SetContainer(root.GetContainer())
//	//增加调用函数
//	root.Cron.AddFunc(spec, func() {
//		defer func() {
//			// 捕捉可能出现的panic
//			if err := recover(); err != nil {
//				log.Println(err)
//			}
//		}()
//		err := cronCmd.ExecuteContext(ctx)
//		if err != nil {
//			log.Println(err)
//		}
//	})
//}
