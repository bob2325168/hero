package demo

import (
	demoService "github.com/bob2325168/gohero/app/provider/demo"
	"github.com/bob2325168/gohero/framework/gin"
)

type DemoApi struct {
	service *Service
}

func NewDemoApi() *DemoApi {
	s := NewService()
	return &DemoApi{service: s}
}

func Register(r *gin.Engine) error {
	apiService := NewDemoApi()
	r.Bind(&demoService.DemoProvider{})

	r.GET("/demo/demo", apiService.Demo)
	r.GET("/demo/demo2", apiService.Demo2)
	r.POST("/demo/demo_post", apiService.DemoPost)
	r.GET("/demo/orm", apiService.DemoOrm)
	return nil
}

// Demo godoc
// @Summary 获取所有用户
// @Description 获取所有用户
// @Produce  json
// @Tags demo
// @Success 200 array []UserDTO
// @Router /demo/demo [get]
func (api *DemoApi) Demo(c *gin.Context) {
	//appService := c.MustMake(contract.AppKey).(contract.App)
	//baseFolder := appService.BaseFolder()
	//c.JSON(200, baseFolder)
	users := api.service.GetUsers()
	usersDTO := UserModelsToUserDTOs(users)
	c.JSON(200, usersDTO)
}

// Demo2 godoc
// @Summary 获取所有学生
// @Description 获取所有学生
// @Produce  json
// @Tags demo
// @Success 200 array []UserDTO
// @Router /demo/demo2 [get]
func (api *DemoApi) Demo2(c *gin.Context) {
	demoProvider := c.MustMake(demoService.DemoKey).(demoService.IService)
	students := demoProvider.GetAllStudent()
	usersDTO := StudentsToUserDTOs(students)
	c.JSON(200, usersDTO)
}

func (api *DemoApi) DemoPost(c *gin.Context) {
	type Foo struct {
		Name string
	}
	foo := &Foo{}
	err := c.BindJSON(&foo)
	if err != nil {
		c.AbortWithError(500, err)
	}
	c.JSON(200, nil)
}
