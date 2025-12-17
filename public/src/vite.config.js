import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      'vue': 'vue/dist/vue.esm-bundler.js',
    },
  },
  server: {
    port: 3000,
    proxy: {
      // 代理API请求到后端
      '/api': {
        target: 'http://192.168.56.109', // 后端服务地址
        changeOrigin: true,
        secure: false
      }
    }
  },
  build: {
    outDir: 'dist',
    assetsDir: 'assets',
    sourcemap: false
  }
})
