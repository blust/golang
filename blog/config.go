package blog

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port    int    //端口号
	Static  string //静态文件目录
	Charset string //编码
}

var BlogConfig *Config

func init() {
	BlogConfig = &Config{
		Port:    8080,
		Static:  "/static/",
		Charset: "utf-8",
	}

	conf := "myblog.conf"
	if FileExists(conf) {
		f, err := os.Open(conf)
		if err != nil {
			panic(err)
		}

		defer f.Close()

		r := bufio.NewReader(f)

		for {
			b, _, err := r.ReadLine()
			if err != nil {
				if err == io.EOF {
					break
				}
				panic(err)
			}

			s := strings.TrimSpace(string(b))
			//就当是注释吧...
			if strings.Index(s, "#") == 0 {
				continue
			}

			index := strings.Index(s, "=")
			if index < 0 {
				continue
			}

			key := strings.ToLower(strings.TrimSpace(s[:index]))
			value := strings.TrimSpace(s[index+1:])

			switch key {
			case "port":
				BlogConfig.Port, _ = strconv.Atoi(value)
			case "static":
				BlogConfig.Static = value
			}
		}
	}
}
