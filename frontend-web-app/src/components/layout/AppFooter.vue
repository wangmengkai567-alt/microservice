<template>
  <footer class="bg-white/60 backdrop-blur-sm border-t border-gray-100 mt-auto py-8">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex flex-col md:flex-row justify-between items-center space-y-4 md:space-y-0">
        <div class="flex items-center space-x-2">
          <div
            class="w-6 h-6 bg-indigo-100 text-indigo-600 rounded-md flex items-center justify-center"
          >
            <svg
              class="w-3.5 h-3.5"
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
          <span class="text-sm font-medium text-gray-600"
            >© {{ currentYear }} 电商平台. 保留所有权利.</span
          >
        </div>
        <div class="flex space-x-6">
          <a
            href="#"
            class="text-sm font-medium text-gray-500 hover:text-indigo-600 transition-colors focus:outline-none focus:ring-2 focus:ring-indigo-500 rounded-md px-2 py-1"
          >
            关于我们
          </a>
          <a
            href="#"
            class="text-sm font-medium text-gray-500 hover:text-indigo-600 transition-colors focus:outline-none focus:ring-2 focus:ring-indigo-500 rounded-md px-2 py-1"
          >
            联系方式
          </a>
          <a
            href="#"
            class="text-sm font-medium text-gray-500 hover:text-indigo-600 transition-colors focus:outline-none focus:ring-2 focus:ring-indigo-500 rounded-md px-2 py-1"
          >
            隐私政策
          </a>
        </div>
      </div>
    </div>
  </footer>
</template>

<script setup lang="ts">
import { computed } from 'vue';

const currentYear = computed(() => new Date().getFullYear());
</script>
