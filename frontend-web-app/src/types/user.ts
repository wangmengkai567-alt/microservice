// 用户类型定义

export interface User {
  id: string;
  username: string;
  createdAt: string;
}

export interface AuthState {
  token: string | null;
  user: User | null;
  isAuthenticated: boolean;
}

export interface RegisterDto {
  username: string;
  password: string;
}

export interface LoginDto {
  username: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  user: User;
}
