<script setup lang="ts">
import { computed } from 'vue';

export interface SelectOption {
  value: string | number;
  label: string;
  disabled?: boolean;
}

interface Props {
  modelValue: string | number;
  options: SelectOption[];
  label?: string;
  placeholder?: string;
  error?: string;
  disabled?: boolean;
  required?: boolean;
  id?: string;
}

const props = withDefaults(defineProps<Props>(), {
  disabled: false,
  required: false,
  placeholder: '请选择',
});

const emit = defineEmits<{
  'update:modelValue': [value: string | number];
  change: [value: string | number];
}>();

const selectId = computed(() => props.id || `select-${Math.random().toString(36).substr(2, 9)}`);

const hasError = computed(() => !!props.error);

const handleChange = (event: Event) => {
  const target = event.target as HTMLSelectElement;
  const value = target.value;
  emit('update:modelValue', value);
  emit('change', value);
};
</script>

<template>
  <div class="w-full">
    <label v-if="label" :for="selectId" class="block text-sm font-medium text-gray-700 mb-1">
      {{ label }}
      <span v-if="required" class="text-red-500" aria-label="必填">*</span>
    </label>
    <select
      :id="selectId"
      :value="modelValue"
      :disabled="disabled"
      :required="required"
      :aria-invalid="hasError"
      :aria-describedby="hasError ? `${selectId}-error` : undefined"
      :aria-required="required"
      @change="handleChange"
      class="block w-full px-3 py-2 border rounded-md shadow-sm focus:outline-none focus:ring-2 transition-colors"
      :class="{
        'border-red-300 text-red-900 focus:ring-red-500 focus:border-red-500': hasError,
        'border-gray-300 focus:ring-blue-500 focus:border-blue-500': !hasError && !disabled,
        'bg-gray-100 cursor-not-allowed': disabled,
      }"
    >
      <option value="" disabled>{{ placeholder }}</option>
      <option
        v-for="option in options"
        :key="option.value"
        :value="option.value"
        :disabled="option.disabled"
      >
        {{ option.label }}
      </option>
    </select>
    <p v-if="hasError" :id="`${selectId}-error`" class="mt-1 text-sm text-red-600" role="alert">
      {{ error }}
    </p>
  </div>
</template>
