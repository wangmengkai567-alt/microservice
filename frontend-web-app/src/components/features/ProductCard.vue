<template>
  <div
    class="group bg-white rounded-2xl shadow-sm border border-gray-100 hover:shadow-xl hover:-translate-y-1 transition-all duration-300 overflow-hidden flex flex-col h-full"
  >
    <!-- 商品图片 -->
    <div class="h-48 bg-gray-100 relative overflow-hidden">
      <img
        v-if="product.image_url || product.imageUrl"
        :src="product.image_url || product.imageUrl"
        :alt="product.name"
        class="w-full h-full object-cover"
      />
      <div v-else class="absolute inset-0 flex items-center justify-center text-gray-400">
        <svg class="w-12 h-12" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="1.5"
            d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"
          />
        </svg>
      </div>
    </div>

    <div class="p-5 flex flex-col flex-grow">
      <h3
        class="text-lg font-bold text-gray-900 mb-2 group-hover:text-indigo-600 transition-colors"
      >
        {{ product.name }}
      </h3>
      <p class="text-sm text-gray-500 mb-4 line-clamp-2 flex-grow">
        {{ product.description || '暂无商品描述' }}
      </p>

      <div class="flex items-end justify-between mb-5 pt-4 border-t border-gray-50">
        <div>
          <span class="text-xs text-gray-400 font-medium">价格</span>
          <div class="text-2xl font-bold text-indigo-600 flex items-baseline">
            {{ formatPrice(product.price) }}
          </div>
        </div>
        <div class="text-right">
          <span class="text-xs text-gray-400 font-medium block mb-1">库存状态</span>
          <span
            class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
            :class="stockBadgeClass"
          >
            {{ product.stock }} 件
          </span>
        </div>
      </div>

      <div class="flex space-x-2 mt-auto">
        <router-link
          :to="`/products/${product.id}`"
          class="flex-1 inline-flex justify-center items-center px-4 py-2 border border-indigo-100 text-sm font-medium rounded-xl text-indigo-700 bg-indigo-50 hover:bg-indigo-100 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition-colors"
        >
          查看详情
        </router-link>
        <button
          v-if="showActions"
          @click="$emit('edit', product)"
          class="px-3 py-2 border border-gray-200 text-sm font-medium rounded-xl text-gray-600 bg-white hover:bg-gray-50 hover:text-gray-900 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition-colors shadow-sm"
          aria-label="编辑商品"
          title="编辑"
        >
          <svg
            class="w-4 h-4"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"
            ></path>
          </svg>
        </button>
        <button
          v-if="showActions"
          @click="$emit('delete', product)"
          class="px-3 py-2 border border-red-100 text-sm font-medium rounded-xl text-red-600 bg-red-50 hover:bg-red-100 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 transition-colors shadow-sm"
          aria-label="删除商品"
          title="删除"
        >
          <svg
            class="w-4 h-4"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
            xmlns="http://www.w3.org/2000/svg"
          >
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
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import type { Product } from '@/types/product';
import { formatPrice } from '@/utils/format';

interface Props {
  product: Product;
  showActions?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  showActions: false,
});

defineEmits<{
  edit: [product: Product];
  delete: [product: Product];
}>();

const stockBadgeClass = computed(() => {
  if (props.product.stock === 0) {
    return 'bg-red-100 text-red-800';
  } else if (props.product.stock < 10) {
    return 'bg-orange-100 text-orange-800';
  }
  return 'bg-green-100 text-green-800';
});
</script>
