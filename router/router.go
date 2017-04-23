package router

import (
	"github.com/zhjx922/myblog/blog"
	"github.com/zhjx922/myblog/controllers"
)

func init()  {
	blog.AddRouter("/", &controllers.HomeController{})
	blog.AddRouter("/home", &controllers.HomeController{})
}