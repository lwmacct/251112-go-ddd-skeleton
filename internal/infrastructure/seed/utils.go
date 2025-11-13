package seed

import (
	"math/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

// generateULID 生成 ULID
func generateULID() string {
	entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
	return ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()
}

// stringPtr 返回字符串指针（辅助函数）
func stringPtr(s string) *string {
	return &s
}
