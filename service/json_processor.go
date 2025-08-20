package service

import (
	"context"
	"encoding/json"
	"regexp"
	"strings"

	"sojson/zlog"
)

var (
	JSONProcessorService = &jsonProcessorService{}
)

// jsonProcessorService JSON处理服务
type jsonProcessorService struct{}

// UnescapeJSON 去除JSON转义
func (s *jsonProcessorService) UnescapeJSON(text string) (string, error) {
	text = strings.TrimSpace(text)

	// 去除外层引号
	if strings.HasPrefix(text, `"`) && strings.HasSuffix(text, `"`) && len(text) > 1 {
		text = text[1 : len(text)-1]
	}

	// 只有在去除外层引号后才进行转义字符的处理
	text = strings.ReplaceAll(text, `\"`, `"`)
	text = strings.ReplaceAll(text, `\/`, `/`)
	text = strings.ReplaceAll(text, `\\`, `\`)
	text = strings.ReplaceAll(text, `\n`, "\n")
	text = strings.ReplaceAll(text, `\r`, "\r")
	text = strings.ReplaceAll(text, `\t`, "\t")
	text = strings.ReplaceAll(text, `\b`, "\b")
	text = strings.ReplaceAll(text, `\f`, "\f")

	return text, nil
}

// FormatJSON 格式化JSON
func (s *jsonProcessorService) FormatJSON(ctx context.Context, text string, indent int) (string, error) {
	var jsonObj interface{}

	// 解析JSON
	if err := json.Unmarshal([]byte(text), &jsonObj); err != nil {
		zlog.Errorf(ctx, "FormatJSON: json.Unmarshal failed, input text length: %d, error: %v", len(text), err)
		return "", err
	}

	// 格式化输出
	var result []byte
	var err error

	if indent <= 0 {
		result, err = json.Marshal(jsonObj)
		if err != nil {
			zlog.Errorf(ctx, "FormatJSON: json.Marshal failed, input text length: %d, error: %v", len(text), err)
			return "", err
		}
	} else {
		indentStr := strings.Repeat(" ", indent)
		result, err = json.MarshalIndent(jsonObj, "", indentStr)
		if err != nil {
			zlog.Errorf(ctx, "FormatJSON: json.MarshalIndent failed, input text length: %d, indent: %d, error: %v", len(text), indent, err)
			return "", err
		}
	}

	zlog.Infof(ctx, "FormatJSON: successfully formatted JSON, input length: %d, output length: %d, indent: %d", len(text), len(result), indent)
	return string(result), nil
}

// ProcessJSON 完整处理：先去除转义，再格式化
func (s *jsonProcessorService) ProcessJSON(ctx context.Context, text string, indent int) (string, error) {
	// 先修复未引号包裹的时间字段
	fixed := s.FixUnquotedTimeFields(text)
	zlog.Debugf(ctx, "ProcessJSON: FixUnquotedTimeFields, original length: %d, fixed length: %d", len(text), len(fixed))

	// 去除转义
	unescaped, err := s.UnescapeJSON(fixed)
	if err != nil {
		zlog.Errorf(ctx, "ProcessJSON: UnescapeJSON failed, input text length: %d, error: %v", len(fixed), err)
		return "", err
	}

	// 格式化
	formatted, err := s.FormatJSON(ctx, unescaped, indent)
	if err != nil {
		zlog.Errorf(ctx, "ProcessJSON: FormatJSON failed, unescaped text length: %d, indent: %d, error: %v", len(unescaped), indent, err)
		return "", err
	}

	zlog.Infof(ctx, "ProcessJSON: successfully processed JSON, input length: %d, output length: %d, indent: %d", len(text), len(formatted), indent)
	return formatted, nil
}

// FixUnquotedTimeFields 修复未被引号包裹的时间字段
func (s *jsonProcessorService) FixUnquotedTimeFields(text string) string {
	// 修复常见的时间字段格式问题
	// 匹配模式: "time":2025-08-18T08:04:19.827Z (没有引号的时间值)
	// 替换为: "time":"2025-08-18T08:04:19.827Z" (加上引号)

	// 时间格式的正则表达式模式
	// 匹配: "field_name":YYYY-MM-DDTHH:MM:SS.sssZ后面跟着逗号或右大括号
	timePattern := `("(?:time|timestamp|created_at|updated_at|date)"\s*:\s*)(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d{3})?Z?)(\s*[,}])`

	// 使用正则表达式替换，保留逗号或右大括号
	re := regexp.MustCompile(timePattern)
	result := re.ReplaceAllString(text, `${1}"${2}"${3}`)

	return result
}

// ProcessJSONWithTimeFix 专门用于处理包含未引号时间字段的JSON
func (s *jsonProcessorService) ProcessJSONWithTimeFix(ctx context.Context, text string, indent int) (string, error) {
	zlog.Infof(ctx, "ProcessJSONWithTimeFix: processing JSON with time field fix, input length: %d", len(text))

	// 先修复时间字段，然后正常处理
	return s.ProcessJSON(ctx, text, indent)
}

// ValidateJSON 验证JSON格式
func (s *jsonProcessorService) ValidateJSON(text string) error {
	var jsonObj interface{}
	return json.Unmarshal([]byte(text), &jsonObj)
}
