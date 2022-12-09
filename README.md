# GoHero
GoHero是一个注重协议，注重开发效率的前后端一体化框架。hero框架的目标在于将go+vue模型的应用开发最简化，并且提供统一，一体化的脚手架工具促进业务开发。

hero框架在于其制定了一系列的基本协议，在具体的业务逻辑中，可以通过每个协议独特的Key来从全局容器中获取已经注入的服务实例。所有的具体应用开发，在业务逻辑中，都是按照hero约定的协议进行逻辑处理，从而脱离了具体的每个服务所定义的个性化差异。

# 框架特色：

## 基于协议

服务与服务间的协议是基于协议进行交互的。

## 前后端协同(TODO)

前后端协同开发

## 命令行

有充分的命令行工具

# 环境变量

## 设置

hade 支持使用应用默认下的隐藏文件 `.env` 来配置各个机器不同的环境变量。

```
APP_ENV=development

DB_PASSWORD=mypassword
```

环境变量的设置可以在配置文件中通过 `env([环境变量])` 来获取到。

比如：

```
mysql:
    hostname: 127.0.0.1
    username: root
    password:  env(DB_PASSWORD)
    timeout: 1
    readtime: 2.3
```


## 应用环境

hero 启动应用的默认应用环境为 development。

你可以通过设置 .env 文件中的 APP_ENV 设置应用环境。

应用环境建议选择：
- development // 开发使用
- production // 线上使用
- testing //测试环境

应用环境对应配置的文件夹，配置服务会去对应应用环境的文件夹中寻找配置。

比如应用环境为 development，在代码中使用
```
configService := container.MustMake(contract.ConfigKey).(contract.Config)
url := configService.GetString("app.url")
```

查找文件为：`config/development/app.yaml`

通过命令`./hade env`也可以获取当前应用环境：

```
[~/Documents/workspace/hade_workspace/demo5]$ ./hade env
environment: development
```

# 命令

## 指南

hero 允许自定义命令，挂载到 hero 上。并且提供了`./hero command` 系列命令。

```
[~/Documents/workspace/hade_workspace/demo5]$ ./hero command
all about commond

Usage:
  hero command [flags]
  hero command [command]

Available Commands:
  list        show all command list
  new         create a command

Flags:
  -h, --help   help for command

Use "hero command [command] --help" for more information about a command.
```

- list  // 列出当前所有已经挂载的命令列表
- new   // 创建一个新的自定义命令

## 创建

创建一个新命令，可以使用 `./hero command new`

这是一个交互式的命令行工具。

```
[~/demo ]$ ./hero command new
create a new command...
? please input command name: test
? please input file name(default: command name):
create new command success，file path: /Users/ZhangLi/demo/app/console/command/test.go
please remember add command to console/kernel.go
```

创建完成之后，会在应用的 app/console/command/ 目录下创建一个新的文件。

## 自定义

hero 中的命令使用的是 cobra 库。 https://github.com/spf13/cobra

## 挂载

编写完自定义命令后，请记得挂载到 console/kernel.go 中。
