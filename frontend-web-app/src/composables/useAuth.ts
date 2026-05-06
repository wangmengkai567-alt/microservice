import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '@/stores/auth';
import { useUIStore } from '@/stores/ui';
import type { LoginDto, RegisterDto } from '@/types/user';

export function useAuth() {
  const authStore = useAuthStore();
  const uiStore = useUIStore();
  const router = useRouter();

  const isLoading = ref(false);
  const error = ref<string | null>(null);

  // 计算属性
  const user = computed(() => authStore.user);
  const isAuthenticated = computed(() => authStore.isAuthenticated);

  // 登录
  async function login(credentials: LoginDto) {
    try {
      isLoading.value = true;
      error.value = null;

      await authStore.login(credentials);

      uiStore.showSuccess('登录成功');

      // 跳转到目标页面或商品列表
      const redirect = router.currentRoute.value.query.redirect as string;
      router.push(redirect || '/products');
    } catch (err: any) {
      const message = err.message || '登录失败';
      error.value = message;
      uiStore.showError(message);
      throw err;
    } finally {
      isLoading.value = false;
    }
  }

  // 注册
  async function register(data: RegisterDto) {
    try {
      isLoading.value = true;
      error.value = null;

      await authStore.register(data);

      uiStore.showSuccess('注册成功，请登录');
      router.push('/login');
    } catch (err: any) {
      const message = err.message || '注册失败';
      error.value = message;
      uiStore.showError(message);
      throw err;
    } finally {
      isLoading.value = false;
    }
  }

  // 退出登录
  function logout() {
    authStore.logout();
    uiStore.showSuccess('已退出登录');
    router.push('/login');
  }

  return {
    user,
    isAuthenticated,
    isLoading,
    error,
    login,
    register,
    logout,
  };
}
