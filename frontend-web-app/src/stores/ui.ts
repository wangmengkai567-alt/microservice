import { defineStore } from 'pinia';
import { ref } from 'vue';

export const useUIStore = defineStore('ui', () => {
  // 状态
  const isLoading = ref(false);
  const error = ref<string | null>(null);
  const successMessage = ref<string | null>(null);

  // 设置加载状态
  function setLoading(loading: boolean) {
    isLoading.value = loading;
  }

  // 设置错误消息
  function setError(err: string | null) {
    error.value = err;

    // 自动在3秒后清除错误消息
    if (err) {
      setTimeout(() => {
        if (error.value === err) {
          error.value = null;
        }
      }, 3000);
    }
  }

  // 设置成功消息
  function setSuccess(message: string | null) {
    successMessage.value = message;

    // 自动在3秒后清除成功消息
    if (message) {
      setTimeout(() => {
        if (successMessage.value === message) {
          successMessage.value = null;
        }
      }, 3000);
    }
  }

  // 清除所有消息
  function clearMessages() {
    error.value = null;
    successMessage.value = null;
  }

  // 显示错误提示（便捷方法）
  function showError(message: string) {
    setError(message);
    console.error('UI Error:', message);
  }

  // 显示成功提示（便捷方法）
  function showSuccess(message: string) {
    setSuccess(message);
  }

  return {
    // 状态
    isLoading,
    error,
    successMessage,
    // 方法
    setLoading,
    setError,
    setSuccess,
    clearMessages,
    showError,
    showSuccess,
  };
});
