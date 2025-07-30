package pkg

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func Http(addr string, path string) error {
	// path := "./"
	// if len(r) > 0 {
	// 	path = r[0]
	// }

	// 切换工作目录
	err := os.Chdir(path)
	if err != nil {
		log.Fatalf("无法切换到目录 %s: %v", path, err)
	}
	// set static dir
	http.Handle("/", http.FileServer(http.Dir(".")))

	// show in terminal and click jump
	show := addr
	if strings.Split(addr, ":")[0] == "" {
		show = fmt.Sprintf("http://localhost%s", addr)
	}
	fmt.Printf("Listening: \x1b]8;;%s\x1b\\%s\x1b]8;;\x1b\\\n", show, show)
	err = http.ListenAndServe(addr, nil)
	return err
}
