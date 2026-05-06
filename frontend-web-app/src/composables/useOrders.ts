import { computed } from 'vue';
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query';
import { orderApi } from '@/services/api/order';
import { useUIStore } from '@/stores/ui';
import type { Order, CreateOrderDto } from '@/types/order';

/**
 * 订单列表Composable
 * 使用TanStack Query实现缓存
 */
export function useOrders(page: number = 1, pageSize: number = 10) {
  const uiStore = useUIStore();
  const queryClient = useQueryClient();

  // 查询订单列表，带缓存
  const {
    data: ordersData,
    isLoading,
    error: queryError,
    refetch,
  } = useQuery({
    queryKey: ['orders', page, pageSize],
    queryFn: async () => {
      const response = await orderApi.getOrders(page, pageSize);
      return response.data;
    },
    staleTime: 2 * 60 * 1000, // 2分钟缓存有效期（订单数据更新频繁）
    gcTime: 5 * 60 * 1000, // 5分钟后垃圾回收
  });

  // 计算属性 - 按创建时间倒序排列
  const orders = computed(() => {
    const orderList = ordersData.value?.orders || [];
    // 确保按创建时间倒序排列（最新的在前）
    return [...orderList].sort((a, b) => {
      return new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime();
    });
  });

  const total = computed(() => ordersData.value?.total || 0);
  const currentPage = computed(() => ordersData.value?.page || page);
  const error = computed(() => queryError.value?.message || null);

  // 创建订单Mutation
  const createMutation = useMutation({
    mutationFn: (data: CreateOrderDto) =>
      orderApi.createOrder({
        product_id: Number(data.productId),
        quantity: data.quantity,
      } as any),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['orders'] });
      // 库存已扣减，刷新商品缓存
      queryClient.invalidateQueries({ queryKey: ['products'] });
      queryClient.invalidateQueries({ queryKey: ['product'] });
      uiStore.showSuccess('订单创建成功');
    },
    onError: (err: any) => {
      const message = err.response?.data?.message || '创建订单失败';
      if (message.includes('库存') || message.includes('stock')) {
        uiStore.showError('库存不足，无法创建订单');
      } else {
        uiStore.showError(message);
      }
    },
  });

  // 取消订单Mutation
  const cancelMutation = useMutation({
    mutationFn: (id: string) => orderApi.cancelOrder(id),
    onSuccess: (_response, id) => {
      queryClient.invalidateQueries({ queryKey: ['orders'] });
      queryClient.invalidateQueries({ queryKey: ['order', id] });
      // 库存已恢复，刷新商品缓存
      queryClient.invalidateQueries({ queryKey: ['products'] });
      queryClient.invalidateQueries({ queryKey: ['product'] });
      uiStore.showSuccess('订单取消成功');
    },
    onError: (err: any) => {
      const message = err.response?.data?.message || '取消订单失败';
      // 特殊处理业务规则错误
      if (
        message.includes('无法取消') ||
        message.includes('cannot cancel') ||
        message.includes('not pending')
      ) {
        uiStore.showError('无法取消该订单，订单状态不是待处理');
      } else {
        uiStore.showError(message);
      }
    },
  });

  // 支付订单Mutation
  const payMutation = useMutation({
    mutationFn: (id: string) => orderApi.payOrder(id),
    onSuccess: (_response, id) => {
      queryClient.invalidateQueries({ queryKey: ['orders'] });
      queryClient.invalidateQueries({ queryKey: ['order', id] });
      uiStore.showSuccess('订单支付成功');
    },
    onError: (err: any) => {
      const message = err.response?.data?.message || '支付订单失败';
      uiStore.showError(message);
    },
  });

  // 删除订单Mutation
  const deleteMutation = useMutation({
    mutationFn: (id: string) => orderApi.deleteOrder(id),
    onSuccess: (_response, id) => {
      queryClient.invalidateQueries({ queryKey: ['orders'] });
      queryClient.removeQueries({ queryKey: ['order', id] });
      uiStore.showSuccess('订单删除成功');
    },
    onError: (err: any) => {
      const message = err.response?.data?.message || '删除订单失败';
      if (message.includes('cancelled')) {
        uiStore.showError('仅支持删除已取消的订单');
      } else {
        uiStore.showError(message);
      }
    },
  });

  // 创建订单
  async function createOrder(data: CreateOrderDto): Promise<Order> {
    const result = await createMutation.mutateAsync(data);
    return result.data;
  }

  // 取消订单
  async function cancelOrder(id: string): Promise<void> {
    await cancelMutation.mutateAsync(id);
  }

  // 支付订单
  async function payOrder(id: string): Promise<void> {
    await payMutation.mutateAsync(id);
  }

  async function deleteOrder(id: string): Promise<void> {
    await deleteMutation.mutateAsync(id);
  }

  return {
    // 数据
    orders,
    total,
    currentPage,
    isLoading,
    error,
    // 方法
    refetch,
    createOrder,
    cancelOrder,
    payOrder,
    deleteOrder,
    // Mutation状态
    isCreating: computed(() => createMutation.isPending.value),
    isCancelling: computed(() => cancelMutation.isPending.value),
    isPaying: computed(() => payMutation.isPending.value),
    isDeleting: computed(() => deleteMutation.isPending.value),
  };
}

/**
 * 单个订单详情Composable
 * 使用TanStack Query实现缓存
 */
export function useOrderDetail(id: string) {
  const uiStore = useUIStore();
  const queryClient = useQueryClient();

  // 查询单个订单详情，带缓存
  const {
    data: orderData,
    isLoading,
    error: queryError,
    refetch,
  } = useQuery({
    queryKey: ['order', id],
    queryFn: async () => {
      const response = await orderApi.getOrder(id);
      return response.data;
    },
    staleTime: 2 * 60 * 1000, // 2分钟缓存有效期
    gcTime: 5 * 60 * 1000, // 5分钟后垃圾回收
    enabled: !!id, // 只有当id存在时才执行查询
  });

  // 计算属性
  const order = computed(() => orderData.value || null);
  const error = computed(() => {
    if (queryError.value) {
      const err = queryError.value as any;
      if (err.response?.status === 404) {
        return '订单未找到';
      }
      return err.message || '获取订单详情失败';
    }
    return null;
  });

  // 判断订单是否可以取消（状态为pending）
  const canCancel = computed(() => {
    return order.value?.status === 'pending';
  });
  // 判断订单是否可以支付（状态为pending）
  const canPay = computed(() => {
    return order.value?.status === 'pending';
  });
  const canDelete = computed(() => {
    return order.value?.status === 'cancelled';
  });

  // 监听错误并显示提示
  if (error.value) {
    uiStore.showError(error.value);
  }

  // 取消订单Mutation（用于详情页）
  const cancelMutation = useMutation({
    mutationFn: () => orderApi.cancelOrder(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['orders'] });
      queryClient.invalidateQueries({ queryKey: ['order', id] });
      // 库存已恢复，刷新商品缓存
      queryClient.invalidateQueries({ queryKey: ['products'] });
      queryClient.invalidateQueries({ queryKey: ['product'] });
      uiStore.showSuccess('订单取消成功');
    },
    onError: (err: any) => {
      const message = err.response?.data?.message || '取消订单失败';
      if (
        message.includes('无法取消') ||
        message.includes('cannot cancel') ||
        message.includes('not pending')
      ) {
        uiStore.showError('无法取消该订单，订单状态不是待处理');
      } else {
        uiStore.showError(message);
      }
    },
  });

  const deleteMutation = useMutation({
    mutationFn: () => orderApi.deleteOrder(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['orders'] });
      queryClient.removeQueries({ queryKey: ['order', id] });
      uiStore.showSuccess('订单删除成功');
    },
    onError: (err: any) => {
      const message = err.response?.data?.message || '删除订单失败';
      uiStore.showError(message);
    },
  });

  // 支付订单Mutation
  const payMutation = useMutation({
    mutationFn: () => orderApi.payOrder(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['orders'] });
      queryClient.invalidateQueries({ queryKey: ['order', id] });
      uiStore.showSuccess('订单支付成功');
    },
    onError: (err: any) => {
      const message = err.response?.data?.message || '支付订单失败';
      uiStore.showError(message);
    },
  });

  // 取消订单
  async function cancelOrder(): Promise<void> {
    await cancelMutation.mutateAsync();
  }

  // 支付订单
  async function payOrder(): Promise<void> {
    await payMutation.mutateAsync();
  }

  async function deleteOrder(): Promise<void> {
    await deleteMutation.mutateAsync();
  }

  return {
    order,
    isLoading,
    error,
    canCancel,
    canPay,
    canDelete,
    refetch,
    cancelOrder,
    payOrder,
    deleteOrder,
    isCancelling: computed(() => cancelMutation.isPending.value),
    isPaying: computed(() => payMutation.isPending.value),
    isDeleting: computed(() => deleteMutation.isPending.value),
  };
}
