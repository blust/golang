package controllers

import (
	"fmt"

	"github.com/zhjx922/myblog/blog"
)

type HomeController struct {
	blog.BlogController
}

func (this *HomeController) Test() {
	//this.Session.Set("a", "b")
	//fmt.Printf("a:%s\n", this.Session.Get("a").(string))
	fmt.Printf("controller:%s,action:%s\n", this.ControllerName, this.ActionName)
	this.Assign["Welcome"] = "Welcome To zhjx922DeBlog!"
	this.Template = "views/test.html"
}
