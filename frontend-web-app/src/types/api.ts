// API响应和错误类型定义

export interface ApiResponse<T = any> {
  data: T;
  message?: string;
  code?: number;
}

export interface ApiError {
  message: string;
  code?: number;
  status?: number;
  errors?: Record<string, string[]>;
}

export interface PaginationParams {
  page?: number;
  pageSize?: number;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  pageSize: number;
  totalPages: number;
}
