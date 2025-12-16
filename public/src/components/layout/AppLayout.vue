<template>
  <el-container class="layout-container">
    <el-aside width="200px" class="aside">
      <div class="logo">
        <img :src="systemSettings.logo || '/logo.png'" alt="Logo">
        <span>{{ systemSettings.systemName }}</span>
      </div>
      <el-menu
        class="menu"
        :default-active="activeMenu"
        :router="true"
        :background-color="systemSettings.theme?.navMode === 'dark' ? '#001529' : '#fff'"
        :text-color="systemSettings.theme?.navMode === 'dark' ? '#fff' : '#303133'"
        :active-text-color="systemSettings.theme?.primaryColor || '#409EFF'"
      >
        <el-menu-item index="/dashboard">
          <el-icon><Monitor /></el-icon>
          <span>仪表盘</span>
        </el-menu-item>
        <el-menu-item index="/users">
          <el-icon><User /></el-icon>
          <span>用户管理</span>
        </el-menu-item>
        <el-menu-item index="/roles">
          <el-icon><Lock /></el-icon>
          <span>角色管理</span>
        </el-menu-item>
        <el-menu-item index="/settings">
          <el-icon><Setting /></el-icon>
          <span>系统设置</span>
        </el-menu-item>
      </el-menu>
    </el-aside>
    <el-container>
      <el-header class="header">
        <div class="header-left">
          <el-button @click="toggleCollapse" link>
            <el-icon :size="20"><Expand /></el-icon>
          </el-button>
          <el-breadcrumb separator="/">
            <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
            <el-breadcrumb-item>{{ currentRoute }}</el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        <div class="header-right">
          <el-dropdown>
            <span class="user-info">
              <el-avatar :size="32" src="/avatar.png" />
              {{ userInfo.username }}
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item>个人中心</el-dropdown-item>
                <el-dropdown-item>修改密码</el-dropdown-item>
                <el-dropdown-item divided @click="handleLogout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>
      <el-main class="main">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Monitor, User, Lock, Setting, Expand } from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const isCollapse = ref(false)

// 模拟数据 - 实际应该从API获取
const systemSettings = ref({
  systemName: '后台管理系统',
  logo: '/assets/images/logo.png',
  theme: {
    primaryColor: '#409EFF',
    navMode: 'light',
    darkMode: false
  }
})

// 从 localStorage 获取用户信息
const userInfo = ref({
  username: 'Admin'
})

// 组件挂载时从 localStorage 读取用户信息
onMounted(() => {
  const storedUser = localStorage.getItem('user')

  if (storedUser) {
    try {
      const userData = JSON.parse(storedUser)
      userInfo.value = userData
    } catch (error) {
      console.error('解析用户信息失败:', error)
      // 解析失败也跳转到登录页
      router.push('/login')
    }
  } else {
    // 如果没有用户信息，跳转到登录页
    console.log('未找到用户信息，跳转到登录页')
    router.push('/login')
  }
})

const activeMenu = computed(() => route.path)
const currentRoute = computed(() => {
  const routeMap = {
    '/dashboard': '仪表盘',
    '/users': '用户管理',
    '/roles': '角色管理',
    '/settings': '系统设置'
  }
  return routeMap[route.path] || '未知页面'
})

const toggleCollapse = () => {
  isCollapse.value = !isCollapse.value
}

const handleLogout = async () => {
  try {
    // 调用登出API
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    await router.push('/login')
  } catch (error) {
    console.error('登出失败:', error)
  }
}
</script>

<style scoped>
.layout-container {
  height: 100vh;
}

.aside {
  background-color: v-bind('systemSettings.theme.navMode === "dark" ? "#001529" : "#fff"');
  border-right: 1px solid #e6e6e6;
}

.logo {
  height: 60px;
  padding: 10px;
  display: flex;
  align-items: center;
  color: v-bind('systemSettings.theme.navMode === "dark" ? "#fff" : "#303133"');
  border-bottom: 1px solid #e6e6e6;
}

.logo img {
  height: 32px;
  margin-right: 12px;
}

.menu {
  border-right: none;
}

.header {
  background-color: #fff;
  border-bottom: 1px solid #e6e6e6;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 20px;
}

.header-right {
  display: flex;
  align-items: center;
}

.user-info {
  display: flex;
  align-items: center;
  cursor: pointer;
  gap: 8px;
}

.main {
  background-color: #f5f7fa;
  padding: 20px;
}
</style>
