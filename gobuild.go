package gobuild

import (
	"errors"
	"github.com/go-cmd/cmd"
	"os"
	"path/filepath"
	"strings"
)

// Start 编译特定路径的Go项目
func Start() *OutputName {
	return &OutputName{}
}

type OutputName struct {
	name string // 输出执行文件名称
}

// SetOutputName 设置编译后的可执行文件名称
func (slf *OutputName) SetOutputName(name string) *ProjectPath {
	slf.name = name
	return &ProjectPath{receive: slf}
}

type ProjectPath struct {
	receive *OutputName
	path string // 项目路径
}

// SetProjectPath 设置待编译项目路径
func (slf *ProjectPath) SetProjectPath(path string) *CompileFiles {
	slf.path = path
	return &CompileFiles{receive: slf}
}

type CompileFiles struct {
	receive *ProjectPath
	files 	[]string
}

// SetCompileFiles 设置待编译文件，不设置任何文件时为main.go
func (slf *CompileFiles) SetCompileFiles(filename ...string) *Exec {
	slf.files = filename
	return &Exec{receive: slf}
}

type Exec struct {
	receive *CompileFiles
	goBin string
}

func (slf *Exec) initEnv(gopath ...string) error {
	for _, goRoot := range append(gopath, os.Getenv("GOROOT")) {
		goBin := filepath.Join(goRoot, "bin", "go")
		if fileExist(goBin) {
			slf.goBin = goBin
			return nil
		}
	}
	return errors.New("not found go executable file, please check GOROOT")
}

// CompileDefault 编译当前默认操作系统环境
func (slf *Exec) CompileDefault(gopath ...string) (cmd.Status, error) {
	if err := slf.initEnv(gopath...); err != nil {
		return cmd.Status{}, err
	}

	var args = []string{"build"}
	if outputName := strings.TrimSpace(slf.receive.receive.receive.name); outputName != "" {
		args = append(args, "-o", outputName)
	}
	for _, file := range slf.receive.files {
		args = append(args, filepath.Join(slf.receive.receive.path, file))
	}
	buildCmd := cmd.NewCmd(slf.goBin, args...)
	buildCmd.Dir = slf.receive.receive.path

	status := <- buildCmd.Start()
	if status.Error != nil {
		return cmd.Status{}, status.Error
	}
	return status, nil
}

// fileExist 检查文件是否存在
func fileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}