package cmd

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/exec"
	"path"
)

func Run() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	})))
	// 获取命令行参数
	// os.Args[0] 是命令本身的名字 os.Args[1:] 是传递给命令的参数
	args := os.Args[1:]

	SelectCmd(args)
}

// SelectCmd 选择命令
func SelectCmd(args []string) {
	if len(args) == 0 {
		println("Usage: cs [command]")
		return
	}
	switch args[0] {
	case "build":
		if len(args) > 1 && args[1] != "" {
			CrossCompile(args[1])
		} else {
			CrossCompile("")
		}
	case "v":
		fmt.Println("cs version: 0.0.1")
	default:
		fmt.Println("Command Not Found")
	}
}

// CrossCompile 交叉编译
func CrossCompile(fileName string) {
	buildPath := "bin"

	system := []string{"linux", "windows", "darwin"}
	arch := []string{"amd64", "arm64"}
	for _, p := range system {
		// os.Setenv 设置环境变量
		for _, a := range arch {
			buildPath = "bin/" + p + "-" + a
			os.Setenv("GOOS", p)
			os.Setenv("GOARCH", a)

			// 保存路径
			path := ""
			buildPath = "bin/" + p + "-" + a
			if fileName != "" {
				path = buildPath + "/" + fileName
			} else {
				path = buildPath + "/main-" + p
			}
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
				fmt.Errorf("error: %s", err.Error())
			}
			log.Printf("CrossCompile %s Success\n", p)
		}
	}
	// 获取当前工作目录
	dir, _ := os.Getwd()
	log.Printf("CrossCompile Build Success, Path: %s\n", path.Join(dir, "bin"))
}
