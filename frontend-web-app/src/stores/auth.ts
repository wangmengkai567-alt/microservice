import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import { userApi } from '@/services/api/user';
import {
  setAuthToken,
  getAuthToken,
  clearAuthData,
  setUserInfo,
  getUserInfo,
} from '@/utils/storage';
import type { User, LoginDto, RegisterDto } from '@/types/user';

export const useAuthStore = defineStore('auth', () => {
  // 状态
  const token = ref<string | null>(null);
  const user = ref<User | null>(null);

  // 计算属性
  const isAuthenticated = computed(() => !!token.value && !!user.value);

  // 初始化：从localStorage恢复认证状态
  function initAuth() {
    const savedToken = getAuthToken();
    const savedUser = getUserInfo<User>();

    if (savedToken && savedUser) {
      // 验证Token格式和是否过期
      try {
        // 检查 token 格式是否为 JWT (应该有3个部分，用.分隔)
        const parts = savedToken.split('.');
        if (parts.length !== 3) {
          console.warn('Invalid token format: not a valid JWT');
          clearAuth();
          return;
        }

        // 解析JWT Token的payload部分
        const payload = JSON.parse(atob(parts[1]));

        // 检查是否有过期时间
        if (payload.exp) {
          const exp = payload.exp * 1000; // 转换为毫秒
          const now = Date.now();

          if (exp > now) {
            // Token未过期
            token.value = savedToken;
            user.value = savedUser;
          } else {
            // Token已过期，清除数据
            console.warn('Token expired');
            clearAuth();
          }
        } else {
          // 没有过期时间，直接使用
          token.value = savedToken;
          user.value = savedUser;
        }
      } catch (error) {
        // Token解析失败，清除数据
        console.error('Failed to parse token:', error);
        clearAuth();
      }
    }
  }

  // 设置认证信息
  function setAuth(newToken: string, newUser: User) {
    token.value = newToken;
    user.value = newUser;
    setAuthToken(newToken);
    setUserInfo(newUser);
  }

  // 清除认证信息
  function clearAuth() {
    token.value = null;
    user.value = null;
    clearAuthData();
  }

  // 登录
  async function login(credentials: LoginDto) {
    try {
      const response = await userApi.login(credentials);
      const { token: newToken, user: newUser } = response.data;
      setAuth(newToken, newUser);
      return response.data;
    } catch (error) {
      clearAuth();
      throw error;
    }
  }

  // 注册
  async function register(data: RegisterDto) {
    try {
      await userApi.register(data);
    } catch (error) {
      throw error;
    }
  }

  // 退出登录
  function logout() {
    clearAuth();
  }

  // 初始化认证状态
  initAuth();

  return {
    // 状态
    token,
    user,
    isAuthenticated,
    // 方法
    setAuth,
    clearAuth,
    login,
    register,
    logout,
  };
});
