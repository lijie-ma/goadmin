<template>
  <div class="dashboard-container">
    <el-container>
      <!-- 侧边栏 -->
      <el-aside width="200px">
        <div class="logo">{{ systemName }}</div>
        <el-menu default-active="1" class="sidebar-menu" background-color="#304156" text-color="#bfcbd9" active-text-color="#409EFF">
          <el-menu-item index="1">
            <el-icon><House /></el-icon>
            <span>首页</span>
          </el-menu-item>
          <el-menu-item index="2">
            <el-icon><User /></el-icon>
            <span>用户管理</span>
          </el-menu-item>
          <el-menu-item index="3">
            <el-icon><Setting /></el-icon>
            <span>系统设置</span>
          </el-menu-item>
        </el-menu>
      </el-aside>

      <el-container>
        <!-- 头部 -->
        <el-header>
          <div class="header-content">
            <div class="header-title">后台管理系统</div>
            <div class="header-right">
              <el-dropdown>
                <span class="user-info">
                  <el-icon><User /></el-icon>
                  管理员
                  <el-icon class="el-icon--right"><ArrowDown /></el-icon>
                </span>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item @click="handleLogout">退出登录</el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>
          </div>
        </el-header>

        <!-- 主体内容 -->
        <el-main>
          <div class="welcome-card">
            <h2>欢迎使用{{ systemName }}后台管理系统</h2>
            <p>这是后台管理系统的首页，展示了全局系统设置的使用</p>

            <el-card style="margin-top: 20px;">
              <h3>系统全局配置</h3>
              <p>系统名称：{{ systemName }}</p>
              <p>验证码状态：{{ needCaptcha ? '已开启' : '已关闭' }}</p>
              <p>所有配置：{{ JSON.stringify(systemSettings, null, 2) }}</p>
            </el-card>
          </div>
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { systemName, needCaptcha, systemSettings } from '../stores/settings'
import { House, User, Setting, ArrowDown } from '@element-plus/icons-vue'

const router = useRouter()

// 退出登录
const handleLogout = () => {
  localStorage.removeItem('token')
  router.push('/login')
}
</script>

<style scoped>
.dashboard-container {
  height: 100vh;
}

.logo {
  height: 60px;
  line-height: 60px;
  text-align: center;
  color: #fff;
  font-size: 20px;
  font-weight: bold;
  background-color: #2b2f3a;
}

.sidebar-menu {
  border-right: none;
  height: calc(100vh - 60px);
}

.el-header {
  background-color: #fff;
  border-bottom: 1px solid #e6e6e6;
  padding: 0 20px;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 60px;
}

.header-title {
  font-size: 18px;
  font-weight: bold;
}

.user-info {
  display: flex;
  align-items: center;
  cursor: pointer;
  color: #606266;
}

.user-info:hover {
  color: #409eff;
}

.welcome-card {
  padding: 20px;
}

.welcome-card h2 {
  margin-bottom: 10px;
  color: #303133;
}

.welcome-card p {
  color: #606266;
  margin-bottom: 10px;
}
</style>
