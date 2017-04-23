package blog

import (
	"fmt"
	"net/http"
	"time"
)

var BlogServer *BServer

func init() {
	BlogServer = &BServer{Version: "1.0.0", ServerName: "nginx/1.6.2", Powered: "PHP/8.0.3"}
}

type BServer struct {
	Version    string
	ServerName string
	Powered    string
	IsStop     chan bool
}

func (b *BServer) Run() {

	BlogServer.IsStop = make(chan bool, 1)

	addr := fmt.Sprintf(":%d", BlogConfig.Port)
	s := &http.Server{
		Addr:           addr,
		Handler:        &BlogHandle{},
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil {
			time.Sleep(100 * time.Microsecond)
			BlogServer.IsStop <- true
		}
	}()

	<-BlogServer.IsStop
}
