import { createApp } from 'vue'
import ElementPlus from 'element-plus'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import 'element-plus/dist/index.css'
import App from './App.vue'
import router from './router'
import axios from 'axios'
import i18n from './i18n'

// 创建Vue应用
const app = createApp(App)

// 使用Element Plus并配置语言
app.use(ElementPlus, {
  locale: i18n.global.locale.value === 'zh' ? 'zh-cn' : 'en',
})

// 使用i18n
app.use(i18n)

// 注册所有Element Plus图标
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

// 使用路由
app.use(router)

// 配置axios
axios.defaults.baseURL = 'http://192.168.56.109/'
axios.defaults.headers.common['Content-Type'] = 'application/json'

// 请求拦截器
axios.interceptors.request.use(
  config => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    // 添加 Accept-Language header
    const language = localStorage.getItem('language') || 'zh'
    config.headers['Accept-Language'] = language === 'zh' ? 'zh-CN' : 'en-US'
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// 响应拦截器
axios.interceptors.response.use(
  response => {
    return response
  },
  error => {
    if (error.response && error.response.status === 401) {
      // 未授权，清除token并跳转到登录页
      localStorage.removeItem('token')
      router.push('/login')
    }
    return Promise.reject(error)
  }
)

// 全局属性
app.config.globalProperties.$axios = axios

// 挂载应用
app.mount('#app')
