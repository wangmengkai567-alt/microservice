<script setup lang="ts">
import { computed } from 'vue';

interface Props {
  modelValue: string | number;
  type?: 'text' | 'password' | 'email' | 'number';
  label?: string;
  placeholder?: string;
  error?: string;
  disabled?: boolean;
  required?: boolean;
  id?: string;
}

const props = withDefaults(defineProps<Props>(), {
  type: 'text',
  disabled: false,
  required: false,
});

const emit = defineEmits<{
  'update:modelValue': [value: string | number];
  blur: [event: FocusEvent];
}>();

const inputId = computed(() => props.id || `input-${Math.random().toString(36).substr(2, 9)}`);

const hasError = computed(() => !!props.error);

const handleInput = (event: Event) => {
  const target = event.target as HTMLInputElement;
  const value = props.type === 'number' ? Number(target.value) : target.value;
  emit('update:modelValue', value);
};

const handleBlur = (event: FocusEvent) => {
  emit('blur', event);
};
</script>

<template>
  <div class="w-full relative">
    <label
      v-if="label"
      :for="inputId"
      class="block text-sm font-medium text-gray-700 mb-1.5 transition-colors"
      :class="{ 'text-red-500': hasError }"
    >
      {{ label }}
      <span v-if="required" class="text-red-400 ml-0.5" aria-label="必填">*</span>
    </label>
    <div class="relative">
      <input
        :id="inputId"
        :type="type"
        :value="modelValue"
        :placeholder="placeholder"
        :disabled="disabled"
        :required="required"
        :aria-invalid="hasError"
        :aria-describedby="hasError ? `${inputId}-error` : undefined"
        :aria-required="required"
        @input="handleInput"
        @blur="handleBlur"
        class="block w-full px-4 py-3 border rounded-xl shadow-sm focus:outline-none focus:ring-2 focus:border-transparent transition-all duration-200 bg-white/50 focus:bg-white"
        :class="{
          'border-red-300 text-red-900 placeholder-red-300 focus:ring-red-500 bg-red-50/30':
            hasError,
          'border-gray-200 focus:ring-indigo-500 hover:border-gray-300': !hasError && !disabled,
          'bg-gray-50 text-gray-500 border-gray-200 cursor-not-allowed': disabled,
        }"
      />
      <div
        v-if="hasError"
        class="absolute inset-y-0 right-0 pr-3 flex items-center pointer-events-none"
      >
        <svg class="h-5 w-5 text-red-500" fill="currentColor" viewBox="0 0 20 20">
          <path
            fill-rule="evenodd"
            d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z"
            clip-rule="evenodd"
          />
        </svg>
      </div>
    </div>
    <p
      v-if="hasError"
      :id="`${inputId}-error`"
      class="mt-1.5 text-sm text-red-500 flex items-center"
      role="alert"
    >
      {{ error }}
    </p>
  </div>
</template>
