import { createRouter, createWebHistory } from 'vue-router';
import { useAuthStore } from '@/stores/auth';

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      redirect: '/products',
    },
    {
      path: '/login',
      name: 'Login',
      component: () => import('@/views/auth/LoginPage.vue'),
      meta: { title: '登录' },
    },
    {
      path: '/register',
      name: 'Register',
      component: () => import('@/views/auth/RegisterPage.vue'),
      meta: { title: '注册' },
    },
    {
      path: '/products',
      name: 'ProductList',
      component: () => import('@/views/products/ProductListPage.vue'),
      meta: { title: '商品列表', requiresAuth: true },
    },
    {
      path: '/products/new',
      name: 'ProductCreate',
      component: () => import('@/views/products/ProductFormPage.vue'),
      meta: { title: '创建商品', requiresAuth: true },
    },
    {
      path: '/products/:id',
      name: 'ProductDetail',
      component: () => import('@/views/products/ProductDetailPage.vue'),
      meta: { title: '商品详情', requiresAuth: true },
    },
    {
      path: '/products/:id/edit',
      name: 'ProductEdit',
      component: () => import('@/views/products/ProductFormPage.vue'),
      meta: { title: '编辑商品', requiresAuth: true },
    },
    {
      path: '/orders',
      name: 'OrderList',
      component: () => import('@/views/orders/OrderListPage.vue'),
      meta: { title: '订单列表', requiresAuth: true },
    },
    {
      path: '/orders/:id',
      name: 'OrderDetail',
      component: () => import('@/views/orders/OrderDetailPage.vue'),
      meta: { title: '订单详情', requiresAuth: true },
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'NotFound',
      component: () => import('@/views/NotFoundPage.vue'),
      meta: { title: '页面未找到' },
    },
  ],
});

// 导航守卫
router.beforeEach((to, _from, next) => {
  const authStore = useAuthStore();

  // 需要认证的路由
  const requiresAuth = to.meta.requiresAuth;

  if (requiresAuth && !authStore.isAuthenticated) {
    // 保存原始路径，登录后跳转回来
    next({
      path: '/login',
      query: { redirect: to.fullPath },
    });
  } else {
    next();
  }
});

// 更新页面标题
router.afterEach((to) => {
  const title = to.meta.title as string;
  document.title = title ? `${title} - 电商平台` : '电商平台';
});

export default router;
