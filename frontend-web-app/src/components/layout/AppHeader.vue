<template>
  <header class="bg-white/80 backdrop-blur-md sticky top-0 z-50 border-b border-gray-100 shadow-sm">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between items-center h-16">
        <!-- Logo and Title -->
        <div class="flex items-center">
          <router-link
            to="/products"
            class="flex items-center space-x-2 text-xl font-extrabold text-indigo-600 hover:text-indigo-700 transition-colors focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 rounded-lg"
          >
            <div
              class="w-8 h-8 bg-gradient-to-br from-indigo-500 to-blue-600 rounded-lg flex items-center justify-center shadow-sm text-white"
            >
              <svg
                class="w-5 h-5"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z"
                ></path>
              </svg>
            </div>
            <span>电商平台</span>
          </router-link>
        </div>

        <!-- Navigation -->
        <nav class="hidden md:flex space-x-1" role="navigation" aria-label="主导航">
          <router-link
            to="/products"
            class="text-gray-600 hover:text-indigo-600 px-4 py-2 rounded-xl text-sm font-medium transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-indigo-500"
            active-class="text-indigo-700 bg-indigo-50 shadow-sm"
          >
            商品列表
          </router-link>
          <router-link
            to="/orders"
            class="text-gray-600 hover:text-indigo-600 px-4 py-2 rounded-xl text-sm font-medium transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-indigo-500"
            active-class="text-indigo-700 bg-indigo-50 shadow-sm"
          >
            我的订单
          </router-link>
        </nav>

        <!-- User Info and Logout -->
        <div v-if="isAuthenticated" class="flex items-center space-x-4">
          <div class="hidden sm:flex items-center space-x-2">
            <div
              class="w-8 h-8 rounded-full bg-indigo-100 flex items-center justify-center text-indigo-700 font-bold text-sm"
            >
              {{ user?.username?.charAt(0).toUpperCase() || 'U' }}
            </div>
            <span class="text-sm font-medium text-gray-700">
              {{ user?.username }}
            </span>
          </div>
          <button
            @click="handleLogout"
            class="inline-flex items-center px-4 py-2 border border-gray-200 text-sm font-medium rounded-xl text-gray-600 bg-white hover:bg-gray-50 hover:text-red-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 transition-all duration-200 shadow-sm"
            :disabled="isLoggingOut"
          >
            <svg
              class="w-4 h-4 mr-1.5"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"
              ></path>
            </svg>
            {{ isLoggingOut ? '退出中...' : '退出登录' }}
          </button>
        </div>

        <!-- Mobile menu button -->
        <div class="md:hidden">
          <button
            @click="mobileMenuOpen = !mobileMenuOpen"
            type="button"
            class="inline-flex items-center justify-center p-2 rounded-md text-gray-700 hover:text-indigo-600 hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-indigo-500"
            :aria-expanded="mobileMenuOpen"
            aria-label="打开菜单"
          >
            <svg
              class="h-6 w-6"
              :class="{ hidden: mobileMenuOpen, block: !mobileMenuOpen }"
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M4 6h16M4 12h16M4 18h16"
              />
            </svg>
            <svg
              class="h-6 w-6"
              :class="{ block: mobileMenuOpen, hidden: !mobileMenuOpen }"
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M6 18L18 6M6 6l12 12"
              />
            </svg>
          </button>
        </div>
      </div>
    </div>

    <!-- Mobile menu -->
    <div v-if="mobileMenuOpen" class="md:hidden border-t border-gray-200">
      <div class="px-2 pt-2 pb-3 space-y-1">
        <router-link
          to="/products"
          class="block px-3 py-2 rounded-md text-base font-medium text-gray-700 hover:text-indigo-600 hover:bg-gray-50"
          active-class="text-indigo-600 bg-indigo-50"
          @click="mobileMenuOpen = false"
        >
          商品列表
        </router-link>
        <router-link
          to="/orders"
          class="block px-3 py-2 rounded-md text-base font-medium text-gray-700 hover:text-indigo-600 hover:bg-gray-50"
          active-class="text-indigo-600 bg-indigo-50"
          @click="mobileMenuOpen = false"
        >
          我的订单
        </router-link>
      </div>
      <div v-if="isAuthenticated" class="pt-4 pb-3 border-t border-gray-200">
        <div class="px-4 flex items-center justify-between">
          <span class="text-sm text-gray-700">
            欢迎，<span class="font-medium">{{ user?.username }}</span>
          </span>
          <button
            @click="handleLogout"
            class="px-4 py-2 text-sm font-medium text-white bg-indigo-600 rounded-md hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
            :disabled="isLoggingOut"
          >
            {{ isLoggingOut ? '退出中...' : '退出登录' }}
          </button>
        </div>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '@/stores/auth';
import { storeToRefs } from 'pinia';

const router = useRouter();
const authStore = useAuthStore();
const { user, isAuthenticated } = storeToRefs(authStore);

const mobileMenuOpen = ref(false);
const isLoggingOut = ref(false);

async function handleLogout() {
  if (isLoggingOut.value) return;

  isLoggingOut.value = true;
  try {
    await authStore.logout();
    router.push('/login');
  } catch (error) {
    console.error('Logout failed:', error);
  } finally {
    isLoggingOut.value = false;
  }
}
</script>
