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
                <el-dropdown-item @click="showChangePasswordDialog">修改密码</el-dropdown-item>
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

  <!-- 修改密码对话框 -->
  <el-dialog
    v-model="changePasswordVisible"
    title="修改密码"
    width="400px"
    :close-on-click-modal="false"
  >
    <el-form
      ref="passwordFormRef"
      :model="passwordForm"
      :rules="passwordRules"
      label-width="100px"
    >
      <el-form-item label="原密码" prop="oldPassword">
        <el-input
          v-model="passwordForm.oldPassword"
          type="password"
          show-password
          placeholder="请输入原密码"
        />
      </el-form-item>
      <el-form-item label="新密码" prop="newPassword">
        <el-input
          v-model="passwordForm.newPassword"
          type="password"
          show-password
          placeholder="请输入新密码"
        />
      </el-form-item>
      <el-form-item label="确认新密码" prop="confirmPassword">
        <el-input
          v-model="passwordForm.confirmPassword"
          type="password"
          show-password
          placeholder="请再次输入新密码"
        />
      </el-form-item>
    </el-form>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="changePasswordVisible = false">取消</el-button>
        <el-button type="primary" @click="handleChangePassword">确认</el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, computed, onMounted, reactive } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Monitor, User, Lock, Setting, Expand } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import md5 from 'js-md5'

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
    callback(new Error('请再次输入新密码'))
  } else if (value !== passwordForm.newPassword) {
    callback(new Error('两次输入密码不一致!'))
  } else {
    callback()
  }
}

const passwordRules = reactive({
  oldPassword: [
    { required: true, message: '请输入原密码', trigger: 'blur' }
  ],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, max: 20, message: '密码长度在 6 到 20 个字符', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请再次输入新密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ]
})

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
        const response = await fetch('/api/admin/v1/user/changePwd', {
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
</style>
