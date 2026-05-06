package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"user-service/config"
)

func init() {
	// 测试时直接设置全局配置，不依赖配置文件
	config.Global.JWT.Secret = "test_secret_key"
	config.Global.JWT.Expire = 3600 // 1小时
}

func TestGenerateToken(t *testing.T) {
	t.Run("正常生成 Token", func(t *testing.T) {
		token, err := GenerateToken(1)
		require.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("不同 userID 生成不同 Token", func(t *testing.T) {
		token1, _ := GenerateToken(1)
		token2, _ := GenerateToken(2)
		assert.NotEqual(t, token1, token2)
	})
}

func TestParseToken(t *testing.T) {
	secret := config.Global.JWT.Secret

	t.Run("正常解析有效 Token", func(t *testing.T) {
		token, err := GenerateToken(42)
		require.NoError(t, err)

		claims, err := ParseToken(token, secret)
		require.NoError(t, err)
		assert.Equal(t, uint(42), claims.UserID)
	})

	t.Run("使用错误 secret 解析失败", func(t *testing.T) {
		token, err := GenerateToken(1)
		require.NoError(t, err)

		_, err = ParseToken(token, "wrong_secret")
		assert.Error(t, err)
	})

	t.Run("解析非法 Token 字符串失败", func(t *testing.T) {
		_, err := ParseToken("this.is.not.a.valid.token", secret)
		assert.Error(t, err)
	})

	t.Run("解析空字符串失败", func(t *testing.T) {
		_, err := ParseToken("", secret)
		assert.Error(t, err)
	})
}
