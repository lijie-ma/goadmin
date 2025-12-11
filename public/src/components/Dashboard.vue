<template>
  <div class="dashboard">
    <h1>仪表板</h1>

    <!-- 统计卡片 -->
    <div class="stats-cards">
      <el-row :gutter="20">
        <el-col :xs="24" :sm="12" :lg="6">
          <el-card class="stat-card">
            <div class="stat-content">
              <div class="stat-icon users">
                <el-icon :size="24"><User /></el-icon>
              </div>
              <div class="stat-data">
                <div class="stat-value">1,234</div>
                <div class="stat-label">用户总数</div>
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :xs="24" :sm="12" :lg="6">
          <el-card class="stat-card">
            <div class="stat-content">
              <div class="stat-icon views">
                <el-icon :size="24"><View /></el-icon>
              </div>
              <div class="stat-data">
                <div class="stat-value">23,456</div>
                <div class="stat-label">页面访问量</div>
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :xs="24" :sm="12" :lg="6">
          <el-card class="stat-card">
            <div class="stat-content">
              <div class="stat-icon orders">
                <el-icon :size="24"><ShoppingCart /></el-icon>
              </div>
              <div class="stat-data">
                <div class="stat-value">345</div>
                <div class="stat-label">订单数量</div>
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :xs="24" :sm="12" :lg="6">
          <el-card class="stat-card">
            <div class="stat-content">
              <div class="stat-icon revenue">
                <el-icon :size="24"><Coin /></el-icon>
              </div>
              <div class="stat-data">
                <div class="stat-value">¥234,567</div>
                <div class="stat-label">营业收入</div>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- 最近活动表格 -->
    <el-card class="activity-card">
      <template #header>
        <div class="card-header">
          <span>最近活动</span>
        </div>
      </template>
      <el-table :data="recentActivities" style="width: 100%">
        <el-table-column prop="time" label="时间" width="180" />
        <el-table-column prop="user" label="用户" width="120" />
        <el-table-column prop="action" label="操作" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.status === 'success' ? 'success' : 'danger'">
              {{ scope.row.status === 'success' ? '成功' : '失败' }}
            </el-tag>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 快捷操作 -->
    <el-card class="quick-actions">
      <template #header>
        <div class="card-header">
          <span>快捷操作</span>
        </div>
      </template>
      <el-row :gutter="20">
        <el-col :span="6">
          <el-button type="primary" @click="handleQuickAction('users')">
            <el-icon><User /></el-icon>
            用户管理
          </el-button>
        </el-col>
        <el-col :span="6">
          <el-button type="success" @click="handleQuickAction('roles')">
            <el-icon><Lock /></el-icon>
            角色管理
          </el-button>
        </el-col>
        <el-col :span="6">
          <el-button type="warning" @click="handleQuickAction('settings')">
            <el-icon><Setting /></el-icon>
            系统设置
          </el-button>
        </el-col>
        <el-col :span="6">
          <el-button type="info" @click="handleQuickAction('help')">
            <el-icon><QuestionFilled /></el-icon>
            帮助中心
          </el-button>
        </el-col>
      </el-row>
    </el-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import {
  User,
  View,
  ShoppingCart,
  Coin,
  Lock,
  Setting,
  QuestionFilled
} from '@element-plus/icons-vue'

const router = useRouter()

const recentActivities = ref([
  {
    time: '2024-12-11 16:30:00',
    user: 'Admin',
    action: '登录系统',
    status: 'success'
  },
  {
    time: '2024-12-11 16:25:00',
    user: 'User001',
    action: '修改密码',
    status: 'success'
  },
  {
    time: '2024-12-11 16:20:00',
    user: 'User002',
    action: '上传文件',
    status: 'success'
  },
  {
    time: '2024-12-11 16:15:00',
    user: 'User003',
    action: '删除数据',
    status: 'failed'
  }
])

const handleQuickAction = (action) => {
  const routes = {
    users: '/users',
    roles: '/roles',
    settings: '/settings',
    help: '/help'
  }
  if (routes[action]) {
    router.push(routes[action])
  }
}
</script>

<style scoped>
.dashboard {
  padding: 20px;
}

h1 {
  margin-bottom: 24px;
  font-weight: 500;
  font-size: 24px;
  color: #303133;
}

.stats-cards {
  margin-bottom: 24px;
}

.stat-card {
  margin-bottom: 20px;
}

.stat-content {
  display: flex;
  align-items: center;
  padding: 20px;
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 16px;
  color: white;
}

.stat-icon.users {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.stat-icon.views {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.stat-icon.orders {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
}

.stat-icon.revenue {
  background: linear-gradient(135deg, #fa709a 0%, #fee140 100%);
}

.stat-value {
  font-size: 24px;
  font-weight: bold;
  color: #303133;
  line-height: 1;
  margin-bottom: 8px;
}

.stat-label {
  font-size: 14px;
  color: #909399;
}

.activity-card {
  margin-bottom: 24px;
}

.quick-actions {
  margin-bottom: 24px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.el-button {
  width: 100%;
  margin-bottom: 10px;
}

.el-button .el-icon {
  margin-right: 8px;
}
</style>
