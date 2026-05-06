<template>
  <MainLayout>
    <div class="flex flex-col items-center justify-center py-20 px-4 sm:px-6 lg:px-8 bg-white rounded-3xl shadow-xl border border-gray-100 max-w-4xl mx-auto my-12 relative overflow-hidden">
      <!-- Decorative background elements -->
      <div class="absolute -top-24 -right-24 w-64 h-64 bg-indigo-50 rounded-full mix-blend-multiply filter blur-3xl opacity-70"></div>
      <div class="absolute -bottom-24 -left-24 w-64 h-64 bg-blue-50 rounded-full mix-blend-multiply filter blur-3xl opacity-70"></div>

      <div class="text-center relative z-10">
        <h1 class="text-9xl font-black text-transparent bg-clip-text bg-gradient-to-r from-indigo-500 to-blue-500 mb-4 tracking-tighter drop-shadow-sm">404</h1>
        <div class="inline-flex items-center justify-center w-20 h-20 rounded-full bg-indigo-50 mb-6">
          <svg class="w-10 h-10 text-indigo-500" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
          </svg>
        </div>
        <h2 class="text-3xl sm:text-4xl font-extrabold text-gray-900 tracking-tight mb-3">页面未找到</h2>
        <p class="text-lg text-gray-500 max-w-md mx-auto mb-10">抱歉，您访问的页面不存在或已被移除。请检查URL是否正确，或者返回主页继续浏览。</p>
        
        <div class="flex flex-col sm:flex-row items-center justify-center gap-4">
          <Button @click="router.back()" variant="secondary" class="w-full sm:w-auto flex items-center justify-center">
            <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"></path>
            </svg>
            返回上一页
          </Button>
          <Button @click="router.push('/products')" class="w-full sm:w-auto shadow-md hover:shadow-lg transition-all hover:-translate-y-0.5 flex items-center justify-center">
            <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"></path>
            </svg>
            返回首页
          </Button>
        </div>
      </div>
    </div>
  </MainLayout>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router';
import MainLayout from '@/components/layout/MainLayout.vue';
import Button from '@/components/common/Button.vue';

const router = useRouter();
</script>
