package cmd

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

const defaultBuildPath = "bin"

func Run() {
	// 设置默认日志处理器
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// 解析命令行参数
	flag.Parse()
	args := flag.Args()

	SelectCmd(args)
}

func SelectCmd(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: cs [command]")
		return
	}
	switch args[0] {
	case "build":
		if args[1] != "" {
			CrossCompile(args[1])
		}
		CrossCompile("")
	case "v":
		fmt.Println("cs version: 0.0.1")
	default:
		fmt.Println("Command Not Found")
	}
}

func CrossCompile(name string) {

	systems := []string{"linux", "windows", "darwin"}
	architectures := []string{"amd64", "arm64"}

	// 获取当前工作目录
	dir, err := os.Getwd()
	if err != nil {
		log.Printf("Error getting current directory: %v", err)
		return
	}
	// 获取最后一个目录名
	lastDir := filepath.Base(dir)
	if name != "" {
		lastDir = name
	}
	for _, p := range systems {
		for _, a := range architectures {
			go func() {
				buildPath := path.Join(defaultBuildPath, fmt.Sprintf("%s-%s", p, a))
				if err := os.MkdirAll(buildPath, 0755); err != nil {
					log.Fatalf("Error creating directory: %v", err)
				}

				os.Setenv("GOOS", p)
				os.Setenv("GOARCH", a)

				outputFile := path.Join(buildPath, lastDir)
				if p == "windows" {
					outputFile += ".exe"
				}

				cmd := exec.Command("go", "build", "-o", outputFile)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				if err = cmd.Run(); err != nil {
					log.Fatalf("Error compiling for %s-%s: %v", p, a, err)
				}
				log.Printf("CrossCompile %s-%s Success\n", p, a)
			}()
		}
	}

	log.Printf("CrossCompile Build Success, Path: %s\n", path.Join(dir, defaultBuildPath))
}
