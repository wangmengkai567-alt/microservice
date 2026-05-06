<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import MainLayout from '@/components/layout/MainLayout.vue';
import OrderItem from '@/components/features/OrderItem.vue';
import Pagination from '@/components/common/Pagination.vue';
import Button from '@/components/common/Button.vue';
import Loading from '@/components/common/Loading.vue';
import Modal from '@/components/common/Modal.vue';
import { useOrders } from '@/composables/useOrders';
import type { Order } from '@/types/order';

const router = useRouter();
const currentPage = ref(1);
const pageSize = 10;

const {
  orders,
  total,
  isLoading,
  error,
  refetch,
  cancelOrder,
  payOrder,
  deleteOrder,
  isCancelling,
  isPaying,
  isDeleting,
} = useOrders(currentPage.value, pageSize);

const showModal = ref(false);
const modalType = ref<'pay' | 'cancel' | 'delete'>('pay');
const targetOrder = ref<Order | null>(null);

function handlePayOrder(order: Order) {
  modalType.value = 'pay';
  targetOrder.value = order;
  showModal.value = true;
}

function handleCancelOrder(order: Order) {
  modalType.value = 'cancel';
  targetOrder.value = order;
  showModal.value = true;
}

function handleDeleteOrder(order: Order) {
  modalType.value = 'delete';
  targetOrder.value = order;
  showModal.value = true;
}

async function confirmAction() {
  if (!targetOrder.value) return;
  const orderId = targetOrder.value.id;

  try {
    if (modalType.value === 'pay') {
      await payOrder(orderId);
    } else if (modalType.value === 'cancel') {
      await cancelOrder(orderId);
    } else if (modalType.value === 'delete') {
      await deleteOrder(orderId);
    }
  } catch (err) {
    console.error(`Failed to ${modalType.value} order:`, err);
  } finally {
    showModal.value = false;
    targetOrder.value = null;
  }
}

function cancelModal() {
  showModal.value = false;
  targetOrder.value = null;
}

function handlePageChange(page: number) {
  currentPage.value = page;
}

// 倒计时结束，自动取消订单（静默，不弹确认框）
async function handleOrderExpired(order: Order) {
  try {
    await cancelOrder(order.id);
    // 取消成功后刷新列表
    refetch();
  } catch (err) {
    console.error('Failed to auto-cancel expired order:', err);
    // 即使取消失败也刷新，可能后端已经自动取消了
    refetch();
  }
}
</script>

<template>
  <MainLayout>
    <div
      class="mb-8 flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 bg-white p-6 rounded-2xl shadow-sm border border-gray-100"
    >
      <div>
        <h1
          class="text-2xl sm:text-3xl font-extrabold text-gray-900 tracking-tight flex items-center"
        >
          <svg
            class="w-8 h-8 mr-3 text-indigo-500"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01"
            ></path>
          </svg>
          我的订单
        </h1>
        <p class="mt-2 text-sm text-gray-500">查看和管理您的所有订单</p>
      </div>
      <Button
        @click="router.push('/products')"
        class="shadow-md hover:shadow-lg transition-all hover:-translate-y-0.5"
      >
        <svg
          class="w-5 h-5 mr-1.5"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
          xmlns="http://www.w3.org/2000/svg"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z"
          ></path>
        </svg>
        继续购物
      </Button>
    </div>

    <!-- 加载状态 -->
    <div v-if="isLoading" class="flex justify-center py-12">
      <Loading />
    </div>

    <!-- 错误提示 -->
    <div v-else-if="error" class="bg-red-50 border border-red-200 rounded-md p-4">
      <p class="text-sm text-red-800">{{ error }}</p>
    </div>

    <!-- 空列表 -->
    <div
      v-else-if="orders.length === 0"
      class="text-center py-20 bg-white rounded-2xl shadow-sm border border-gray-100 flex flex-col items-center justify-center"
    >
      <div class="w-24 h-24 bg-gray-50 rounded-full flex items-center justify-center mb-6">
        <svg class="h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="1.5"
            d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
          ></path>
        </svg>
      </div>
      <h3 class="text-lg font-medium text-gray-900 mb-2">暂无订单</h3>
      <p class="text-gray-500 max-w-sm mx-auto mb-6">
        您还没有任何订单。去看看有没有喜欢的商品吧！
      </p>
      <Button @click="router.push('/products')" class="shadow-sm hover:shadow transition-all">
        <svg
          class="w-5 h-5 mr-1.5"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
          xmlns="http://www.w3.org/2000/svg"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z"
          ></path>
        </svg>
        去购物
      </Button>
    </div>

    <!-- 订单列表 -->
    <div v-else>
      <div class="space-y-4 mb-8">
        <OrderItem
          v-for="order in orders"
          :key="order.id"
          :order="order"
          :is-cancelling="isCancelling"
          :is-paying="isPaying"
          :is-deleting="isDeleting"
          @cancel="handleCancelOrder"
          @pay="handlePayOrder"
          @delete="handleDeleteOrder"
          @expired="handleOrderExpired"
        />
      </div>

      <!-- 分页 -->
      <Pagination
        v-if="total > pageSize"
        :current-page="currentPage"
        :total-pages="Math.ceil(total / pageSize)"
        @page-change="handlePageChange"
      />
    </div>

    <!-- 操作确认对话框 -->
    <Modal
      :is-open="showModal"
      :title="modalType === 'pay' ? '确认支付' : modalType === 'cancel' ? '确认取消' : '确认删除'"
      @close="cancelModal"
    >
      <template #default>
        <p class="text-sm text-gray-600">
          <template v-if="modalType === 'pay'">
            确定要支付订单
            <span class="font-medium text-gray-900">{{ targetOrder?.orderNo }}</span>
            吗？金额：<span class="font-medium text-indigo-600 font-mono"
              >¥{{ targetOrder?.totalPrice.toFixed(2) }}</span
            >
          </template>
          <template v-else-if="modalType === 'cancel'">
            确定要取消订单
            <span class="font-medium text-gray-900">{{ targetOrder?.orderNo }}</span> 吗？
          </template>
          <template v-else-if="modalType === 'delete'">
            确定要删除订单
            <span class="font-medium text-gray-900">{{ targetOrder?.orderNo }}</span>
            吗？删除后无法恢复。
          </template>
        </p>
      </template>
      <template #footer>
        <div class="flex justify-end space-x-3">
          <Button
            variant="secondary"
            @click="cancelModal"
            :disabled="isPaying || isCancelling || isDeleting"
          >
            取消
          </Button>
          <Button
            @click="confirmAction"
            :loading="isPaying || isCancelling || isDeleting"
            :variant="modalType === 'pay' ? 'primary' : 'danger'"
          >
            {{
              modalType === 'pay' ? '确认支付' : modalType === 'cancel' ? '确认取消' : '确认删除'
            }}
          </Button>
        </div>
      </template>
    </Modal>
  </MainLayout>
</template>
