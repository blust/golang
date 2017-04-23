package blog

type Router struct {
	Routers map[string]ControllerInterface
}

var BlogRouter *Router

func init() {
	BlogRouter = &Router{}
	BlogRouter.Routers = make(map[string]ControllerInterface)
}

func AddRouter(path string, c ControllerInterface, method ...string) {
	BlogRouter.Routers[path] = c
}
