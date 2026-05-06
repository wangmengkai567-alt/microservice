<script setup lang="ts">
import { computed } from 'vue';

interface Props {
  type?: 'button' | 'submit' | 'reset';
  variant?: 'primary' | 'secondary' | 'danger' | 'ghost';
  size?: 'sm' | 'md' | 'lg';
  loading?: boolean;
  disabled?: boolean;
  fullWidth?: boolean;
  ariaLabel?: string;
}

const props = withDefaults(defineProps<Props>(), {
  type: 'button',
  variant: 'primary',
  size: 'md',
  loading: false,
  disabled: false,
  fullWidth: false,
});

const emit = defineEmits<{
  click: [event: MouseEvent];
}>();

const isDisabled = computed(() => props.disabled || props.loading);

const buttonClasses = computed(() => {
  const classes = [
    'inline-flex items-center justify-center font-medium rounded-xl',
    'focus:outline-none focus:ring-2 focus:ring-offset-2',
    'transition-all duration-200',
    'disabled:opacity-50 disabled:cursor-not-allowed',
  ];

  // Size classes
  if (props.size === 'sm') {
    classes.push('px-3 py-1.5 text-sm');
  } else if (props.size === 'lg') {
    classes.push('px-6 py-3 text-lg');
  } else {
    classes.push('px-4 py-2.5 text-sm');
  }

  // Variant classes
  if (props.variant === 'primary') {
    classes.push(
      'bg-indigo-600 text-white hover:bg-indigo-700 hover:shadow-md hover:-translate-y-0.5',
      'focus:ring-indigo-500',
      'disabled:hover:bg-indigo-600 disabled:hover:translate-y-0 disabled:hover:shadow-none'
    );
  } else if (props.variant === 'secondary') {
    classes.push(
      'bg-white border border-gray-200 text-gray-700 hover:bg-gray-50 hover:text-gray-900 shadow-sm',
      'focus:ring-indigo-500',
      'disabled:hover:bg-white'
    );
  } else if (props.variant === 'danger') {
    classes.push(
      'bg-red-600 text-white hover:bg-red-700 hover:shadow-md hover:-translate-y-0.5',
      'focus:ring-red-500',
      'disabled:hover:bg-red-600 disabled:hover:translate-y-0 disabled:hover:shadow-none'
    );
  } else if (props.variant === 'ghost') {
    classes.push(
      'bg-transparent text-gray-700 hover:bg-gray-100',
      'focus:ring-gray-500',
      'disabled:hover:bg-transparent'
    );
  }

  // Full width
  if (props.fullWidth) {
    classes.push('w-full');
  }

  return classes.join(' ');
});

const handleClick = (event: MouseEvent) => {
  if (!isDisabled.value) {
    emit('click', event);
  }
};
</script>

<template>
  <button
    :type="type"
    :disabled="isDisabled"
    :aria-label="ariaLabel"
    :aria-busy="loading"
    :class="buttonClasses"
    @click="handleClick"
  >
    <!-- Loading spinner -->
    <svg
      v-if="loading"
      class="animate-spin -ml-1 mr-2 h-4 w-4"
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
    <span v-if="loading" class="sr-only">加载中...</span>
    <slot />
  </button>
</template>
