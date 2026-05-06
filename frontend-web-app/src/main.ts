import { createApp } from 'vue';
import { createPinia } from 'pinia';
import { VueQueryPlugin } from '@tanstack/vue-query';
import App from './App.vue';
import router from './router';
import './styles/global.css';

const app = createApp(App);

// 全局错误处理
app.config.errorHandler = (err, _instance, info) => {
  console.error('Vue Error Handler:', err, info);
  // TODO: 可以发送到错误监控服务
};

app.use(createPinia());
app.use(router);
app.use(VueQueryPlugin);

app.mount('#app');
