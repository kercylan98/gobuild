package gobuild

import (
	"fmt"
	"testing"
)

func TestStart(t *testing.T) {
	status, err := Start().
		SetOutputName("result").
		SetProjectPath("/Users/kercylan/Coding.localized/Go/gobuild/main").
		SetCompileFiles().
		CompileDefault()

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("cmd:", status.Cmd)
	for _, s := range status.Stdout {
		t.Log(s)
	}
}
