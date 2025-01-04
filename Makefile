SourceDirectory := "."
TargetDirectory := "./result"
SummaryFile := "./result/all.md"
Token := "sk-bZf8kBSewGAYqJ85CgDZzJtGyBO1AcBdA6OKdy0ntNkUtob6"
export OPENAI_BASE_URL='https://api.chatanywhere.tech/v1'

Question:= "添加一个子命令，该命令的功能是输出一个目录下的所有文件,安装目录结构生成一个图片， 引入图形化库（如 Graphviz）来生成可视化的结构图。"

analyze:
	 go run entry/main.go analyze  -d $(SourceDirectory)  -t $(Token) -o $(TargetDirectory)
.PHONY: analyze

question:
	 go run entry/main.go question -s $(SummaryFile)  -t $(Token)  $(Question)
