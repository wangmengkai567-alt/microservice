// 订单类型定义

export type OrderStatus = 'pending' | 'paid' | 'completed' | 'cancelled';

export interface Order {
  id: string;
  orderNo: string; // 订单号
  userId: string;
  productId: string;
  productName: string;
  quantity: number;
  totalPrice: number;
  status: OrderStatus;
  createdAt: string;
  updatedAt: string;
}

export interface OrderFormData {
  productId: string | number;
  quantity: number;
}

export interface CreateOrderDto {
  productId: string | number;
  quantity: number;
}

export interface OrderListResponse {
  orders: Order[];
  total: number;
  page: number;
  pageSize: number;
}
