<script setup lang="ts">
import { computed } from 'vue';

interface Props {
  currentPage: number;
  totalPages: number;
  maxVisiblePages?: number;
}

const props = withDefaults(defineProps<Props>(), {
  maxVisiblePages: 5,
});

const emit = defineEmits<{
  'update:currentPage': [page: number];
  'page-change': [page: number];
}>();

const hasPrevious = computed(() => props.currentPage > 1);
const hasNext = computed(() => props.currentPage < props.totalPages);

// Calculate visible page numbers
const visiblePages = computed(() => {
  const pages: number[] = [];
  const half = Math.floor(props.maxVisiblePages / 2);

  let start = Math.max(1, props.currentPage - half);
  let end = Math.min(props.totalPages, start + props.maxVisiblePages - 1);

  // Adjust start if we're near the end
  if (end - start + 1 < props.maxVisiblePages) {
    start = Math.max(1, end - props.maxVisiblePages + 1);
  }

  for (let i = start; i <= end; i++) {
    pages.push(i);
  }

  return pages;
});

const showFirstPage = computed(() => visiblePages.value[0] > 1);
const showLastPage = computed(
  () => visiblePages.value[visiblePages.value.length - 1] < props.totalPages
);

const goToPage = (page: number) => {
  if (page >= 1 && page <= props.totalPages && page !== props.currentPage) {
    emit('update:currentPage', page);
    emit('page-change', page);
  }
};

const goToPrevious = () => {
  if (hasPrevious.value) {
    goToPage(props.currentPage - 1);
  }
};

const goToNext = () => {
  if (hasNext.value) {
    goToPage(props.currentPage + 1);
  }
};

// Keyboard navigation
const handleKeyDown = (event: KeyboardEvent, page: number) => {
  if (event.key === 'Enter' || event.key === ' ') {
    event.preventDefault();
    goToPage(page);
  }
};

const handlePreviousKeyDown = (event: KeyboardEvent) => {
  if (event.key === 'Enter' || event.key === ' ') {
    event.preventDefault();
    goToPrevious();
  }
};

const handleNextKeyDown = (event: KeyboardEvent) => {
  if (event.key === 'Enter' || event.key === ' ') {
    event.preventDefault();
    goToNext();
  }
};
</script>

<template>
  <nav class="flex items-center justify-center space-x-2" role="navigation" aria-label="分页导航">
    <!-- Previous button -->
    <button
      type="button"
      :disabled="!hasPrevious"
      :aria-label="hasPrevious ? '上一页' : '已是第一页'"
      class="px-3 py-2 text-sm font-medium rounded-lg transition-colors focus:outline-none focus:ring-2 focus:ring-indigo-500 shadow-sm"
      :class="{
        'text-gray-700 bg-white border border-gray-200 hover:bg-gray-50 hover:text-indigo-600':
          hasPrevious,
        'text-gray-400 bg-gray-50 border border-gray-100 cursor-not-allowed': !hasPrevious,
      }"
      @click="goToPrevious"
      @keydown="handlePreviousKeyDown"
    >
      <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
      </svg>
      <span class="sr-only">上一页</span>
    </button>

    <!-- First page -->
    <button
      v-if="showFirstPage"
      type="button"
      aria-label="第 1 页"
      class="px-3 py-2 text-sm font-medium rounded-lg transition-colors focus:outline-none focus:ring-2 focus:ring-indigo-500 text-gray-700 bg-white border border-gray-200 hover:bg-gray-50 shadow-sm"
      @click="goToPage(1)"
      @keydown="handleKeyDown($event, 1)"
    >
      1
    </button>

    <!-- Ellipsis before -->
    <span v-if="showFirstPage && visiblePages[0] > 2" class="px-2 text-gray-500"> ... </span>

    <!-- Page numbers -->
    <button
      v-for="page in visiblePages"
      :key="page"
      type="button"
      :aria-label="`第 ${page} 页`"
      :aria-current="page === currentPage ? 'page' : undefined"
      class="px-3.5 py-2 text-sm font-medium rounded-lg transition-colors focus:outline-none focus:ring-2 focus:ring-indigo-500 shadow-sm"
      :class="{
        'text-white bg-indigo-600 border border-indigo-600 shadow-md': page === currentPage,
        'text-gray-700 bg-white border border-gray-200 hover:bg-gray-50 hover:text-indigo-600':
          page !== currentPage,
      }"
      @click="goToPage(page)"
      @keydown="handleKeyDown($event, page)"
    >
      {{ page }}
    </button>

    <!-- Ellipsis after -->
    <span
      v-if="showLastPage && visiblePages[visiblePages.length - 1] < totalPages - 1"
      class="px-2 text-gray-500"
    >
      ...
    </span>

    <!-- Last page -->
    <button
      v-if="showLastPage"
      type="button"
      :aria-label="`第 ${totalPages} 页`"
      class="px-3 py-2 text-sm font-medium rounded-lg transition-colors focus:outline-none focus:ring-2 focus:ring-indigo-500 text-gray-700 bg-white border border-gray-200 hover:bg-gray-50 shadow-sm"
      @click="goToPage(totalPages)"
      @keydown="handleKeyDown($event, totalPages)"
    >
      {{ totalPages }}
    </button>

    <!-- Next button -->
    <button
      type="button"
      :disabled="!hasNext"
      :aria-label="hasNext ? '下一页' : '已是最后一页'"
      class="px-3 py-2 text-sm font-medium rounded-lg transition-colors focus:outline-none focus:ring-2 focus:ring-indigo-500 shadow-sm"
      :class="{
        'text-gray-700 bg-white border border-gray-200 hover:bg-gray-50 hover:text-indigo-600':
          hasNext,
        'text-gray-400 bg-gray-50 border border-gray-100 cursor-not-allowed': !hasNext,
      }"
      @click="goToNext"
      @keydown="handleNextKeyDown"
    >
      <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
      </svg>
      <span class="sr-only">下一页</span>
    </button>
  </nav>
</template>
