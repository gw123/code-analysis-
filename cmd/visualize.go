package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var visualizeCmd = &cobra.Command{
	Use:   "visualize [directory]",
	Short: "Visualize the directory structure as a graph",
	Args:  cobra.ExactArgs(1), // 需要一个参数：目录路径
	Run: func(cmd *cobra.Command, args []string) {
		dir := args[0]
		err := visualizeDirectory(dir)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	},
}

// visualizeDirectory 遍历目录并生成可视化图像
func visualizeDirectory(dir string) error {
	var files []string

	// 遍历目录并获取所有的 .go 文件
	err := WalkDir(dir, func(path string) {
		if filepath.Ext(path) == ".go" && !isIgnored(path) {
			files = append(files, path)
		}
	})
	if err != nil {
		return err
	}

	// 生成图形内容
	dot := generateDotStructure(files)

	// 渲染图形
	return renderGraph(dot)
}

// 根据文件生成 DOT 格式字符串
func generateDotStructure(files []string) string {
	var sb strings.Builder
	sb.WriteString("digraph G {\n")

	for _, file := range files {
		relPath, _ := filepath.Rel(filepath.Dir(file), file)
		sb.WriteString(fmt.Sprintf("    \"%s\";\n", relPath))
	}

	sb.WriteString("}\n")
	return sb.String()
}

// 使用 Graphviz 渲染图形成 PNG
func renderGraph(dot string) error {
	// 保存 DOT 文件
	err := os.WriteFile("structure.dot", []byte(dot), 0644)
	if err != nil {
		return fmt.Errorf("error writing DOT file: %v", err)
	}

	// 调用 Graphviz 生成图像
	cmd := exec.Command("dot", "-Tpng", "structure.dot", "-o", "structure.png")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error generating image: %v", err)
	}

	fmt.Println("Structure visualization saved as structure.png")
	return nil
}

func init() {
	rootCmd.AddCommand(visualizeCmd)
}
