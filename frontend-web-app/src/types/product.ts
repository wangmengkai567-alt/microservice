// 商品类型定义

export interface Product {
  id: string | number;
  name: string;
  description: string;
  price: number;
  stock: number;
  imageUrl?: string;
  image_url?: string;
  createdAt?: string;
  updatedAt?: string;
}

export interface ProductFormData {
  name: string;
  description: string;
  price: number;
  stock: number;
  imageUrl?: string;
}

export interface CreateProductDto {
  name: string;
  description: string;
  price: number;
  stock: number;
  image_url?: string;
}

export interface UpdateProductDto {
  name?: string;
  description?: string;
  price?: number;
  stock?: number;
  image_url?: string;
}

export interface ProductListResponse {
  products: Product[];
  total: number;
  page: number;
  pageSize: number;
}
