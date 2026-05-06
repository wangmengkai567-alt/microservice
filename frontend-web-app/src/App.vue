<script setup lang="ts">
import Toast from '@/components/common/Toast.vue';
</script>

<template>
  <div id="app" class="min-h-screen bg-gray-50">
    <router-view />
    <Toast />
  </div>
</template>
