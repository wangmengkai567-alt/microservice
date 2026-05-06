import apiClient from './client';
import type { AxiosProgressEvent } from 'axios';
import type {
  Product,
  CreateProductDto,
  UpdateProductDto,
  ProductListResponse,
} from '@/types/product';

/**
 * 商品API服务
 */
export const productApi = {
  /**
   * 获取商品列表
   */
  getProducts: (page: number = 1, pageSize: number = 10) => {
    return apiClient.get<ProductListResponse>('/api/v1/products', {
      params: { page, page_size: pageSize },
    });
  },

  /**
   * 获取单个商品详情
   */
  getProduct: (id: string) => {
    return apiClient.get<Product>(`/api/v1/products/${id}`);
  },

  /**
   * 创建商品
   */
  createProduct: (data: CreateProductDto) => {
    return apiClient.post<Product>('/api/v1/products', data);
  },

  /**
   * 更新商品
   */
  updateProduct: (id: string, data: UpdateProductDto) => {
    return apiClient.put<Product>(`/api/v1/products/${id}`, data);
  },

  /**
   * 删除商品
   */
  deleteProduct: (id: string) => {
    return apiClient.delete<void>(`/api/v1/products/${id}`);
  },

  uploadImage: (file: File, onProgress?: (percent: number) => void) => {
    const formData = new FormData();
    formData.append('image', file);
    return apiClient.post<{ image_url: string }>('/api/v1/products/upload-image', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
      onUploadProgress: (event: AxiosProgressEvent) => {
        if (!onProgress || !event.total) return;
        const percent = Math.min(100, Math.round((event.loaded / event.total) * 100));
        onProgress(percent);
      },
    });
  },
};
