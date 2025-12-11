<template>
  <div class="dashboard-container">
    <!-- 统计卡片 -->
    <el-row :gutter="20" class="stats-row">
      <el-col :xs="24" :sm="12" :md="12" :lg="6">
        <el-card class="stats-card">
          <div class="stats-item">
            <div class="stats-icon" style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);">
              <el-icon :size="30"><User /></el-icon>
            </div>
            <div class="stats-info">
              <div class="stats-value">{{ statsData.users }}</div>
              <div class="stats-label">用户总数</div>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :xs="24" :sm="12" :md="12" :lg="6">
        <el-card class="stats-card">
          <div class="stats-item">
            <div class="stats-icon" style="background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);">
              <el-icon :size="30"><View /></el-icon>
            </div>
            <div class="stats-info">
              <div class="stats-value">{{ statsData.views }}</div>
              <div class="stats-label">今日访问</div>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :xs="24" :sm="12" :md="12" :lg="6">
        <el-card class="stats-card">
          <div class="stats-item">
            <div class="stats-icon" style="background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);">
              <el-icon :size="30"><ShoppingCart /></el-icon>
            </div>
            <div class="stats-info">
              <div class="stats-value">{{ statsData.orders }}</div>
              <div class="stats-label">订单数量</div>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :xs="24" :sm="12" :md="12" :lg="6">
        <el-card class="stats-card">
          <div class="stats-item">
            <div class="stats-icon" style="background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);">
              <el-icon :size="30"><Coin /></el-icon>
            </div>
            <div class="stats-info">
              <div class="stats-value">￥{{ statsData.revenue }}</div>
              <div class="stats-label">营业收入</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 图表区域 -->
    <el-row :gutter="20" class="chart-row">
      <el-col :xs="24" :sm="24" :md="24" :lg="16">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>数据趋势</span>
              <el-radio-group v-model="chartData.period" size="small">
                <el-radio-button label="week">本周</el-radio-button>
                <el-radio-button label="month">本月</el-radio-button>
                <el-radio-button label="year">本年</el-radio-button>
              </el-radio-group>
            </div>
          </template>
          <div class="chart-container">
            <div class="chart-placeholder">
              <el-icon :size="48" color="#409EFF"><DataLine /></el-icon>
              <p>图表数据加载中...</p>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :xs="24" :sm="24" :md="24" :lg="8">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>系统信息</span>
            </div>
          </template>
          <div class="system-info">
            <div v-for="(item, index) in systemInfo" :key="index" class="info-item">
              <span class="label">{{ item.label }}:</span>
              <span class="value">{{ item.value }}</span>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { User, View, ShoppingCart, Coin, DataLine } from '@element-plus/icons-vue'

// 统计数据
const statsData = reactive({
  users: 2458,
  views: 1234,
  orders: 456,
  revenue: '45,678.00'
})

// 图表数据
const chartData = reactive({
  period: 'week',
  // 这里可以添加实际的图表数据
})

// 系统信息
const systemInfo = ref([
  { label: '系统版本', value: 'v1.0.0' },
  { label: '服务器环境', value: 'CentOS 7.9' },
  { label: '数据库版本', value: 'MySQL 8.0' },
  { label: '缓存服务', value: 'Redis 6.2' },
  { label: '运行时间', value: '30天' }
])

onMounted(() => {
  // 这里可以添加初始化逻辑，如获取实际数据等
})
</script>

<style scoped>
.dashboard-container {
  padding: 20px;
}

.stats-row {
  margin-bottom: 20px;
}

.stats-card {
  cursor: pointer;
  transition: transform 0.3s;
  margin-bottom: 20px;
}

@media (max-width: 768px) {
  .stats-card {
    margin-bottom: 25px;
  }
}

.stats-card:hover {
  transform: translateY(-5px);
}

.stats-item {
  display: flex;
  align-items: center;
}

.stats-icon {
  width: 60px;
  height: 60px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  margin-right: 16px;
}

.stats-info {
  flex: 1;
}

.stats-value {
  font-size: 24px;
  font-weight: bold;
  color: #303133;
  line-height: 1.2;
  margin-bottom: 4px;
}

.stats-label {
  font-size: 14px;
  color: #909399;
}

.chart-row {
  margin-bottom: 20px;
}

.chart-row .el-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.chart-container {
  height: 400px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.chart-placeholder {
  text-align: center;
  color: #909399;
}

.chart-placeholder p {
  margin-top: 12px;
}

.system-info {
  .info-item {
    display: flex;
    margin-bottom: 12px;
    line-height: 1.8;

    .label {
      color: #606266;
      margin-right: 8px;
      min-width: 80px;
    }

    .value {
      color: #303133;
      flex: 1;
    }
  }
}
</style>
