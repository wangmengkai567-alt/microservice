<template>
  <MainLayout>
    <div class="max-w-3xl mx-auto">
      <div
        class="mb-8 bg-white p-6 rounded-2xl shadow-sm border border-gray-100 flex items-center justify-between"
      >
        <div class="flex items-center">
          <button
            @click="router.push('/products')"
            class="mr-4 p-2 rounded-full hover:bg-gray-100 text-gray-500 transition-colors"
          >
            <svg
              class="w-6 h-6"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M10 19l-7-7m0 0l7-7m-7 7h18"
              ></path>
            </svg>
          </button>
          <div>
            <h1 class="text-2xl font-extrabold text-gray-900 tracking-tight">
              {{ isEditMode ? '编辑商品' : '创建新商品' }}
            </h1>
            <p class="mt-1 text-sm text-gray-500">
              {{ isEditMode ? '修改商品信息并保存' : '填写下方表单以添加新商品到您的店铺' }}
            </p>
          </div>
        </div>
      </div>

      <div
        v-if="isEditMode && isLoading"
        class="flex justify-center py-20 bg-white rounded-2xl shadow-sm border border-gray-100"
      >
        <Loading />
      </div>

      <div
        v-else-if="isEditMode && error"
        class="bg-red-50 border border-red-100 rounded-2xl p-6 text-center"
      >
        <svg
          class="w-12 h-12 text-red-400 mx-auto mb-3"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
          xmlns="http://www.w3.org/2000/svg"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
          ></path>
        </svg>
        <p class="text-red-800 font-medium">{{ error }}</p>
        <Button @click="router.push('/products')" variant="secondary" class="mt-6">
          返回商品列表
        </Button>
      </div>

      <div v-else class="bg-white rounded-2xl shadow-sm border border-gray-100 p-8 relative overflow-hidden">
        <div class="absolute top-0 left-0 w-full h-1 bg-gradient-to-r from-indigo-500 to-blue-500"></div>
        <ProductForm
          :initial-data="product || undefined"
          :is-submitting="isSubmitting"
          :submit-text="isEditMode ? '保存修改' : '确认创建'"
          @submit="handleSubmit"
          @cancel="handleCancel"
        />
      </div>
    </div>
  </MainLayout>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import MainLayout from '@/components/layout/MainLayout.vue';
import ProductForm from '@/components/features/ProductForm.vue';
import Button from '@/components/common/Button.vue';
import Loading from '@/components/common/Loading.vue';
import { useProducts, useProductDetail } from '@/composables/useProducts';
import type { ProductFormData } from '@/types/product';

const router = useRouter();
const route = useRoute();

const productId = computed(() => route.params.id as string);
const isEditMode = computed(() => !!productId.value && productId.value !== 'new');

// For edit mode
const { product, isLoading, error } = isEditMode.value
  ? useProductDetail(productId.value)
  : { product: ref(null), isLoading: ref(false), error: ref(null) };

// For create/update operations
const { createProduct, updateProduct } = useProducts();
const isSubmitting = ref(false);

async function handleSubmit(data: ProductFormData) {
  isSubmitting.value = true;
  try {
    if (isEditMode.value) {
      await updateProduct(productId.value, data);
    } else {
      await createProduct(data);
    }
    router.push('/products');
  } catch (err) {
    console.error('Failed to save product:', err);
  } finally {
    isSubmitting.value = false;
  }
}

function handleCancel() {
  router.push('/products');
}
</script>
