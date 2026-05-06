package pkg

import (
	"fmt"
	"time"
	"math/rand"
)

// GenerateOrderNo 生成订单号
// 格式: ORD + 年月日时分秒 + 6位随机数
// 例如: ORD20240308153045123456
func GenerateOrderNo() string {
	now := time.Now()
	timestamp := now.Format("20060102150405")
	random := rand.Intn(900000) + 100000 // 生成 100000-999999 的随机数
	return fmt.Sprintf("ORD%s%d", timestamp, random)
}
