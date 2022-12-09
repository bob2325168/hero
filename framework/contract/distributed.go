package contract

import "time"

const DistributedKey = "hero:distributed"

type Distributed interface {
	// Select 分布式选择器，所有节点对某个服务进行抢占，只能选择其中一个节点
	// ServiceName 服务器名称， appId 当前的appid， holdTime 分布式选择器持有的时间
	Select(serviceName string, appId string, holdTime time.Duration) (selectAppId string, err error)
}
