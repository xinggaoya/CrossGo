package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func Run() {
	// 获取命令行参数
	// os.Args[0] 是命令本身的名字 os.Args[1:] 是传递给命令的参数
	args := os.Args[1:]

	for _, arg := range args {
		// -test=123 解析成 -test 123
		mapValue := strings.Split(arg, "=")
		key := mapValue[0]
		value := mapValue[1]

		// 打印参数
		println(key, value)
	}
	CrossCompile()
}

// CrossCompile 交叉编译
func CrossCompile() {
	buildPath := "bin"

	platforms := []string{"linux", "windows", "darwin"}
	for _, p := range platforms {
		// os.Setenv 设置环境变量
		os.Setenv("GOOS", p)
		os.Setenv("GOARCH", "amd64")

		// 保存路径
		path := buildPath + "/main-" + p
		// windows 下修改路径
		if p == "windows" {
			path += ".exe"
		}

		// 执行命令
		cmd := exec.Command("go", "build", "-o", path, "main.go")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Errorf("Error: %s", err.Error())
		}
		log.Printf("CrossCompile %s Success\n", p)
	}
	// 获取当前工作目录
	dir, _ := os.Getwd()
	log.Printf("CrossCompile Build Success, Path: %s/%s\n", dir, buildPath)
}
