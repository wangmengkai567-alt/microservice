package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"user-service/config"
	"user-service/internal/model"
)

// ---- Mock 实现 ----

type mockUserRepo struct {
	users  map[string]*model.User // username -> user
	createErr error
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{users: make(map[string]*model.User)}
}

func (m *mockUserRepo) Create(user *model.User) error {
	if m.createErr != nil {
		return m.createErr
	}
	if _, exists := m.users[user.Username]; exists {
		return errors.New("username already exists")
	}
	user.ID = uint(len(m.users) + 1) // 模拟自增 ID
	m.users[user.Username] = user
	return nil
}

func (m *mockUserRepo) FindByUsername(username string) (*model.User, error) {
	user, ok := m.users[username]
	if !ok {
		return nil, errors.New("record not found")
	}
	return user, nil
}

// ---- 测试 ----

func init() {
	config.Global.JWT.Secret = "test_secret"
	config.Global.JWT.Expire = 3600
}

func TestUserService_Register(t *testing.T) {
	t.Run("注册成功", func(t *testing.T) {
		repo := newMockUserRepo()
		svc := NewUserService(repo)

		err := svc.Register("alice", "alice@example.com", "password123")
		assert.NoError(t, err)
		// 验证用户确实被存入
		user, _ := repo.FindByUsername("alice")
		assert.Equal(t, "alice", user.Username)
		// 密码应该被哈希，不是明文
		assert.NotEqual(t, "password123", user.Password)
	})

	t.Run("用户名为空时注册失败", func(t *testing.T) {
		svc := NewUserService(newMockUserRepo())
		err := svc.Register("", "a@example.com", "password123")
		assert.EqualError(t, err, "username, email and password cannot be empty")
	})

	t.Run("邮箱为空时注册失败", func(t *testing.T) {
		svc := NewUserService(newMockUserRepo())
		err := svc.Register("alice", "", "password123")
		assert.EqualError(t, err, "username, email and password cannot be empty")
	})

	t.Run("密码少于6位时注册失败", func(t *testing.T) {
		svc := NewUserService(newMockUserRepo())
		err := svc.Register("alice", "alice@example.com", "123")
		assert.EqualError(t, err, "password must be at least 6 characters")
	})

	t.Run("用户名重复时注册失败", func(t *testing.T) {
		repo := newMockUserRepo()
		svc := NewUserService(repo)

		require.NoError(t, svc.Register("alice", "alice@example.com", "password123"))
		err := svc.Register("alice", "alice2@example.com", "password123")
		assert.Error(t, err)
	})

	t.Run("数据库错误时注册失败", func(t *testing.T) {
		repo := newMockUserRepo()
		repo.createErr = errors.New("db connection lost")
		svc := NewUserService(repo)

		err := svc.Register("alice", "alice@example.com", "password123")
		assert.Error(t, err)
	})
}

func TestUserService_Login(t *testing.T) {
	// 预先注册一个用户
	repo := newMockUserRepo()
	svc := NewUserService(repo)
	require.NoError(t, svc.Register("bob", "bob@example.com", "correctpass"))

	t.Run("正确用户名和密码登录成功，返回 Token", func(t *testing.T) {
		token, err := svc.Login("bob", "correctpass")
		require.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("用户不存在时登录失败", func(t *testing.T) {
		_, err := svc.Login("nobody", "password")
		assert.EqualError(t, err, "user not found")
	})

	t.Run("密码错误时登录失败", func(t *testing.T) {
		_, err := svc.Login("bob", "wrongpass")
		assert.EqualError(t, err, "password incorrect")
	})

	t.Run("用户名为空时登录失败", func(t *testing.T) {
		_, err := svc.Login("", "password")
		assert.EqualError(t, err, "username and password cannot be empty")
	})

	t.Run("密码为空时登录失败", func(t *testing.T) {
		_, err := svc.Login("bob", "")
		assert.EqualError(t, err, "username and password cannot be empty")
	})
}
