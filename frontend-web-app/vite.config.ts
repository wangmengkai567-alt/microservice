import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import { fileURLToPath, URL } from "node:url";
import Icons from "unplugin-icons/vite";

export default defineConfig({
  plugins: [
    vue(),
    Icons({
      compiler: "vue3",
      autoInstall: true,
    }),
  ],
  resolve: {
    alias: {
      "@": fileURLToPath(new URL("./src", import.meta.url)),
    },
  },
  server: {
    port: 5173,
    proxy: {
      "/api": {
        target: "http://localhost:8080",
        changeOrigin: true,
      },
    },
  },
  build: {
    rollupOptions: {
      output: {
        manualChunks: {
          "vue-vendor": ["vue", "vue-router", "pinia"],
          "query-vendor": ["@tanstack/vue-query", "axios"],
          "ui-vendor": ["@headlessui/vue"],
        },
      },
    },
  },
});
