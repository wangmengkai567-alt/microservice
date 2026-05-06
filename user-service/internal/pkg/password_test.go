package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
	t.Run("正常哈希密码", func(t *testing.T) {
		hash, err := HashPassword("password123")
		require.NoError(t, err)
		assert.NotEmpty(t, hash)
		// bcrypt 哈希结果不等于原始密码
		assert.NotEqual(t, "password123", hash)
	})

	t.Run("相同密码每次哈希结果不同（bcrypt 内置随机 salt）", func(t *testing.T) {
		hash1, err1 := HashPassword("samepassword")
		hash2, err2 := HashPassword("samepassword")
		require.NoError(t, err1)
		require.NoError(t, err2)
		// 两次哈希结果不同，证明 salt 是随机的
		assert.NotEqual(t, hash1, hash2)
	})
}

func TestCheckPassword(t *testing.T) {
	password := "mypassword"
	hash, err := HashPassword(password)
	require.NoError(t, err)

	t.Run("正确密码验证通过", func(t *testing.T) {
		err := CheckPassword(hash, password)
		assert.NoError(t, err)
	})

	t.Run("错误密码验证失败", func(t *testing.T) {
		err := CheckPassword(hash, "wrongpassword")
		assert.Error(t, err)
	})

	t.Run("空密码验证失败", func(t *testing.T) {
		err := CheckPassword(hash, "")
		assert.Error(t, err)
	})
}
