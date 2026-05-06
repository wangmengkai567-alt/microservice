<template>
  <div
    class="min-h-screen flex flex-col bg-gradient-to-br from-gray-50 to-gray-100 font-sans text-gray-900"
  >
    <AppHeader />
    <main class="flex-1 max-w-7xl w-full mx-auto px-4 sm:px-6 lg:px-8 py-8 md:py-10">
      <slot />
    </main>
    <AppFooter />
  </div>
</template>

<script setup lang="ts">
import AppHeader from './AppHeader.vue';
import AppFooter from './AppFooter.vue';
</script>
