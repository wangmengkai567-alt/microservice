import apiClient from './client';
import type { Order, CreateOrderDto, OrderListResponse } from '@/types/order';

function normalizeOrder(raw: any): Order {
  return {
    id: String(raw.id),
    orderNo: raw.orderNo ?? raw.order_no ?? '',
    userId: String(raw.userId ?? raw.user_id ?? ''),
    productId: String(raw.productId ?? raw.product_id ?? ''),
    productName: raw.productName ?? raw.product_name ?? '',
    quantity: Number(raw.quantity ?? 0),
    totalPrice: Number(raw.totalPrice ?? raw.total_price ?? 0),
    status: raw.status ?? 'pending',
    createdAt: raw.createdAt ?? raw.created_at ?? '',
    updatedAt: raw.updatedAt ?? raw.updated_at ?? '',
  };
}

/**
 * 订单API服务
 */
export const orderApi = {
  /**
   * 获取订单列表
   */
  getOrders: (page: number = 1, pageSize: number = 10) => {
    return apiClient
      .get('/api/v1/orders', {
        params: { page, page_size: pageSize },
      })
      .then((resp) => {
        const data = resp.data || {};
        const orders = Array.isArray(data.orders) ? data.orders.map(normalizeOrder) : [];
        const result: OrderListResponse = {
          orders,
          total: Number(data.total ?? 0),
          page,
          pageSize,
        };
        return { ...resp, data: result };
      });
  },

  /**
   * 获取单个订单详情
   */
  getOrder: (id: string) => {
    return apiClient.get(`/api/v1/orders/${id}`).then((resp) => {
      return { ...resp, data: normalizeOrder(resp.data) };
    });
  },

  /**
   * 创建订单
   */
  createOrder: (data: CreateOrderDto) => {
    return apiClient.post('/api/v1/orders', data).then((resp) => {
      return { ...resp, data: normalizeOrder(resp.data) };
    });
  },

  /**
   * 取消订单
   */
  cancelOrder: (id: string) => {
    return apiClient.post<{ message: string }>(`/api/v1/orders/${id}/cancel`);
  },

  /**
   * 支付订单
   */
  payOrder: (id: string) => {
    return apiClient.post<{ message: string }>(`/api/v1/orders/${id}/pay`);
  },

  /**
   * 删除订单（已取消订单）
   */
  deleteOrder: (id: string) => {
    return apiClient.delete<{ message: string }>(`/api/v1/orders/${id}`);
  },
};
