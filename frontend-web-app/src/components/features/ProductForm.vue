<template>
  <form @submit.prevent="handleSubmit" class="space-y-6">
    <Input
      v-model="formData.name"
      label="商品名称"
      type="text"
      placeholder="请输入商品名称"
      :error="errors.name"
      required
      aria-required="true"
    />

    <div>
      <label for="description" class="block text-sm font-medium text-gray-700 mb-1.5">
        商品描述
      </label>
      <textarea
        id="description"
        v-model="formData.description"
        rows="4"
        class="w-full px-4 py-3 border border-gray-200 rounded-xl shadow-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-all duration-200 bg-white/50 focus:bg-white"
        :class="{ 'border-red-300 focus:ring-red-500 bg-red-50/30': errors.description }"
        placeholder="请输入详细的商品描述，例如材质、产地、使用方法等..."
        aria-required="true"
      ></textarea>
      <p v-if="errors.description" class="mt-1.5 text-sm text-red-500 flex items-center">
        <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
          ></path>
        </svg>
        {{ errors.description }}
      </p>
    </div>

    <Input
      v-model.number="formData.price"
      label="价格"
      type="number"
      step="0.01"
      min="0"
      placeholder="请输入价格"
      :error="errors.price"
      required
      aria-required="true"
    />

    <Input
      v-model.number="formData.stock"
      label="库存"
      type="number"
      min="0"
      placeholder="请输入商品库存数量"
      :error="errors.stock"
      required
      aria-required="true"
    />

    <div>
      <label for="product-image" class="block text-sm font-medium text-gray-700 mb-1.5"> 商品图片 </label>
      <input
        id="product-image"
        type="file"
        accept="image/png,image/jpeg,image/jpg,image/webp,image/gif"
        class="block w-full text-sm text-gray-600 file:mr-4 file:py-2 file:px-4 file:rounded-lg file:border-0 file:text-sm file:font-medium file:bg-indigo-50 file:text-indigo-700 hover:file:bg-indigo-100"
        :disabled="isUploadingImage"
        @change="handleFileChange"
      />
      <p class="mt-1 text-xs text-gray-500">支持 JPG/PNG/WEBP/GIF，大小不超过 5MB。</p>
      <p v-if="isUploadingImage" class="mt-1.5 text-sm text-indigo-600">
        正在上传图片... {{ uploadProgress }}%
      </p>
      <p v-if="errors.imageUrl" class="mt-1.5 text-sm text-red-500">{{ errors.imageUrl }}</p>
      <button
        v-if="errors.imageUrl && selectedFile"
        type="button"
        class="mt-2 text-sm text-indigo-600 hover:text-indigo-700"
        @click="retryUploadImage"
      >
        重新上传
      </button>
      <div v-if="previewImageUrl" class="mt-3">
        <img :src="previewImageUrl" alt="预览图" class="h-24 w-24 rounded-lg object-cover border border-gray-200" />
      </div>
    </div>

    <div class="flex justify-end space-x-3 pt-4">
      <Button type="button" variant="secondary" @click="$emit('cancel')" :disabled="isSubmitting">
        取消
      </Button>
      <Button type="submit" :loading="isSubmitting" :disabled="isSubmitting || isUploadingImage">
        {{ submitText }}
      </Button>
    </div>
  </form>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';
import Input from '@/components/common/Input.vue';
import Button from '@/components/common/Button.vue';
import { productApi } from '@/services/api/product';
import { productSchema } from '@/utils/validation';
import type { Product, ProductFormData } from '@/types/product';
import { z } from 'zod';

interface Props {
  initialData?: Product;
  isSubmitting?: boolean;
  submitText?: string;
}

const props = withDefaults(defineProps<Props>(), {
  isSubmitting: false,
  submitText: '提交',
});

const emit = defineEmits<{
  submit: [data: ProductFormData];
  cancel: [];
}>();

const formData = ref<ProductFormData>({
  name: props.initialData?.name || '',
  description: props.initialData?.description || '',
  price: props.initialData?.price || 0,
  stock: props.initialData?.stock || 0,
  imageUrl: props.initialData?.imageUrl || props.initialData?.image_url || '',
});

const errors = ref<Record<string, string>>({});
const previewImageUrl = ref<string>(formData.value.imageUrl || '');
const isUploadingImage = ref(false);
const uploadProgress = ref(0);
const selectedFile = ref<File | null>(null);

// Watch for initialData changes (for edit mode)
watch(
  () => props.initialData,
  (newData) => {
    if (newData) {
      formData.value = {
        name: newData.name,
        description: newData.description,
        price: newData.price,
        stock: newData.stock,
        imageUrl: newData.imageUrl || (newData as any).image_url || '',
      };
      previewImageUrl.value = formData.value.imageUrl || '';
    }
  },
  { immediate: true }
);

async function handleFileChange(event: Event) {
  const input = event.target as HTMLInputElement;
  const file = input.files?.[0];
  if (!file) return;

  selectedFile.value = file;
  await uploadImage(file);
  input.value = '';
}

async function retryUploadImage() {
  if (!selectedFile.value) return;
  await uploadImage(selectedFile.value);
}

async function uploadImage(file: File) {
  errors.value.imageUrl = '';
  isUploadingImage.value = true;
  uploadProgress.value = 0;
  try {
    const resp = await productApi.uploadImage(file, (percent) => {
      uploadProgress.value = percent;
    });
    const imageUrl = resp.data?.image_url || '';
    formData.value.imageUrl = imageUrl;
    previewImageUrl.value = imageUrl;
    uploadProgress.value = 100;
  } catch (_error) {
    errors.value.imageUrl = '图片上传失败，请重试';
  } finally {
    isUploadingImage.value = false;
  }
}

function handleSubmit() {
  // Clear previous errors
  errors.value = {};

  // Validate form data
  try {
    productSchema.parse(formData.value);
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
