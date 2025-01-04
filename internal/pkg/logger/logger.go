package logger

import (
	"fmt"
	"os"
	"sync"
)

// Logger 结构体表示日志记录器
type Logger struct {
	logFile *os.File
	mutex   sync.Mutex
}

// NewLogger 创建新的 Logger 实例
func NewLogger(filePath string) (*Logger, error) {
	logFile, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %v", err)
	}

	return &Logger{
		logFile: logFile,
	}, nil
}

// LogDetail 记录详细日志
func (l *Logger) LogDetail(text string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if _, err := l.logFile.WriteString(text + "\n"); err != nil {
		fmt.Printf("Failed to write log: %v\n", err)
	}
}

// Close 关闭日志文件
func (l *Logger) Close() {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if err := l.logFile.Close(); err != nil {
		fmt.Printf("Failed to close log file: %v\n", err)
	}
}
