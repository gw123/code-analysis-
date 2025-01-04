package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

// classifyCmd 定义了分类文件的命令
var classifyCmd = &cobra.Command{
	Use:   "classify [directory]",
	Short: "Classify files in the specified directory by file type",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dir := args[0]
		fileMap := make(map[string][]string) // 用于存储文件类型及其文件路径

		err := WalkDir(dir, func(path string) {
			ext := strings.ToLower(filepath.Ext(path))
			fileMap[ext] = append(fileMap[ext], path) // 按文件类型分类
		})

		if err != nil {
			fmt.Println("Error walking the path:", err)
			return
		}

		// 输出分类结果
		fmt.Println("Classified files:")
		for ext, files := range fileMap {
			fmt.Printf("\nFiles of type '%s':\n", ext)
			for _, file := range files {
				fmt.Printf(" - %s\n", file)
			}
		}
	},
}

// WalkDir 遍历目录并执行回调函数
func WalkDir(dir string, callback func(path string)) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && !isIgnored(filepath.Base(path)) {
			callback(path)
		}
		return nil
	})
}

// isIgnored 检查目录或文件是否被忽略（这里可以修改为实际的忽略逻辑）
func isIgnored(name string) bool {
	ignored := []string{"vendor", ".git", ".svn"} // 示例忽略列表
	for _, ig := range ignored {
		if strings.Contains(name, ig) {
			return true
		}
	}
	return false
}

func init() {
	rootCmd.AddCommand(classifyCmd) // 将子命令添加到根命令
}
