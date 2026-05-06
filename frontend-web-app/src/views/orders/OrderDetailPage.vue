<script setup lang="ts">
import { computed, ref, onMounted, onUnmounted, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import MainLayout from '@/components/layout/MainLayout.vue';
import Button from '@/components/common/Button.vue';
import Loading from '@/components/common/Loading.vue';
import Modal from '@/components/common/Modal.vue';
import { useOrderDetail } from '@/composables/useOrders';
import { formatPrice, formatDate } from '@/utils/format';

const route = useRoute();
const router = useRouter();

const orderId = computed(() => route.params.id as string);
const {
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
  isCancelling,
  isPaying,
  isDeleting,
} = useOrderDetail(orderId.value);

const showModal = ref(false);
const modalType = ref<'cancel' | 'pay' | 'delete' | ''>('');

const displayUpdatedAt = computed(() => order.value?.updatedAt || order.value?.createdAt || '');

// 倒计时
const TIMEOUT_SECONDS = 60;
const remainingSeconds = ref(-1);
let timer: ReturnType<typeof setInterval> | null = null;

function calcRemaining() {
  if (!order.value) return 0;
  const raw = order.value.createdAt;
  const created = new Date(raw.replace(' ', 'T') + 'Z').getTime();
  if (isNaN(created)) return 0;
  const elapsed = Math.floor((Date.now() - created) / 1000);
  return Math.max(0, TIMEOUT_SECONDS - elapsed);
}

function startCountdown() {
  if (order.value?.status !== 'pending') return;
  remainingSeconds.value = calcRemaining();
  if (remainingSeconds.value === 0) {
    autoCancel();
    return;
  }
  timer = setInterval(() => {
    remainingSeconds.value = calcRemaining();
    if (remainingSeconds.value === 0) {
      clearInterval(timer!);
      autoCancel();
    }
  }, 1000);
}

// 倒计时结束自动取消（静默）
async function autoCancel() {
  try {
    await cancelOrder();
  } catch (err) {
    console.error('Auto cancel failed:', err);
  } finally {
    refetch();
  }
}

onMounted(() => {
  if (order.value) startCountdown();
});
onUnmounted(() => {
  if (timer) clearInterval(timer);
});

// order 异步加载完成后启动倒计时
watch(
  () => order.value,
  (val) => {
    if (val?.status === 'pending' && timer === null) startCountdown();
  }
);

const countdownText = computed(() => {
  const s = remainingSeconds.value;
  return `${String(Math.floor(s / 60)).padStart(2, '0')}:${String(s % 60).padStart(2, '0')}`;
});

function getStatusText(status: string) {
  const statusMap: Record<string, string> = {
    pending: '待支付',
    paid: '已支付',
    completed: '已完成',
    cancelled: '已取消',
  };
  return statusMap[status] || status;
}

function getStatusClass(status: string) {
  const classMap: Record<string, string> = {
    pending: 'bg-yellow-100 text-yellow-800',
    paid: 'bg-blue-100 text-blue-800',
    completed: 'bg-green-100 text-green-800',
    cancelled: 'bg-gray-100 text-gray-800',
  };
  return classMap[status] || 'bg-gray-100 text-gray-800';
}

function confirmCancel() {
  modalType.value = 'cancel';
  showModal.value = true;
}

function handlePayOrder() {
  if (!order.value) return;
  modalType.value = 'pay';
  showModal.value = true;
}

function handleDeleteOrder() {
  if (!order.value) return;
  modalType.value = 'delete';
  showModal.value = true;
}

async function confirmAction() {
  if (modalType.value === 'cancel') {
    try {
      await cancelOrder();
      showModal.value = false;
    } catch (err) {
      console.error('Failed to cancel order:', err);
    }
  } else if (modalType.value === 'pay') {
    try {
      await payOrder();
      showModal.value = false;
    } catch (err) {
      console.error('Failed to pay order:', err);
    }
  } else if (modalType.value === 'delete') {
    try {
      await deleteOrder();
      showModal.value = false;
      router.push('/orders');
    } catch (err) {
      console.error('Failed to delete order:', err);
    }
  }
}

function cancelModal() {
  showModal.value = false;
}
</script>

<template>
  <MainLayout>
    <!-- 返回按钮 -->
    <button
      @click="router.push('/orders')"
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
      返回订单列表
    </button>

    <!-- 加载状态 -->
    <div v-if="isLoading" class="flex justify-center py-12">
      <Loading />
    </div>

    <!-- 错误提示 -->
    <div v-else-if="error" class="bg-red-50 border border-red-200 rounded-md p-4">
      <p class="text-sm text-red-800">{{ error }}</p>
      <Button @click="router.push('/orders')" variant="secondary" class="mt-4">
        返回订单列表
      </Button>
    </div>

    <!-- 订单详情 -->
    <div
      v-else-if="order"
      class="bg-white shadow-xl rounded-2xl border border-gray-100 overflow-hidden max-w-4xl mx-auto relative"
    >
      <div
        v-if="order.status === 'pending'"
        class="absolute top-0 left-0 w-full h-1 bg-yellow-400"
      ></div>
      <div
        v-else-if="order.status === 'paid'"
        class="absolute top-0 left-0 w-full h-1 bg-blue-500"
      ></div>
      <div
        v-else-if="order.status === 'completed'"
        class="absolute top-0 left-0 w-full h-1 bg-green-500"
      ></div>
      <div
        v-else-if="order.status === 'cancelled'"
        class="absolute top-0 left-0 w-full h-1 bg-gray-300"
      ></div>

      <div class="p-8 sm:p-10">
        <div
          class="flex flex-col sm:flex-row justify-between items-start sm:items-center mb-8 gap-4"
        >
          <div>
            <h1 class="text-3xl font-extrabold text-gray-900 mb-2 tracking-tight">订单详情</h1>
            <p class="text-sm text-gray-500 flex items-center">
              订单号:
              <span
                class="ml-2 font-mono bg-gray-50 px-2 py-0.5 rounded text-gray-700 border border-gray-200"
                >{{ order.orderNo }}</span
              >
            </p>
          </div>
          <div class="flex flex-col items-end space-y-2">
            <span
              :class="[
                'px-3 py-1 text-sm font-medium rounded-full border',
                getStatusClass(order.status),
              ]"
            >
              {{ getStatusText(order.status) }}
            </span>
            <span
              v-if="order.status === 'pending' && remainingSeconds > 0"
              class="text-xs text-red-500 font-mono bg-red-50 px-2 py-1 rounded-md border border-red-100 flex items-center"
            >
              <svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
                ></path>
              </svg>
              {{ countdownText }} 后自动取消
            </span>
          </div>
        </div>

        <div class="bg-gray-50 rounded-xl p-6 sm:p-8 mb-8 border border-gray-100">
          <h3 class="text-lg font-bold text-gray-900 mb-6 border-b border-gray-200 pb-4">
            商品信息
          </h3>
          <dl class="grid grid-cols-1 gap-x-6 gap-y-6 sm:grid-cols-2">
            <div class="sm:col-span-2 flex items-center justify-between">
              <dt class="text-sm font-medium text-gray-500">商品名称</dt>
              <dd class="text-base font-semibold text-gray-900">{{ order.productName }}</dd>
            </div>

            <div
              class="flex items-center justify-between border-t border-gray-200 sm:border-t-0 pt-4 sm:pt-0"
            >
              <dt class="text-sm font-medium text-gray-500">购买数量</dt>
              <dd class="text-sm font-medium text-gray-900">{{ order.quantity }} 件</dd>
            </div>

            <div
              class="flex items-center justify-between border-t border-gray-200 sm:border-t-0 pt-4 sm:pt-0"
            >
              <dt class="text-sm font-medium text-gray-500">订单总价</dt>
              <dd class="text-2xl font-black text-indigo-600 font-mono">
                {{ formatPrice(order.totalPrice) }}
              </dd>
            </div>
          </dl>
        </div>

        <div
          class="grid grid-cols-1 gap-x-4 gap-y-6 sm:grid-cols-2 text-sm text-gray-500 bg-white border border-gray-100 rounded-xl p-6 mb-8 shadow-sm"
        >
          <div class="flex flex-col">
            <dt class="font-medium text-gray-400 mb-1">创建时间</dt>
            <dd class="text-gray-800">{{ formatDate(order.createdAt) }}</dd>
          </div>
          <div class="flex flex-col">
            <dt class="font-medium text-gray-400 mb-1">更新时间</dt>
            <dd class="text-gray-800">{{ formatDate(displayUpdatedAt) }}</dd>
          </div>
        </div>

        <!-- 操作按钮 -->
        <div
          v-if="canPay || canCancel || canDelete"
          class="mt-8 flex flex-col sm:flex-row sm:justify-end gap-3 pt-6 border-t border-gray-100"
        >
          <Button
            v-if="canDelete"
            @click="handleDeleteOrder"
            variant="secondary"
            :loading="isDeleting"
            :disabled="isDeleting"
            class="sm:mr-auto"
          >
            <svg class="w-4 h-4 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
              ></path>
            </svg>
            删除订单
          </Button>
          <Button v-if="canCancel" @click="confirmCancel" variant="danger" :disabled="isCancelling">
            取消订单
          </Button>
          <Button
            v-if="canPay"
            @click="handlePayOrder"
            :loading="isPaying"
            :disabled="isPaying"
            class="shadow-md hover:shadow-lg transition-all hover:-translate-y-0.5"
          >
            立即支付
          </Button>
        </div>
      </div>
    </div>

    <!-- 确认对话框 -->
    <Modal
      :is-open="showModal"
      :title="modalType === 'pay' ? '确认支付' : modalType === 'cancel' ? '确认取消' : '确认删除'"
      @close="cancelModal"
    >
      <template #default>
        <p class="text-sm text-gray-600">
          <template v-if="modalType === 'pay'">
            确定要支付订单
            <span class="font-medium text-gray-900">{{ order?.orderNo }}</span> 吗？金额：<span
              class="font-medium text-indigo-600 font-mono"
              >¥{{ order?.totalPrice.toFixed(2) }}</span
            >
          </template>
          <template v-else-if="modalType === 'cancel'">
            确定要取消订单
            <span class="font-medium text-gray-900">{{ order?.orderNo }}</span> 吗？此操作无法撤销。
          </template>
          <template v-else-if="modalType === 'delete'">
            确定要删除订单
            <span class="font-medium text-gray-900">{{ order?.orderNo }}</span> 吗？删除后无法恢复。
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
