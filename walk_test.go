package code

import (
	"codetest/internal/usecase/web_api"
	"fmt"
	"testing"
)

func TestWalkDir(t *testing.T) {
	parse := web_api.NewParser()
	WalkDir("/Users/gaowei7/code/go/src/gitlabee.com/licloud-workflow-service", func(path string) {
		fmt.Print(path)
		file, err := parse.ParseByFile(path)
		if err != nil {
			return
		}
		file.PrintResults()

	})
}
