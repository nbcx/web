package main

import (
	"log"

	"fmt"
	"net/http"
	"strings"

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
	Addr string `name:"addr" short:"a" value:":8080" usage:"service monitoring address"`
}

func (c *Server) GetUse() string {
	return "web {addr} {path}"
}

func (c *Server) GetLong() string {
	return "static web service"
}

func (c *Server) Exec(r ...string) error {
	path := "./"
	if len(r) > 0 {
		path = r[0]
	}
	// set static dir
	http.Handle("/", http.FileServer(http.Dir(path)))

	// show in terminal and click jump
	show := c.Addr
	if strings.Split(c.Addr, ":")[0] == "" {
		show = fmt.Sprintf("http://localhost%s", c.Addr)
	}
	fmt.Printf("Listening: \x1b]8;;%s\x1b\\%s\x1b]8;;\x1b\\\n", show, show)
	err := http.ListenAndServe(c.Addr, nil)
	return err
}
