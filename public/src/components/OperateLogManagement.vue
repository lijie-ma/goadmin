<template>
  <div class="operate-log-management">
    <h1>{{ t('operateLog.title') }}</h1>
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('operateLog.title') }}</span>
        </div>
      </template>

      <!-- 搜索栏 -->
      <div class="search-bar">
        <el-input
          v-model="searchForm.username"
          :placeholder="t('operateLog.username')"
          clearable
          style="width: 200px; margin-right: 10px;"
        >
          <template #prefix>
            <el-icon><User /></el-icon>
          </template>
        </el-input>
        <el-input
          v-model="searchForm.content"
          :placeholder="t('operateLog.content')"
          clearable
          style="width: 300px; margin-right: 10px;"
        >
          <template #prefix>
            <el-icon><Document /></el-icon>
          </template>
        </el-input>
        <el-input
          v-model="searchForm.ip"
          :placeholder="t('operateLog.ip')"
          clearable
          style="width: 200px; margin-right: 10px;"
        >
          <template #prefix>
            <el-icon><Location /></el-icon>
          </template>
        </el-input>
        <el-date-picker
          v-model="dateRange"
          type="datetimerange"
          :range-separator="t('common.to')"
          :start-placeholder="t('operateLog.startTime')"
          :end-placeholder="t('operateLog.endTime')"
          style="width: 360px; margin-right: 10px;"
          format="YYYY-MM-DD HH:mm:ss"
          value-format="YYYY-MM-DD HH:mm:ss"
        />
        <el-button type="primary" @click="handleSearch">
          <el-icon><Search /></el-icon>
          {{ t('common.search') }}
        </el-button>
        <el-button @click="handleReset">
          <el-icon><Refresh /></el-icon>
          {{ t('common.reset') }}
        </el-button>
      </div>

      <el-table :data="logs" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" :label="t('operateLog.username')" width="120">
          <template #default="scope">
            <el-tag type="primary" size="small">{{ scope.row.username || '-' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="content" :label="t('operateLog.content')" min-width="300">
          <template #default="scope">
            <div class="log-content">{{ scope.row.content }}</div>
          </template>
        </el-table-column>
        <el-table-column prop="ip" :label="t('operateLog.ip')" width="150">
          <template #default="scope">
            <el-tag type="info" size="small">{{ scope.row.ip }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="ctime" :label="t('operateLog.createTime')" width="180">
          <template #default="scope">
            {{ formatDate(scope.row.ctime) }}
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[10, 20, 30, 50, 100]"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
        style="margin-top: 20px; justify-content: flex-end;"
      />
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { Search, Refresh, User, Document, Location } from '@element-plus/icons-vue'
import axios from 'axios'

const { t, locale } = useI18n()

// 响应式数据
const logs = ref([])
const loading = ref(false)
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const dateRange = ref([])

// 搜索表单
const searchForm = ref({
  username: '',
  content: '',
  ip: '',
  start_time: '',
  end_time: ''
})

// 获取操作日志列表
const fetchLogs = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    if (!token) {
      ElMessage.error(t('login.loginFailed'))
      return
    }

    // 处理日期范围
    let startTime = ''
    let endTime = ''
    if (dateRange.value && dateRange.value.length === 2) {
      startTime = dateRange.value[0]
      endTime = dateRange.value[1]
    }

    const response = await axios.get('/api/admin/v1/operate_log/list', {
      headers: {
        'Authorization': `Bearer ${token}`,
        'Accept-Language': locale.value
      },
      params: {
        page: currentPage.value,
        page_size: pageSize.value,
        order_by: 'id desc',
        username: searchForm.value.username,
        content: searchForm.value.content,
        ip: searchForm.value.ip,
        start_time: startTime,
        end_time: endTime
      }
    })

    if (response.data.code === 200) {
      logs.value = response.data.data.list || []
      total.value = response.data.data.total || 0
    } else {
      ElMessage.error(response.data.message || t('common.failed'))
    }
  } catch (error) {
    console.error('获取操作日志列表失败:', error)
    ElMessage.error(t('common.error.systemError'))
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  currentPage.value = 1
  fetchLogs()
}

// 重置
const handleReset = () => {
  searchForm.value = {
    username: '',
    content: '',
    ip: '',
    start_time: '',
    end_time: ''
  }
  dateRange.value = []
  currentPage.value = 1
  fetchLogs()
}

// 分页大小改变
const handleSizeChange = (val) => {
  pageSize.value = val
  currentPage.value = 1
  fetchLogs()
}

// 当前页改变
const handleCurrentChange = (val) => {
  currentPage.value = val
  fetchLogs()
}

// 格式化日期
const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  const seconds = String(date.getSeconds()).padStart(2, '0')
  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
}

// 组件挂载时获取数据
onMounted(() => {
  fetchLogs()
})
</script>

<style scoped>
.operate-log-management {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.search-bar {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  margin-bottom: 20px;
  align-items: center;
}

.log-content {
  word-break: break-all;
  line-height: 1.5;
}
</style>