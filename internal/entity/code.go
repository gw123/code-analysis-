package entity

// AICodeSnippet represents a code snippet record in the database
type AICodeSnippet struct {
	ID              int      `gorm:"primaryKey;autoIncrement" json:"id"`                    // 唯一标识符
	Tags            []string `gorm:"type:varchar[];default:NULL" json:"tags"`               // 代码功能标签 (using pq.StringArray for PostgreSQL array)
	ProjectName     string   `gorm:"type:varchar(128);not null" json:"project_name"`        // 所属项目
	Language        string   `gorm:"type:varchar(64);not null" json:"language"`             // 编程语言
	LanguageVersion string   `gorm:"type:varchar(32);default:NULL" json:"language_version"` // 语言版本
	FilePath        string   `gorm:"type:varchar(1024);not null" json:"file_path"`          // 代码文件路径
	FileName        string   `gorm:"type:varchar(256);not null" json:"file_name"`           // 代码文件名称
	FileType        string   `gorm:"type:varchar(32);not null" json:"file_type"`            // 代码文件类型
	Snippet         string   `gorm:"type:text;not null" json:"snippet"`                     // 代码片段内容
	Desc            string   `gorm:"type:varchar(4096);default:NULL" json:"desc"`           // 代码片段解释
	CodeRaw         string   `gorm:"type:text;default:NULL" json:"code_raw"`                // 原始代码文件 	// 版本号 	// 软删除时间 (可选)
}
