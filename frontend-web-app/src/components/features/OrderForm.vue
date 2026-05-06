<template>
  <form @submit.prevent="handleSubmit" class="space-y-6">
    <div>
      <label class="block text-sm font-medium text-gray-700 mb-1"> 商品信息 </label>
      <div class="bg-gray-50 p-4 rounded-md border border-gray-200">
        <p class="text-sm text-gray-900 font-medium">{{ productName }}</p>
        <p class="text-sm text-gray-600 mt-1">单价: ¥{{ formatPrice(productPrice) }}</p>
        <p class="text-sm text-gray-600">可用库存: {{ productStock }}</p>
      </div>
    </div>

    <Input
      v-model.number="formData.quantity"
      label="购买数量"
      type="number"
      min="1"
      :max="productStock"
      placeholder="请输入购买数量"
      :error="errors.quantity"
      required
      aria-required="true"
    />

    <div v-if="totalPrice > 0" class="bg-indigo-50 p-4 rounded-md border border-indigo-200">
      <p class="text-sm text-gray-700">
        总价: <span class="text-xl font-bold text-indigo-600">{{ formatPrice(totalPrice) }}</span>
      </p>
    </div>

    <div class="flex justify-end space-x-3 pt-4">
      <Button type="button" variant="secondary" @click="$emit('cancel')" :disabled="isSubmitting">
        取消
      </Button>
      <Button type="submit" :loading="isSubmitting" :disabled="isSubmitting || !canSubmit">
        {{ submitText }}
      </Button>
    </div>
  </form>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import Input from '@/components/common/Input.vue';
import Button from '@/components/common/Button.vue';
import { orderSchema } from '@/utils/validation';
import type { OrderFormData } from '@/types/order';
import { formatPrice } from '@/utils/format';
import { z } from 'zod';

interface Props {
  productId: string | number;
  productName: string;
  productPrice: number;
  productStock: number;
  isSubmitting?: boolean;
  submitText?: string;
}

const props = withDefaults(defineProps<Props>(), {
  isSubmitting: false,
  submitText: '确认购买',
});

const emit = defineEmits<{
  submit: [data: OrderFormData];
  cancel: [];
}>();

const formData = ref<OrderFormData>({
  productId: props.productId,
  quantity: 1,
});

const errors = ref<Record<string, string>>({});

const totalPrice = computed(() => {
  return props.productPrice * formData.value.quantity;
});

const canSubmit = computed(() => {
  return (
    formData.value.quantity > 0 &&
    formData.value.quantity <= props.productStock &&
    !props.isSubmitting
  );
});

function handleSubmit() {
  // Clear previous errors
  errors.value = {};

  // Validate form data
  try {
    // Validate quantity is within stock
    if (formData.value.quantity > props.productStock) {
      errors.value.quantity = `数量不能超过库存 (${props.productStock})`;
      return;
    }

    orderSchema.parse(formData.value);
    emit('submit', formData.value);
  } catch (error) {
    if (error instanceof z.ZodError) {
      error.errors.forEach((err) => {
        if (err.path[0]) {
          errors.value[err.path[0] as string] = err.message;
        }
      });
    }
  }
}
</script>
