<script setup lang="ts">
import { ref } from 'vue';
import { useAuth } from '@/composables/useAuth';

const { register, isLoading, error } = useAuth();

const username = ref('');
const email = ref('');
const password = ref('');
const validationError = ref('');

async function handleSubmit() {
  // 表单验证
  if (!username.value) {
    validationError.value = '用户名不能为空';
    return;
  }
  if (!email.value) {
    validationError.value = '邮箱不能为空';
    return;
  }
  if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email.value)) {
    validationError.value = '邮箱格式不正确';
    return;
  }
  if (password.value.length < 6) {
    validationError.value = '密码长度至少为6个字符';
    return;
  }

  validationError.value = '';

  try {
    await register({
      username: username.value,
      password: password.value,
    });
  } catch (err) {
    // 错误已在useAuth中处理
  }
}
</script>

<template>
  <div
    class="min-h-screen flex items-center justify-center bg-gradient-to-br from-indigo-50 via-white to-blue-50 py-12 px-4 sm:px-6 lg:px-8 relative overflow-hidden"
  >
    <!-- Background decorations -->
    <div
      class="absolute -top-40 -left-40 w-96 h-96 bg-blue-100 rounded-full mix-blend-multiply filter blur-3xl opacity-50 animate-blob"
    ></div>
    <div
      class="absolute -bottom-40 -right-40 w-96 h-96 bg-indigo-100 rounded-full mix-blend-multiply filter blur-3xl opacity-50 animate-blob animation-delay-2000"
    ></div>

    <div
      class="max-w-md w-full space-y-8 bg-white/80 backdrop-blur-lg p-10 rounded-2xl shadow-[0_8px_30px_rgb(0,0,0,0.04)] border border-white/20 relative z-10"
    >
      <div>
        <div
          class="mx-auto h-12 w-12 bg-indigo-100 text-indigo-600 rounded-xl flex items-center justify-center mb-6 shadow-sm"
        >
          <svg
            class="w-8 h-8"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z"
            ></path>
          </svg>
        </div>
        <h2 class="text-center text-3xl font-extrabold text-gray-900 tracking-tight">注册新账号</h2>
        <p class="mt-2 text-center text-sm text-gray-500">加入我们，开启您的电商之旅</p>
      </div>
      <form class="mt-8 space-y-6" @submit.prevent="handleSubmit">
        <div class="space-y-4">
          <div>
            <label for="username" class="block text-sm font-medium text-gray-700 mb-1"
              >用户名</label
            >
            <input
              id="username"
              v-model="username"
              type="text"
              required
              class="appearance-none block w-full px-4 py-3 border border-gray-200 rounded-xl shadow-sm placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-all duration-200 sm:text-sm bg-white/50 focus:bg-white"
              placeholder="请输入您的用户名"
            />
          </div>
          <div>
            <label for="email" class="block text-sm font-medium text-gray-700 mb-1">邮箱</label>
            <input
              id="email"
              v-model="email"
              type="email"
              required
              class="appearance-none block w-full px-4 py-3 border border-gray-200 rounded-xl shadow-sm placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-all duration-200 sm:text-sm bg-white/50 focus:bg-white"
              placeholder="请输入您的邮箱"
            />
          </div>
          <div>
            <label for="password" class="block text-sm font-medium text-gray-700 mb-1">密码</label>
            <input
              id="password"
              v-model="password"
              type="password"
              required
              class="appearance-none block w-full px-4 py-3 border border-gray-200 rounded-xl shadow-sm placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-all duration-200 sm:text-sm bg-white/50 focus:bg-white"
              placeholder="请输入密码（至少6个字符）"
            />
          </div>
        </div>

        <div
          v-if="validationError || error"
          class="bg-red-50 border border-red-100 text-red-600 text-sm p-3 rounded-lg text-center flex items-center justify-center"
        >
          <svg class="w-4 h-4 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
            ></path>
          </svg>
          {{ validationError || error }}
        </div>

        <div>
          <button
            type="submit"
            :disabled="isLoading"
            class="group relative w-full flex justify-center py-3 px-4 border border-transparent text-sm font-medium rounded-xl text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed shadow-md hover:shadow-lg transition-all duration-200"
          >
            <span class="absolute left-0 inset-y-0 flex items-center pl-3">
              <svg
                v-if="!isLoading"
                class="h-5 w-5 text-indigo-500 group-hover:text-indigo-400 transition-colors"
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 20 20"
                fill="currentColor"
                aria-hidden="true"
              >
                <path
                  fill-rule="evenodd"
                  d="M8 4a4 4 0 100 8 4 4 0 000-8zM2 8a6 6 0 1110.89 3.476l4.817 4.817a1 1 0 01-1.414 1.414l-4.816-4.816A6 6 0 012 8z"
                  clip-rule="evenodd"
                />
              </svg>
              <svg
                v-else
                class="animate-spin h-5 w-5 text-white"
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
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
            </span>
            {{ isLoading ? '正在注册...' : '立即注册' }}
          </button>
        </div>

        <div class="text-center mt-6">
          <router-link
            to="/login"
            class="text-sm font-medium text-indigo-600 hover:text-indigo-500 transition-colors"
          >
            已有账号？立即登录 &rarr;
          </router-link>
        </div>
      </form>
    </div>
  </div>
</template>
