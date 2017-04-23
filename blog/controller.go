package blog

import (
	"fmt"
	"html/template"
)

type BlogController struct {
	ControllerName string
	ActionName     string
	Template       string
	Http           *BlogHttp
	Assign         map[interface{}]interface{}
	IsGet          bool
	IsPost         bool
	Session        *Session
}

//ControllerInterface 控制器接口
type ControllerInterface interface {
	Init(controllerName, actionName string, blogHttp *BlogHttp)
	BeforeAction()
	AfterAction()
	Render()
}

func (c *BlogController) Init(controllerName, actionName string, blogHttp *BlogHttp) {
	c.ControllerName = controllerName
	c.ActionName = actionName
	c.Http = blogHttp
	c.Assign = make(map[interface{}]interface{})
	c.Session = globalSession.Start(blogHttp.Writer, blogHttp.Reader)
}

func (c *BlogController) BeforeAction() {}

func (c *BlogController) AfterAction() {}

//自动渲染模板
func (c *BlogController) Render() {
	if c.Template != "" {
		c.Http.Writer.Header().Set("Content-Type", "text/html; charset="+BlogConfig.Charset)
		t, err := template.ParseFiles(c.Template)

		if err != nil {
			panic(fmt.Sprintf("Template %s Error.", c.Template))
		}
		t.Execute(c.Http.Writer, c.Assign)
	}
}
