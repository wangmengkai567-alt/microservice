import axios, { AxiosInstance, AxiosError, InternalAxiosRequestConfig } from 'axios';
import { getAuthToken, clearAuthData } from '@/utils/storage';
import type { ApiError } from '@/types/api';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

// 创建Axios实例
const apiClient: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// 请求拦截器：自动添加认证Token
apiClient.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const token = getAuthToken();
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error: AxiosError) => {
    return Promise.reject(error);
  }
);

// 响应拦截器：统一处理错误
apiClient.interceptors.response.use(
  (response) => {
    return response;
  },
  (error: AxiosError<ApiError>) => {
    const { response } = error;

    // 网络错误
    if (!response) {
      const networkError: ApiError = {
        message: '网络连接失败，请检查网络设置',
        code: 0,
      };
      console.error('Network Error:', error);
      return Promise.reject(networkError);
    }

    // 认证错误 - 401未授权
    if (response.status === 401) {
      // 清除认证数据
      clearAuthData();

      // 跳转到登录页面
      if (typeof window !== 'undefined') {
        window.location.href = '/login';
      }

      const authError: ApiError = {
        message: '登录已过期，请重新登录',
        code: 401,
        status: 401,
      };
      console.error('Authentication Error:', authError);
      return Promise.reject(authError);
    }

    // 服务器错误 - 500系列
    if (response.status >= 500) {
      const serverError: ApiError = {
        message: '服务器错误，请稍后重试',
        code: response.status,
        status: response.status,
      };
      console.error('Server Error:', error);
      return Promise.reject(serverError);
    }

    // 业务错误 - 4xx系列
    const businessError: ApiError = {
      message: response.data?.message || '操作失败',
      code: response.data?.code || response.status,
      status: response.status,
      errors: response.data?.errors,
    };
    console.error('Business Error:', businessError);
    return Promise.reject(businessError);
  }
);

export default apiClient;
