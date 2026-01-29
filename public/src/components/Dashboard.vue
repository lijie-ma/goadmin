<template>
  <div class="dashboard">
    <h1>{{ t('dashboardPage.title') }}</h1>

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
                <div class="stat-label">{{ t('dashboardPage.totalUsers') }}</div>
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
                <div class="stat-label">{{ t('dashboardPage.pageViews') }}</div>
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
                <div class="stat-label">{{ t('dashboardPage.orderCount') }}</div>
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
                <div class="stat-label">{{ t('dashboardPage.revenue') }}</div>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- 最近活动表格 -->
    <el-card class="activity-card" v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>{{ t('dashboardPage.recentActivities') }}</span>
        </div>
      </template>
      <el-table :data="recentActivities" style="width: 100%">
        <el-table-column prop="ctime" :label="t('dashboardPage.time')" width="180" />
        <el-table-column prop="username" :label="t('dashboardPage.user')" width="120" />
        <el-table-column prop="content" :label="t('dashboardPage.action')" />
        <el-table-column prop="ip" :label="t('dashboardPage.ip')" width="140" />
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import {
  User,
  View,
  ShoppingCart,
  Coin
} from '@element-plus/icons-vue'
import { useI18n } from 'vue-i18n'
import axios from 'axios'

const { t } = useI18n()
const loading = ref(false)
const recentActivities = ref([])

// 加载最近活动
const loadRecentActivities = async () => {
  loading.value = true
  try {
    const response = await axios.get('/api/admin/v1/operate_log/list', {
      params: {
        page: 1,
        page_size: 5,
        order_by: "id desc"
      },
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })
    if (response.data.code === 200) {
      recentActivities.value = response.data.data.list || []
    } else {
      console.error('加载最近活动失败:', response.data.message)
    }
  } catch (error) {
    console.error('加载最近活动失败:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadRecentActivities()
})
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

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
