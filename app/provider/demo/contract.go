package demo

// Demo 服务的key
const DemoKey = "hero:demo"

// IService Demo服务的接口
type IService interface {
	GetAllStudent() []Student
}

type Student struct {
	ID   int
	Name string
}
