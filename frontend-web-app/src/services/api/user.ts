import apiClient from './client';
import type { RegisterDto, LoginDto, LoginResponse, User } from '@/types/user';

/**
 * 用户API服务
 */
export const userApi = {
  /**
   * 用户注册
   */
  register: (data: RegisterDto) => {
    return apiClient.post<void>('/api/v1/users/register', data);
  },

  /**
   * 用户登录
   */
  login: (data: LoginDto) => {
    return apiClient.post<LoginResponse>('/api/v1/users/login', data);
  },

  /**
   * 获取当前用户信息
   */
  getCurrentUser: () => {
    return apiClient.get<User>('/api/v1/users/me');
  },
};
