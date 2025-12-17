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
        <el-menu-item :index="`/${locale}/dashboard`">
          <el-icon><Monitor /></el-icon>
          <span>{{ t('dashboard') }}</span>
        </el-menu-item>
        <el-menu-item :index="`/${locale}/users`">
          <el-icon><User /></el-icon>
          <span>{{ t('userManagement') }}</span>
        </el-menu-item>
        <el-menu-item :index="`/${locale}/roles`">
          <el-icon><Lock /></el-icon>
          <span>{{ t('roleManagement') }}</span>
        </el-menu-item>
        <el-menu-item :index="`/${locale}/settings`">
          <el-icon><Setting /></el-icon>
          <span>{{ t('systemSettings') }}</span>
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
            <el-breadcrumb-item :to="{ path: '/' }">{{ t('home') }}</el-breadcrumb-item>
            <el-breadcrumb-item>{{ t(currentRoute) }}</el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        <div class="header-right">
          <!-- 语言切换 -->
          <el-dropdown @command="handleLanguageChange" style="margin-right: 20px;">
            <span class="language-selector">
              <el-icon><i class="fas fa-globe"></i></el-icon>
              {{ locale === 'zh' ? '中文' : 'English' }}
              <el-icon class="el-icon--right"><arrow-down /></el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="zh">中文</el-dropdown-item>
                <el-dropdown-item command="en">English</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>

          <!-- 用户菜单 -->
          <el-dropdown>
            <span class="user-info">
              <el-avatar :size="32" src="/avatar.png" />
              {{ userInfo.username }}
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item>{{ t('userCenter') }}</el-dropdown-item>
                <el-dropdown-item @click="showChangePasswordDialog">{{ t('changePassword') }}</el-dropdown-item>
                <el-dropdown-item divided @click="handleLogout">{{ t('logout') }}</el-dropdown-item>
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

  <!-- 修改密码对话框 -->
  <el-dialog
    v-model="changePasswordVisible"
    :title="t('changePassword')"
    width="400px"
    :close-on-click-modal="false"
  >
    <el-form
      ref="passwordFormRef"
      :model="passwordForm"
      :rules="passwordRules"
      label-width="100px"
    >
      <el-form-item :label="t('oldPassword')" prop="oldPassword">
        <el-input
          v-model="passwordForm.oldPassword"
          type="password"
          show-password
          :placeholder="t('enterOldPassword')"
        />
      </el-form-item>
      <el-form-item :label="t('newPassword')" prop="newPassword">
        <el-input
          v-model="passwordForm.newPassword"
          type="password"
          show-password
          :placeholder="t('enterNewPassword')"
        />
      </el-form-item>
      <el-form-item :label="t('confirmNewPassword')" prop="confirmPassword">
        <el-input
          v-model="passwordForm.confirmPassword"
          type="password"
          show-password
          :placeholder="t('enterNewPasswordAgain')"
        />
      </el-form-item>
    </el-form>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="changePasswordVisible = false">{{ t('cancel') }}</el-button>
        <el-button type="primary" @click="handleChangePassword">{{ t('confirm') }}</el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, computed, onMounted, reactive } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Monitor, User, Lock, Setting, Expand, ArrowDown } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'
import md5 from 'js-md5'

const { t, locale } = useI18n()

const route = useRoute()
const router = useRouter()
const isCollapse = ref(false)

// 修改密码相关
const changePasswordVisible = ref(false)
const passwordFormRef = ref()
const passwordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

// 密码验证规则
const validateConfirmPassword = (rule, value, callback) => {
  if (value === '') {
    callback(new Error(t('enterNewPasswordAgain')))
  } else if (value !== passwordForm.newPassword) {
    callback(new Error(t('passwordNotMatch')))
  } else {
    callback()
  }
}

const passwordRules = computed(() => ({
  oldPassword: [
    { required: true, message: t('enterOldPassword'), trigger: 'blur' }
  ],
  newPassword: [
    { required: true, message: t('enterNewPassword'), trigger: 'blur' },
    { min: 6, max: 20, message: t('passwordLengthLimit'), trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: t('enterNewPasswordAgain'), trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ]
}))

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

const activeMenu = computed(() => {
  // 返回完整路径，包含语言前缀
  return route.path
})
const currentRoute = computed(() => {
  // 移除语言前缀以获取实际路径
  const path = route.path.replace(/^\/(zh|en)/, '')
  const routeMap = {
    '/dashboard': 'dashboard',
    '/users': 'userManagement',
    '/roles': 'roleManagement',
    '/settings': 'systemSettings'
  }
  return routeMap[path] || 'unknownPage'
})

// 处理语言切换
const handleLanguageChange = async (lang) => {
  locale.value = lang
  localStorage.setItem('language', lang)

  // 更新路由到新的语言路径
  const currentPath = route.path
  const newPath = currentPath.replace(/^\/(zh|en)/, `/${lang}`)
  await router.push(newPath)

  // 发送语言切换请求到后端
  if (localStorage.getItem('token')) {
    fetch('/api/admin/v1/user/changeLang', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`,
        'Accept-Language': lang
      },
      body: JSON.stringify({ language: lang })
    }).catch(err => console.error('Language change failed:', err))
  }
}

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

// 显示修改密码对话框
const showChangePasswordDialog = () => {
  changePasswordVisible.value = true
  // 重置表单
  passwordForm.oldPassword = ''
  passwordForm.newPassword = ''
  passwordForm.confirmPassword = ''
  // 清除验证状态
  if (passwordFormRef.value) {
    passwordFormRef.value.resetFields()
  }
}

// 处理修改密码
const handleChangePassword = async () => {
  if (!passwordFormRef.value) return

  await passwordFormRef.value.validate(async (valid) => {
    if (valid) {
      try {
        // 获取token
        const token = localStorage.getItem('token')
        if (!token) {
          ElMessage.error('登录已过期，请重新登录')
          router.push('/login')
          return
        }

        // 调用修改密码API
        const response = await fetch('/api/admin/v1/user/change_pwd', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
          },
          body: JSON.stringify({
            old_password: md5(passwordForm.oldPassword),
            new_password: md5(passwordForm.newPassword),
            confirm_password: md5(passwordForm.confirmPassword),
          })
        })

        const data = await response.json()

        if (response.ok && data.code === 200) {
          ElMessage.success('密码修改成功，请重新登录')
          changePasswordVisible.value = false
          // 清除登录信息
          localStorage.removeItem('token')
          localStorage.removeItem('user')
          // 跳转到登录页面
          router.push('/login')
        } else {
          ElMessage.error(data.message || '修改密码失败')
        }
      } catch (error) {
        console.error('修改密码失败:', error)
        ElMessage.error('修改密码失败，请稍后重试')
      }
    } else {
      console.log('表单验证失败')
    }
  })
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

.language-selector {
  display: flex;
  align-items: center;
  cursor: pointer;
  gap: 4px;
  padding: 4px 8px;
  border-radius: 4px;
  transition: background-color 0.3s;
}

.language-selector:hover {
  background-color: #f5f7fa;
}
</style>
