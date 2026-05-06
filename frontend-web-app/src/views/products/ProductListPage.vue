<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import MainLayout from '@/components/layout/MainLayout.vue';
import ProductCard from '@/components/features/ProductCard.vue';
import Pagination from '@/components/common/Pagination.vue';
import Button from '@/components/common/Button.vue';
import Loading from '@/components/common/Loading.vue';
import Modal from '@/components/common/Modal.vue';
import { useProducts } from '@/composables/useProducts';
import type { Product } from '@/types/product';

const router = useRouter();
const currentPage = ref(1);
const pageSize = 10;

const { products, total, isLoading, error, deleteProduct, isDeleting } = useProducts(
  currentPage.value,
  pageSize
);

// Delete confirmation modal
const showDeleteModal = ref(false);
const productToDelete = ref<Product | null>(null);

function handleEdit(product: Product) {
  router.push(`/products/${product.id}/edit`);
}

function confirmDelete(product: Product) {
  productToDelete.value = product;
  showDeleteModal.value = true;
}

async function handleDelete() {
  if (!productToDelete.value) return;

  try {
    await deleteProduct(productToDelete.value.id);
    showDeleteModal.value = false;
    productToDelete.value = null;
  } catch (err) {
    console.error('Failed to delete product:', err);
  }
}

function cancelDelete() {
  showDeleteModal.value = false;
  productToDelete.value = null;
}

function handlePageChange(page: number) {
  currentPage.value = page;
}

function goToCreateProduct() {
  router.push('/products/new');
}
</script>

<template>
  <MainLayout>
    <div
      class="mb-8 flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 bg-white p-6 rounded-2xl shadow-sm border border-gray-100"
    >
      <div>
        <h1
          class="text-2xl sm:text-3xl font-extrabold text-gray-900 tracking-tight flex items-center"
        >
          <svg
            class="w-8 h-8 mr-3 text-indigo-500"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"
            ></path>
          </svg>
          商品列表
        </h1>
        <p class="mt-2 text-sm text-gray-500">管理您的所有商品库存、价格和详情信息</p>
      </div>
      <Button
        @click="goToCreateProduct"
        class="shadow-md hover:shadow-lg transition-all hover:-translate-y-0.5"
      >
        <svg
          class="w-5 h-5 mr-1.5"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
          xmlns="http://www.w3.org/2000/svg"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M12 6v6m0 0v6m0-6h6m-6 0H6"
          ></path>
        </svg>
        创建商品
      </Button>
    </div>

    <!-- 加载状态 -->
    <div v-if="isLoading" class="flex justify-center py-12">
      <Loading />
    </div>

    <!-- 错误提示 -->
    <div v-else-if="error" class="bg-red-50 border border-red-200 rounded-md p-4">
      <p class="text-sm text-red-800">{{ error }}</p>
    </div>

    <!-- 空列表 -->
    <div
      v-else-if="products.length === 0"
      class="text-center py-20 bg-white rounded-2xl shadow-sm border border-gray-100 flex flex-col items-center justify-center"
    >
      <div class="w-24 h-24 bg-gray-50 rounded-full flex items-center justify-center mb-6">
        <svg class="h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="1.5"
            d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4"
          />
        </svg>
      </div>
      <h3 class="text-lg font-medium text-gray-900 mb-2">暂无商品</h3>
      <p class="text-gray-500 max-w-sm mx-auto mb-6">
        您的商品库目前是空的。开始添加商品来丰富您的店铺吧！
      </p>
      <Button @click="goToCreateProduct" class="shadow-sm hover:shadow transition-all">
        <svg
          class="w-5 h-5 mr-1.5"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
          xmlns="http://www.w3.org/2000/svg"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M12 6v6m0 0v6m0-6h6m-6 0H6"
          ></path>
        </svg>
        创建第一个商品
      </Button>
    </div>

    <!-- 商品列表 -->
    <div v-else>
      <div class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3 mb-8">
        <ProductCard
          v-for="product in products"
          :key="product.id"
          :product="product"
          :show-actions="true"
          @edit="handleEdit"
          @delete="confirmDelete"
        />
      </div>

      <!-- 分页 -->
      <Pagination
        v-if="total > pageSize"
        :current-page="currentPage"
        :total-pages="Math.ceil(total / pageSize)"
        @page-change="handlePageChange"
      />
    </div>

    <!-- 删除确认对话框 -->
    <Modal :is-open="showDeleteModal" title="确认删除" @close="cancelDelete">
      <template #default>
        <p class="text-sm text-gray-600">
          确定要删除商品
          <span class="font-medium text-gray-900">{{ productToDelete?.name }}</span>
          吗？此操作无法撤销。
        </p>
      </template>
      <template #footer>
        <div class="flex justify-end space-x-3">
          <Button variant="secondary" @click="cancelDelete" :disabled="isDeleting"> 取消 </Button>
          <Button @click="handleDelete" :loading="isDeleting" variant="danger"> 确认删除 </Button>
        </div>
      </template>
    </Modal>
  </MainLayout>
</template>
