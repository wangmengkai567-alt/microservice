import { computed } from 'vue';
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query';
import { productApi } from '@/services/api/product';
import { useUIStore } from '@/stores/ui';
import type { Product, CreateProductDto, UpdateProductDto } from '@/types/product';
import type { ProductFormData } from '@/types/product';

/**
 * 商品列表Composable
 * 使用TanStack Query实现缓存和自动重新验证
 */
export function useProducts(page: number = 1, pageSize: number = 10) {
  const uiStore = useUIStore();
  const queryClient = useQueryClient();

  // 查询商品列表，带缓存（5分钟有效期）
  const {
    data: productsData,
    isLoading,
    error: queryError,
    refetch,
  } = useQuery({
    queryKey: ['products', page, pageSize],
    queryFn: async () => {
      const response = await productApi.getProducts(page, pageSize);
      return response.data;
    },
    staleTime: 0, // 始终从服务器获取最新库存
    gcTime: 60 * 1000, // 1分钟后垃圾回收
  });

  // 计算属性
  const products = computed(() => productsData.value?.products || []);
  const total = computed(() => productsData.value?.total || 0);
  const currentPage = computed(() => productsData.value?.page || page);
  const error = computed(() => queryError.value?.message || null);

  // 创建商品Mutation
  const createMutation = useMutation({
    mutationFn: (data: ProductFormData) =>
      productApi.createProduct({
        name: data.name,
        description: data.description,
        price: data.price,
        stock: data.stock,
        image_url: data.imageUrl,
      }),
    onSuccess: () => {
      // 清除商品列表缓存，触发重新获取
      queryClient.invalidateQueries({ queryKey: ['products'] });
      uiStore.showSuccess('商品创建成功');
    },
    onError: (err: any) => {
      const message = err.response?.data?.message || '创建商品失败';
      uiStore.showError(message);
    },
  });

  // 更新商品Mutation
  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<ProductFormData> }) =>
      productApi.updateProduct(id, {
        name: data.name,
        description: data.description,
        price: data.price,
        stock: data.stock,
        image_url: data.imageUrl,
      }),
    onSuccess: (_response, { id }) => {
      // 清除商品列表缓存
      queryClient.invalidateQueries({ queryKey: ['products'] });
      // 更新单个商品缓存
      queryClient.invalidateQueries({
        queryKey: ['product', id],
      });
      uiStore.showSuccess('商品更新成功');
    },
    onError: (err: any) => {
      const message = err.response?.data?.message || '更新商品失败';
      uiStore.showError(message);
    },
  });

  // 删除商品Mutation
  const deleteMutation = useMutation({
    mutationFn: (id: string | number) => productApi.deleteProduct(String(id)),
    onSuccess: (_, id) => {
      // 清除商品列表缓存
      queryClient.invalidateQueries({ queryKey: ['products'] });
      // 移除单个商品缓存
      queryClient.removeQueries({ queryKey: ['product', id] });
      uiStore.showSuccess('商品删除成功');
    },
    onError: (err: any) => {
      const message = err.response?.data?.message || '删除商品失败';
      uiStore.showError(message);
    },
  });

  // 创建商品
  async function createProduct(data: CreateProductDto): Promise<Product> {
    const result = await createMutation.mutateAsync(data);
    return result.data;
  }

  // 更新商品
  async function updateProduct(id: string, data: UpdateProductDto): Promise<Product> {
    const result = await updateMutation.mutateAsync({ id, data });
    return result.data;
  }

  // 删除商品
  async function deleteProduct(id: string | number): Promise<void> {
    await deleteMutation.mutateAsync(id);
  }

  return {
    // 数据
    products,
    total,
    currentPage,
    isLoading,
    error,
    // 方法
    refetch,
    createProduct,
    updateProduct,
    deleteProduct,
    // Mutation状态
    isCreating: computed(() => createMutation.isPending.value),
    isUpdating: computed(() => updateMutation.isPending.value),
    isDeleting: computed(() => deleteMutation.isPending.value),
  };
}

/**
 * 单个商品详情Composable
 * 使用TanStack Query实现缓存
 */
export function useProductDetail(id: string) {
  const uiStore = useUIStore();

  // 查询单个商品详情，带缓存
  const {
    data: productData,
    isLoading,
    error: queryError,
    refetch,
  } = useQuery({
    queryKey: ['product', id],
    queryFn: async () => {
      const response = await productApi.getProduct(id);
      return response.data;
    },
    staleTime: 0, // 始终从服务器获取最新库存
    gcTime: 60 * 1000, // 1分钟后垃圾回收
    enabled: !!id,
  });

  // 计算属性
  const product = computed(() => productData.value || null);
  const error = computed(() => {
    if (queryError.value) {
      const err = queryError.value as any;
      if (err.response?.status === 404) {
        return '商品未找到';
      }
      return err.message || '获取商品详情失败';
    }
    return null;
  });

  // 监听错误并显示提示
  if (error.value) {
    uiStore.showError(error.value);
  }

  return {
    product,
    isLoading,
    error,
    refetch,
  };
}
