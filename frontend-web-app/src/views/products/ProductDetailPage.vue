<script setup lang="ts">
import { ref, computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import MainLayout from '@/components/layout/MainLayout.vue';
import Button from '@/components/common/Button.vue';
import Loading from '@/components/common/Loading.vue';
import Modal from '@/components/common/Modal.vue';
import OrderForm from '@/components/features/OrderForm.vue';
import { useProductDetail, useProducts } from '@/composables/useProducts';
import { useOrders } from '@/composables/useOrders';
import { useAuthStore } from '@/stores/auth';
import { formatPrice, formatDate } from '@/utils/format';
import type { OrderFormData } from '@/types/order';

const route = useRoute();
const router = useRouter();
const authStore = useAuthStore();

const productId = computed(() => route.params.id as string);
const { product, isLoading, error } = useProductDetail(productId.value);
const { deleteProduct, isDeleting } = useProducts();
const { createOrder, isCreating } = useOrders();

// Order modal
const showOrderModal = ref(false);

// Delete confirmation modal
const showDeleteModal = ref(false);

function handleBuyNow() {
  if (!authStore.isAuthenticated) {
    router.push({ path: '/login', query: { redirect: route.fullPath } });
    return;
  }
  showOrderModal.value = true;
}

async function handleOrderSubmit(data: OrderFormData) {
  try {
    const order = await createOrder(data);
    showOrderModal.value = false;
    router.push(`/orders/${order.id}`);
  } catch (err) {
    console.error('Failed to create order:', err);
  }
}

function handleEdit() {
  router.push(`/products/${productId.value}/edit`);
}

function confirmDelete() {
  showDeleteModal.value = true;
}

async function handleDelete() {
  try {
    await deleteProduct(productId.value);
    showDeleteModal.value = false;
    router.push('/products');
  } catch (err) {
    console.error('Failed to delete product:', err);
  }
}
</script>

<template>
  <MainLayout>
    <!-- 返回按钮 -->
    <button
      @click="router.back()"
      class="mb-6 inline-flex items-center text-sm font-medium text-gray-500 hover:text-indigo-600 transition-colors focus:outline-none focus:ring-2 focus:ring-indigo-500 rounded-lg px-2 py-1 -ml-2"
    >
      <svg class="w-5 h-5 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M10 19l-7-7m0 0l7-7m-7 7h18"
        />
      </svg>
      返回列表
    </button>

    <!-- 加载状态 -->
    <div v-if="isLoading" class="flex justify-center py-12">
      <Loading />
    </div>

    <!-- 错误提示 -->
    <div v-else-if="error" class="bg-red-50 border border-red-200 rounded-md p-4">
      <p class="text-sm text-red-800">{{ error }}</p>
      <Button @click="router.push('/products')" variant="secondary" class="mt-4">
        返回商品列表
      </Button>
    </div>

    <!-- 商品详情 -->
    <div
      v-else-if="product"
      class="bg-white shadow-xl rounded-2xl border border-gray-100 overflow-hidden"
    >
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-0 lg:gap-8">
        <!-- 左侧图片区域 -->
        <div class="bg-gray-50 h-80 sm:h-96 lg:h-full min-h-[400px] relative">
          <img
            v-if="product.image_url || product.imageUrl"
            :src="product.image_url || product.imageUrl"
            :alt="product.name"
            class="w-full h-full object-cover"
          />
          <div v-else class="absolute inset-0 flex items-center justify-center text-gray-300">
            <svg class="w-24 h-24" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="1"
                d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"
              />
            </svg>
          </div>
        </div>

        <!-- 右侧信息区域 -->
        <div class="p-8 lg:p-10 flex flex-col h-full">
          <div class="flex justify-between items-start mb-4">
            <h1 class="text-3xl sm:text-4xl font-extrabold text-gray-900 tracking-tight">
              {{ product.name }}
            </h1>
            <div class="flex space-x-2 shrink-0 ml-4">
              <button
                @click="handleEdit"
                class="p-2 text-gray-400 hover:text-indigo-600 hover:bg-indigo-50 rounded-full transition-colors focus:outline-none focus:ring-2 focus:ring-indigo-500"
                title="编辑"
              >
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"
                  ></path>
                </svg>
              </button>
              <button
                @click="confirmDelete"
                class="p-2 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-full transition-colors focus:outline-none focus:ring-2 focus:ring-red-500"
                title="删除"
              >
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                  ></path>
                </svg>
              </button>
            </div>
          </div>

          <div class="mb-6 flex items-baseline gap-4">
            <p class="text-4xl font-black text-indigo-600 font-mono">
              {{ formatPrice(product.price) }}
            </p>
            <span
              class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
              :class="
                product.stock > 10
                  ? 'bg-green-100 text-green-800'
                  : product.stock > 0
                    ? 'bg-orange-100 text-orange-800'
                    : 'bg-red-100 text-red-800'
              "
            >
              库存: {{ product.stock }} 件
            </span>
          </div>

          <div class="prose prose-sm sm:prose-base text-gray-600 mb-8 flex-grow">
            <h3 class="text-lg font-semibold text-gray-900 mb-2">商品描述</h3>
            <p class="leading-relaxed">{{ product.description || '暂无描述' }}</p>
          </div>

          <div class="mt-auto space-y-6">
            <div class="border-t border-gray-100 pt-6 grid grid-cols-2 gap-4 text-xs text-gray-500">
              <div>
                <span class="block text-gray-400 mb-1">商品编号</span>
                <span class="font-mono text-gray-700">{{ product.id }}</span>
              </div>
              <div v-if="product.createdAt">
                <span class="block text-gray-400 mb-1">上架时间</span>
                <span>{{ formatDate(product.createdAt) }}</span>
              </div>
            </div>

            <Button
              @click="handleBuyNow"
              :disabled="product.stock === 0"
              class="w-full py-4 text-lg shadow-lg hover:shadow-xl transition-all hover:-translate-y-1"
            >
              <svg
                v-if="product.stock > 0"
                class="w-5 h-5 mr-2"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z"
                ></path>
              </svg>
              <svg
                v-else
                class="w-5 h-5 mr-2"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"
                ></path>
              </svg>
              {{ product.stock === 0 ? '已售罄' : '立即购买' }}
            </Button>
          </div>
        </div>
      </div>
    </div>

    <!-- 订单创建对话框 -->
    <Modal :is-open="showOrderModal" title="创建订单" @close="showOrderModal = false">
      <OrderForm
        v-if="product"
        :product-id="product.id"
        :product-name="product.name"
        :product-price="product.price"
        :product-stock="product.stock"
        :is-submitting="isCreating"
        @submit="handleOrderSubmit"
        @cancel="showOrderModal = false"
      />
    </Modal>

    <!-- 删除确认对话框 -->
    <Modal :is-open="showDeleteModal" title="确认删除" @close="showDeleteModal = false">
      <template #default>
        <p class="text-sm text-gray-600">
          确定要删除商品
          <span class="font-medium text-gray-900">{{ product?.name }}</span> 吗？此操作无法撤销。
        </p>
      </template>
      <template #footer>
        <div class="flex justify-end space-x-3">
          <Button variant="secondary" @click="showDeleteModal = false" :disabled="isDeleting">
            取消
          </Button>
          <Button @click="handleDelete" :loading="isDeleting" variant="danger"> 确认删除 </Button>
        </div>
      </template>
    </Modal>
  </MainLayout>
</template>
