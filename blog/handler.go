package blog

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"
)

type BlogHttp struct {
	Writer http.ResponseWriter
	Reader *http.Request
}

type BlogHandle struct{}

func (*BlogHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer RunTime("ServeHTTP", time.Now())

	w.Header().Set("Server", BlogServer.ServerName)
	w.Header().Set("Version", BlogServer.Version)
	w.Header().Set("X-Powered-By", BlogServer.Powered)

	p := strings.Split(r.URL.Path, "/")

	//@todo 静态Server

	if p[1] == "favicon.ico" {
		return
	}

	controllerName := p[1]
	actionName := "Index"
	if len(p) >= 3 {
		if strings.HasSuffix(p[2], ".php") {
			actionName = p[2][0 : len(p[2])-4]
		}
	}

	isFound := false

	for route, i := range BlogRouter.Routers {
		if strings.Trim(route, "/") == controllerName {
			v := reflect.ValueOf(i)

			controller, ok := v.Interface().(ControllerInterface)
			if !ok {
				panic("controller is not ControllerInterface")
			}

			blogHttp := &BlogHttp{Writer: w, Reader: r}
			controller.Init(strings.Title(controllerName), strings.Title(actionName), blogHttp)

			controller.BeforeAction()

			var in []reflect.Value
			method := v.MethodByName(strings.Title(actionName))

			if method.IsValid() {
				isFound = true
				method.Call(in)
				controller.Render()
			}

			controller.AfterAction()

			break
		}
	}

	if !isFound {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(404)
		fmt.Fprintln(w, "Oh 404.")
	}
}
