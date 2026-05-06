<script setup lang="ts">
interface Props {
  size?: 'sm' | 'md' | 'lg';
  fullScreen?: boolean;
  text?: string;
}

withDefaults(defineProps<Props>(), {
  size: 'md',
  fullScreen: false,
});

const sizeClasses = {
  sm: 'h-6 w-6',
  md: 'h-10 w-10',
  lg: 'h-16 w-16',
};
</script>

<template>
  <div
    v-if="fullScreen"
    class="fixed inset-0 z-50 flex items-center justify-center bg-white bg-opacity-75"
    role="status"
    aria-live="polite"
    aria-label="加载中"
  >
    <div class="text-center">
      <svg
        :class="['animate-spin mx-auto text-blue-600', sizeClasses[size]]"
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 24 24"
        aria-hidden="true"
      >
        <circle
          class="opacity-25"
          cx="12"
          cy="12"
          r="10"
          stroke="currentColor"
          stroke-width="4"
        ></circle>
        <path
          class="opacity-75"
          fill="currentColor"
          d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
        ></path>
      </svg>
      <p v-if="text" class="mt-4 text-sm text-gray-600">{{ text }}</p>
      <span v-else class="sr-only">加载中...</span>
    </div>
  </div>

  <div
    v-else
    class="flex items-center justify-center"
    role="status"
    aria-live="polite"
    aria-label="加载中"
  >
    <svg
      :class="['animate-spin text-blue-600', sizeClasses[size]]"
      xmlns="http://www.w3.org/2000/svg"
      fill="none"
      viewBox="0 0 24 24"
      aria-hidden="true"
    >
      <circle
        class="opacity-25"
        cx="12"
        cy="12"
        r="10"
        stroke="currentColor"
        stroke-width="4"
      ></circle>
      <path
        class="opacity-75"
        fill="currentColor"
        d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
      ></path>
    </svg>
    <span v-if="text" class="ml-3 text-sm text-gray-600">{{ text }}</span>
    <span v-else class="sr-only">加载中...</span>
  </div>
</template>
