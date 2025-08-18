package service

import (
	"encoding/json"
	"strings"
)

// JSONProcessorService JSON处理服务
type JSONProcessorService struct{}

// NewJSONProcessorService 创建JSON处理服务实例
func NewJSONProcessorService() *JSONProcessorService {
	return &JSONProcessorService{}
}

// UnescapeJSON 去除JSON转义
func (s *JSONProcessorService) UnescapeJSON(text string) (string, error) {
	// 去除常见的转义字符
	text = strings.ReplaceAll(text, `\"`, `"`)
	text = strings.ReplaceAll(text, `\/`, `/`)
	text = strings.ReplaceAll(text, `\\`, `\`)
	text = strings.ReplaceAll(text, `\n`, "\n")
	text = strings.ReplaceAll(text, `\r`, "\r")
	text = strings.ReplaceAll(text, `\t`, "\t")
	text = strings.ReplaceAll(text, `\b`, "\b")
	text = strings.ReplaceAll(text, `\f`, "\f")

	// 去除外层引号（如果存在）
	text = strings.TrimSpace(text)
	if strings.HasPrefix(text, `"`) && strings.HasSuffix(text, `"`) && len(text) > 1 {
		text = text[1 : len(text)-1]
	}

	return text, nil
}

// FormatJSON 格式化JSON
func (s *JSONProcessorService) FormatJSON(text string, indent int) (string, error) {
	var jsonObj interface{}

	// 解析JSON
	if err := json.Unmarshal([]byte(text), &jsonObj); err != nil {
		return "", err
	}

	// 格式化输出
	var result []byte
	var err error

	if indent <= 0 {
		result, err = json.Marshal(jsonObj)
	} else {
		indentStr := strings.Repeat(" ", indent)
		result, err = json.MarshalIndent(jsonObj, "", indentStr)
	}

	if err != nil {
		return "", err
	}

	return string(result), nil
}

// ProcessJSON 完整处理：先去除转义，再格式化
func (s *JSONProcessorService) ProcessJSON(text string, indent int) (string, error) {
	// 先去除转义
	unescaped, err := s.UnescapeJSON(text)
	if err != nil {
		return "", err
	}

	// 再格式化
	formatted, err := s.FormatJSON(unescaped, indent)
	if err != nil {
		return "", err
	}

	return formatted, nil
}

// ValidateJSON 验证JSON格式
func (s *JSONProcessorService) ValidateJSON(text string) error {
	var jsonObj interface{}
	return json.Unmarshal([]byte(text), &jsonObj)
}
