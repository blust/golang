package main

import (
	"github.com/zhjx922/myblog/blog"
	_ "github.com/zhjx922/myblog/router"
)

//var SessionManager *blog.BlogSessionManager

func main() {
	//SessionManager = blog.NewManager("PHPSESSIONID", 3600)
	//go SessionManager.GC()
	blog.Start()
}
