package main

import (
	"log"
	"web/pkg"

	"github.com/nbcx/boot"
)

func main() {
	app := boot.First(&Server{})
	if err := boot.Execute(app); err != nil {
		log.Fatalf("server run err: %v", err)
	}
}

type Server struct {
	boot.Default
	Addr  string `name:"addr" short:"a" value:":8080" usage:"service monitoring address"`
	Https bool   `name:"https" short:"s" value:"false" usage:"enable https server"`
}

func (c *Server) GetUse() string {
	return "web {path}"
}

func (c *Server) GetLong() string {
	return "static web service"
}

func (c *Server) Exec(r ...string) error {
	path := "."
	if len(r) > 0 {
		path = r[0]
	}
	if c.Https {
		pkg.Https(c.Addr, path)
	} else {
		return pkg.Http(c.Addr, path)
	}
	return nil
}
