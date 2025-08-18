package env

import (
	"os"
	"runtime"
	"strings"
)

// Environment 环境类型
type Environment string

const (
	// Test 测试环境
	Test Environment = "test"
	// Prod 生产环境
	Prod Environment = "prod"
)

var currentEnv Environment

// Init 初始化环境变量
func Init() {
	runEnv := os.Getenv("RUN_ENV")

	// 如果是Linux系统且没有设置RUN_ENV，直接panic退出
	if runtime.GOOS == "linux" && runEnv == "" {
		panic("RUN_ENV environment variable is required on Linux systems")
	}

	// 如果没有设置RUN_ENV，默认为test环境
	if runEnv == "" {
		currentEnv = Test
		return
	}

	// 根据RUN_ENV的值设置环境
	switch strings.ToLower(runEnv) {
	case "prod", "production":
		currentEnv = Prod
	case "test", "testing", "dev", "development":
		currentEnv = Test
	default:
		// 未知环境默认为test
		currentEnv = Test
	}
}

// GetEnv 获取当前环境
func GetEnv() Environment {
	return currentEnv
}

// IsTest 判断是否为测试环境
func IsTest() bool {
	return currentEnv == Test
}

// IsProd 判断是否为生产环境
func IsProd() bool {
	return currentEnv == Prod
}

// String 返回环境的字符串表示
func (e Environment) String() string {
	return string(e)
}
