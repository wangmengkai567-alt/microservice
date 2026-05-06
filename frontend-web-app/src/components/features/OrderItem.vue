<template>
  <div
    class="bg-white rounded-xl shadow-sm border border-gray-100 hover:shadow-lg hover:-translate-y-1 transition-all duration-300 p-6 overflow-hidden relative"
  >
    <div v-if="order.status === 'pending'" class="absolute top-0 left-0 w-1 h-full bg-yellow-400"></div>
    <div v-else-if="order.status === 'paid'" class="absolute top-0 left-0 w-1 h-full bg-blue-500"></div>
    <div v-else-if="order.status === 'completed'" class="absolute top-0 left-0 w-1 h-full bg-green-500"></div>
    <div v-else-if="order.status === 'cancelled'" class="absolute top-0 left-0 w-1 h-full bg-gray-300"></div>

    <div
      class="flex flex-col md:flex-row md:items-center md:justify-between space-y-5 md:space-y-0 pl-2"
    >
      <div class="flex-1">
        <div class="flex flex-wrap items-center gap-3 mb-3">
          <h3 class="text-lg font-bold text-gray-900 tracking-tight">订单号: <span class="font-mono text-indigo-600">{{ order.orderNo }}</span></h3>
          <span :class="statusBadgeClass">
            {{ statusText }}
          </span>
          <span
            v-if="order.status === 'pending' && remainingSeconds > 0"
            class="text-xs text-red-500 font-mono bg-red-50 px-2 py-1 rounded-md border border-red-100 flex items-center"
          >
            <svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
            {{ countdownText }} 后自动取消
          </span>
        </div>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-y-2 gap-x-6 text-sm text-gray-600 mt-4 bg-gray-50 p-4 rounded-lg">
          <p class="flex items-center"><span class="text-gray-400 w-16">商品:</span> <span class="font-medium text-gray-900">{{ order.productName }}</span></p>
          <p class="flex items-center"><span class="text-gray-400 w-16">数量:</span> <span class="font-medium text-gray-900">{{ order.quantity }} 件</span></p>
          <p class="flex items-center"><span class="text-gray-400 w-16">创建时间:</span> <span>{{ formatDate(order.createdAt) }}</span></p>
          <p class="flex items-center"><span class="text-gray-400 w-16">总价:</span> <span class="text-lg font-bold text-indigo-600 font-mono">{{ formatPrice(order.totalPrice) }}</span></p>
        </div>
      </div>
      <div class="flex flex-wrap gap-2 md:pl-6 md:border-l md:border-gray-100 md:ml-6 min-w-[140px] justify-end md:justify-center flex-col sm:flex-row md:flex-col">
        <router-link
          :to="`/orders/${order.id}`"
          class="inline-flex justify-center items-center px-4 py-2.5 border border-indigo-100 text-sm font-medium rounded-xl text-indigo-700 bg-indigo-50 hover:bg-indigo-100 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition-colors w-full sm:w-auto md:w-full"
        >
          查看详情
        </router-link>
        <button
          v-if="order.status === 'pending'"
          @click="$emit('pay', order)"
          class="inline-flex justify-center items-center px-4 py-2.5 border border-transparent text-sm font-medium rounded-xl text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500 transition-all hover:shadow-md hover:-translate-y-0.5 w-full sm:w-auto md:w-full"
          :disabled="isPaying"
        >
          {{ isPaying ? '支付中...' : '立即支付' }}
        </button>
        <button
          v-if="order.status === 'pending'"
          @click="$emit('cancel', order)"
          class="inline-flex justify-center items-center px-4 py-2.5 border border-red-200 text-sm font-medium rounded-xl text-red-600 bg-white hover:bg-red-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 transition-colors w-full sm:w-auto md:w-full"
          :disabled="isCancelling"
        >
          {{ isCancelling ? '取消中...' : '取消订单' }}
        </button>
        <button
          v-if="order.status === 'cancelled'"
          @click="$emit('delete', order)"
          class="inline-flex justify-center items-center px-4 py-2.5 border border-gray-200 text-sm font-medium rounded-xl text-gray-600 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-500 transition-colors w-full sm:w-auto md:w-full"
          :disabled="isDeleting"
        >
          <svg class="w-4 h-4 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg>
          {{ isDeleting ? '删除中...' : '删除订单' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted, onUnmounted } from 'vue';
import type { Order } from '@/types/order';
import { formatPrice, formatDate } from '@/utils/format';

interface Props {
  order: Order;
  isCancelling?: boolean;
  isPaying?: boolean;
  isDeleting?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  isCancelling: false,
  isPaying: false,
  isDeleting: false,
});

const emit = defineEmits<{
  cancel: [order: Order];
  pay: [order: Order];
  delete: [order: Order];
  expired: [order: Order];
}>();

// 倒计时：基于订单创建时间计算剩余秒数（1分钟超时）
const TIMEOUT_SECONDS = 60;
// -1 表示尚未初始化，避免挂载前误显示"等待取消中..."
const remainingSeconds = ref(-1);
let timer: ReturnType<typeof setInterval> | null = null;

function calcRemaining() {
  const raw = props.order.createdAt;
  // 后端返回 "YYYY-MM-DD HH:mm:ss" 无时区，视为 UTC 解析
  const created = new Date(raw.replace(' ', 'T') + 'Z').getTime();
  if (isNaN(created)) return 0;
  const elapsed = Math.floor((Date.now() - created) / 1000);
  return Math.max(0, TIMEOUT_SECONDS - elapsed);
}

function startCountdown() {
  if (props.order.status !== 'pending') return;
  remainingSeconds.value = calcRemaining();
  if (remainingSeconds.value === 0) {
    emit('expired', props.order);
    return;
  }
  timer = setInterval(() => {
    remainingSeconds.value = calcRemaining();
    if (remainingSeconds.value === 0) {
      clearInterval(timer!);
      emit('expired', props.order);
    }
  }, 1000);
}

onMounted(startCountdown);
onUnmounted(() => {
  if (timer) clearInterval(timer);
});

const countdownText = computed(() => {
  const s = remainingSeconds.value;
  return `${String(Math.floor(s / 60)).padStart(2, '0')}:${String(s % 60).padStart(2, '0')}`;
});

const statusText = computed(() => {
  const statusMap: Record<string, string> = {
    pending: '待支付',
    paid: '已支付',
    completed: '已完成',
    cancelled: '已取消',
  };
  return statusMap[props.order.status] || props.order.status;
});

const statusBadgeClass = computed(() => {
  const baseClass = 'px-2 py-1 text-xs font-medium rounded-full';
  const statusClasses: Record<string, string> = {
    pending: 'bg-yellow-100 text-yellow-800',
    paid: 'bg-blue-100 text-blue-800',
    completed: 'bg-green-100 text-green-800',
    cancelled: 'bg-gray-100 text-gray-800',
  };
  return `${baseClass} ${statusClasses[props.order.status] || 'bg-gray-100 text-gray-800'}`;
});
</script>
