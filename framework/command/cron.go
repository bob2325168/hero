package command

import (
	"fmt"
	"github.com/bob2325168/gohero/framework/cobra"
	"github.com/bob2325168/gohero/framework/contract"
	"github.com/bob2325168/gohero/framework/util"
	"github.com/erikdubbelboer/gspt"
	"github.com/sevlyar/go-daemon"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
)

var cronDaemon = false

func initCronCommand() *cobra.Command {

	cronStartCommand.Flags().BoolVarP(&cronDaemon, "daemon", "d", false, "start sever daemon")
	cronCommand.AddCommand(cronStartCommand)
	cronCommand.AddCommand(cronRestartCommand)
	cronCommand.AddCommand(cronStateCommand)
	cronCommand.AddCommand(cronListCommand)
	cronCommand.AddCommand(cronStopCommand)

	return cronCommand
}

var cronCommand = &cobra.Command{
	Use:   "cron",
	Short: "定时任务相关命令",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			cmd.Help()
		}
		return nil
	},
}

var cronStartCommand = &cobra.Command{

	Use:   "start",
	Short: "启动cron后台进程",
	RunE: func(cmd *cobra.Command, args []string) error {
		//获取容器
		container := cmd.GetContainer()
		//获取容器中的app服务
		appService := container.MustMake(contract.AppKey).(contract.App)
		//设置cron的日志地址和进程ID地址
		pidFolder := appService.RuntimeFolder()
		serverPidFile := filepath.Join(pidFolder, "cron.pid")
		logFolder := appService.LogFolder()
		logFile := filepath.Join(logFolder, "cron.log")
		currentFolder := appService.BaseFolder()

		//如果启动daemon模式
		if cronDaemon {
			//创建一个context
			cntxt := &daemon.Context{
				// 设置pid文件
				PidFileName: serverPidFile,
				PidFilePerm: 0664,
				LogFileName: logFile,
				LogFilePerm: 0640,
				//设置工作路径
				WorkDir: currentFolder,
				//设置所有设置文件的mask，默认为570
				Umask: 027,
				//子进程的参数，按照这个参数设置，子进程命令是： ./gohero cron start --daemon=true
				Args: []string{"", "cron", "start", "--daemon=true"},
			}
			//启动子进程，d不为空表示当前是父进程，为空表示是当前子进程
			d, err := cntxt.Reborn()
			if err != nil {
				return err
			}
			if d != nil {
				// 父进程直接打印成功信息
				fmt.Println("cron server started, pid: ", d.Pid)
				fmt.Println("log file: ", logFile)
				return nil
			}

			// 子进程执行cron.run
			defer cntxt.Release()
			fmt.Println("daemon started")
			gspt.SetProcTitle("gohero cron")
			//cmd.Root().Cron.Run()
			return nil
		}

		// 不启动daemon模式
		fmt.Println("start cron job")
		content := strconv.Itoa(os.Getpid())
		fmt.Println("[PID]", content)
		err := os.WriteFile(serverPidFile, []byte(content), 0664)
		if err != nil {
			return err
		}
		gspt.SetProcTitle("gohero cron")
		//cmd.Root().Cron.Run()
		return nil
	},
}

var cronRestartCommand = &cobra.Command{

	Use:   "restart",
	Short: "重启cron后台进程",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := cmd.GetContainer()
		appService := container.MustMake(contract.AppKey).(contract.App)
		// 获取PID
		pidFile := filepath.Join(appService.RuntimeFolder(), "cron.pid")
		content, err := os.ReadFile(pidFile)
		if err != nil {
			return err
		}
		if content != nil && len(content) > 0 {
			pid, err := strconv.Atoi(string(content))
			if err != nil {
				return err
			}
			if util.CheckProcessExist(pid) {
				if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
					return err
				}
				// 检查进程是否closed
				for i := 0; i < 10; i++ {
					if util.CheckProcessExist(pid) == false {
						break
					}
					time.Sleep(1 * time.Second)
				}
				fmt.Println("kill process: " + strconv.Itoa(pid))
			}
		}
		cronDaemon = true
		return cronStartCommand.RunE(cmd, args)
	},
}

var cronStateCommand = &cobra.Command{

	Use:   "state",
	Short: "cron常驻进程状态",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()
		appService := container.MustMake(contract.AppKey).(contract.App)

		// GetPid
		serverPidFile := filepath.Join(appService.RuntimeFolder(), "cron.pid")

		content, err := os.ReadFile(serverPidFile)
		if err != nil {
			return err
		}

		if content != nil && len(content) > 0 {
			pid, err := strconv.Atoi(string(content))
			if err != nil {
				return err
			}
			if util.CheckProcessExist(pid) {
				fmt.Println("cron server started, pid:", pid)
				return nil
			}
		}
		fmt.Println("no cron server start")
		return nil
	},
}

var cronListCommand = &cobra.Command{
	Use:   "list",
	Short: "列出所有的定时任务",
	RunE: func(c *cobra.Command, args []string) error {

		//cronSpecs := c.Root().CronSpecs
		//var ps [][]string
		//for _, cronSpec := range cronSpecs {
		//	line := []string{cronSpec.Type, cronSpec.Spec, cronSpec.Cmd.Use, cronSpec.Cmd.Short, cronSpec.ServiceName}
		//	ps = append(ps, line)
		//}
		//util.PrettyPrint(ps)
		return nil
	},
}

var cronStopCommand = &cobra.Command{
	Use:   "stop",
	Short: "停止cron常驻进程",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()
		appService := container.MustMake(contract.AppKey).(contract.App)

		// GetPid
		serverPidFile := filepath.Join(appService.RuntimeFolder(), "cron.pid")

		content, err := ioutil.ReadFile(serverPidFile)
		if err != nil {
			return err
		}

		if content != nil && len(content) > 0 {
			pid, err := strconv.Atoi(string(content))
			if err != nil {
				return err
			}
			if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
				return err
			}
			if err := os.WriteFile(serverPidFile, []byte{}, 0644); err != nil {
				return err
			}
			fmt.Println("stop pid:", pid)
		}
		return nil
	},
}
