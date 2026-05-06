package pkg

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateOrderNo(t *testing.T) {
	t.Run("生成的订单号格式正确", func(t *testing.T) {
		orderNo := GenerateOrderNo()
		// 格式：ORD + 14位时间戳 + 6位随机数 = ORD + 20位数字
		assert.True(t, strings.HasPrefix(orderNo, "ORD"))
		assert.Len(t, orderNo, 23) // "ORD" + 20位数字
	})

	t.Run("连续生成的订单号不同（随机数部分不同）", func(t *testing.T) {
		orderNo1 := GenerateOrderNo()
		orderNo2 := GenerateOrderNo()
		// 虽然时间戳可能相同，但随机数大概率不同
		// 这里只验证不是完全相同（极小概率会失败，但可以接受）
		assert.NotEqual(t, orderNo1, orderNo2)
	})

	t.Run("生成100个订单号，全部唯一", func(t *testing.T) {
		orderNos := make(map[string]bool)
		for i := 0; i < 100; i++ {
			orderNo := GenerateOrderNo()
			assert.False(t, orderNos[orderNo], "订单号重复: %s", orderNo)
			orderNos[orderNo] = true
		}
	})
}
