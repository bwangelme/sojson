package service

import (
	"context"
	"testing"

	"sojson/zlog"
)

func init() {
	// 初始化日志系统用于测试
	zlog.InitLogger("debug", "")
}

func TestFixUnquotedTimeFields(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "单个time字段",
			input:    `{"level":30,"time":2025-08-18T08:04:19.827Z}`,
			expected: `{"level":30,"time":"2025-08-18T08:04:19.827Z"}`,
		},
		{
			name:     "time字段在中间",
			input:    `{"level":30,"time":2025-08-18T08:04:19.827Z,"msg":"test"}`,
			expected: `{"level":30,"time":"2025-08-18T08:04:19.827Z","msg":"test"}`,
		},
		{
			name:     "多个时间字段",
			input:    `{"created_at":2025-08-18T08:04:19.827Z,"updated_at":2025-08-18T09:04:19.827Z}`,
			expected: `{"created_at":"2025-08-18T08:04:19.827Z","updated_at":"2025-08-18T09:04:19.827Z"}`,
		},
		{
			name:     "timestamp字段",
			input:    `{"level":30,"timestamp":2025-08-18T08:04:19Z,"msg":"test"}`,
			expected: `{"level":30,"timestamp":"2025-08-18T08:04:19Z","msg":"test"}`,
		},
		{
			name:     "已经有引号的时间字段（不应该改变）",
			input:    `{"level":30,"time":"2025-08-18T08:04:19.827Z","msg":"test"}`,
			expected: `{"level":30,"time":"2025-08-18T08:04:19.827Z","msg":"test"}`,
		},
		{
			name:     "没有时间字段（不应该改变）",
			input:    `{"level":30,"msg":"test message"}`,
			expected: `{"level":30,"msg":"test message"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := JSONProcessorService.FixUnquotedTimeFields(tt.input)
			if result != tt.expected {
				t.Errorf("FixUnquotedTimeFields() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestProcessJSONWithUnquotedTime(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		input       string
		indent      int
		shouldError bool
	}{
		{
			name:        "有效的未引号时间JSON",
			input:       `{"level":30,"time":2025-08-18T08:04:19.827Z,"msg":"test"}`,
			indent:      2,
			shouldError: false,
		},
		{
			name:        "多个未引号时间字段",
			input:       `{"created_at":2025-08-18T08:04:19.827Z,"updated_at":2025-08-18T09:04:19.827Z,"level":30}`,
			indent:      2,
			shouldError: false,
		},
		{
			name:        "无效JSON（即使修复后也无效）",
			input:       `{"level":30,"time":2025-08-18T08:04:19.827Z,"msg":}`,
			indent:      2,
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := JSONProcessorService.ProcessJSON(ctx, tt.input, tt.indent)

			if tt.shouldError {
				if err == nil {
					t.Errorf("ProcessJSON() expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("ProcessJSON() unexpected error = %v", err)
				}
				if result == "" {
					t.Errorf("ProcessJSON() returned empty result")
				}
				// 验证结果是有效的JSON
				if err := JSONProcessorService.ValidateJSON(result); err != nil {
					t.Errorf("ProcessJSON() result is not valid JSON: %v", err)
				}
			}
		})
	}
}

func TestProcessJSONWithTimeFix(t *testing.T) {
	ctx := context.Background()

	// 测试原始问题的JSON
	problemJSON := `{"level":30,"time":2025-08-18T08:04:19.827Z}`

	result, err := JSONProcessorService.ProcessJSONWithTimeFix(ctx, problemJSON, 2)
	if err != nil {
		t.Fatalf("ProcessJSONWithTimeFix() failed: %v", err)
	}

	// 验证结果是有效的JSON
	if err := JSONProcessorService.ValidateJSON(result); err != nil {
		t.Errorf("ProcessJSONWithTimeFix() result is not valid JSON: %v", err)
	}

	// 结果应该包含格式化的JSON
	if len(result) == 0 {
		t.Error("ProcessJSONWithTimeFix() returned empty result")
	}

	t.Logf("原始JSON: %s", problemJSON)
	t.Logf("处理结果: %s", result)
}
