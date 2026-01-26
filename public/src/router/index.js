import { createRouter, createWebHistory } from 'vue-router'
import LoginPage from '../components/LoginPage.vue'
import Dashboard from '../components/Dashboard.vue'
import AppLayout from '../components/layout/AppLayout.vue'

const routes = [
  {
    path: '/:lang?',
    component: { template: '<router-view></router-view>' },
    beforeEnter: (to, from, next) => {
      const lang = to.params.lang
      if (!lang) {
        const defaultLang = localStorage.getItem('language') || 'zh'
        // 避免路径重复，只在根路径时添加语言前缀
        const pathWithoutLang = to.fullPath
        return next(`/${defaultLang}${pathWithoutLang}`)
      }
      if (!['zh', 'en'].includes(lang)) {
        return next('/zh')
      }
      localStorage.setItem('language', lang)
      next()
    },
    children: [
      {
        path: 'login',
        name: 'Login',
        component: LoginPage
      },
      {
        path: '',
        component: AppLayout,
        redirect: to => {
          const lang = to.params.lang || localStorage.getItem('language') || 'zh'
          return `/${lang}/dashboard`
        },
        children: [
          {
            path: 'dashboard',
            name: 'Dashboard',
            component: Dashboard
          },
          {
            path: 'users',
            name: 'Users',
            component: () => import('../components/UserManagement.vue')
          },
          {
            path: 'roles',
            name: 'Roles',
            component: () => import('../components/RoleManagement.vue')
          },
          {
            path: 'settings',
            name: 'Settings',
            component: () => import('../components/SettingsPage.vue')
          },
          {
            path: 'operate-logs',
            name: 'OperateLogs',
            component: () => import('../components/OperateLogManagement.vue')
          }
        ]
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  const lang = to.params.lang || localStorage.getItem('language') || 'zh'

  // 如果没有语言参数，重定向到带语言参数的路径
  if (!to.params.lang && !to.path.startsWith(`/${lang}`)) {
    // 获取不带语言前缀的路径
    const pathWithoutLang = to.fullPath.replace(/^\/[a-z]{2}/, '')
    return next(`/${lang}${pathWithoutLang || '/dashboard'}`)
  }

  // 处理登录页面
  if (to.path.includes('/login')) {
    next()
  } else if (!token) {
    // 未登录，跳转到对应语言的登录页
    next(`/${lang}/login`)
  } else {
    next()
  }
})

export default router
