package utils

import (
	"math"
	"time"
)

// 辅助函数，用于处理可选字段的指针类型
func GetTimePointer(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}
	return &t
}

func GetFloat64Pointer(f float64) *float64 {
	// 检查是否为NaN或Inf，如果是则返回nil
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return nil
	}
	if f == 0 {
		return nil
	}
	return &f
}

func GetIntPointer(i int) *int {
	if i == 0 {
		return nil
	}
	return &i
}

func GetStringPointer(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
