<template>
  <div class="role-management">
    <h1>{{ t('role.management') }}</h1>
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('role.list') }}</span>
          <el-button type="primary" @click="handleAddRole">{{ t('role.addRole') }}</el-button>
        </div>
      </template>
      <el-table :data="roles" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" :label="t('role.name')" />
        <el-table-column prop="description" :label="t('role.description')" />
        <el-table-column prop="permissionList" :label="t('role.permissionList')" min-width="200">
          <template #default="scope">
            <el-tooltip
              v-if="scope.row.permissionList && scope.row.permissionList.length > 0"
              :content="scope.row.permissionList.map(p => p.name).join(', ')"
              placement="top"
              :disabled="scope.row.permissionList.length <= 3"
            >
              <span class="permission-list">
                {{ scope.row.permissionList.slice(0, 3).map(p => p.name).join(', ') }}
                <span v-if="scope.row.permissionList.length > 3">...</span>
              </span>
            </el-tooltip>
            <span v-else class="no-permissions">{{ t('role.noPermissions') }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" :label="t('role.status')" width="180">
          <template #default="scope">
            <el-tag :type="scope.row.status === 1 ? 'success' : 'info'" class="status-tag">
              {{ scope.row.status === 1 ? t('role.enabled') : t('role.disabled') }}
            </el-tag>
            <el-tag v-if="scope.row.system_flag === 1" type="warning" class="system-tag">
              {{ t('role.systemRole') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('role.actions')" width="200">
          <template #default="scope">
            <el-button link type="primary" size="small" @click="handleSetPermissions(scope.row)">{{ t('role.setPermissions') }}</el-button>
            <el-button link type="primary" size="small" @click="handleEditRole(scope.row)">{{ t('role.edit') }}</el-button>
            <el-button
              v-if="scope.row.system_flag !== 1"
              link
              type="danger"
              size="small"
              @click="handleDeleteRole(scope.row)"
            >
              {{ t('role.delete') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      <!-- 分页 -->
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[10, 20, 30, 50]"
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
import { ElMessage, ElMessageBox } from 'element-plus'

const { t, locale } = useI18n()

const roles = ref([])
const loading = ref(false)
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

// 获取角色列表
const fetchRoles = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    if (!token) {
      ElMessage.error(t('login.loginFailed'))
      return
    }

    const response = await fetch(`/api/admin/v1/role/list?page=${currentPage.value}&page_size=${pageSize.value}&order_by=id desc`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Accept-Language': locale.value
      }
    })

    const data = await response.json()

    if (response.ok && data.code === 200) {
      roles.value = data.data.list || []
      total.value = data.data.total || 0

      // 添加计算字段
      roles.value = roles.value.map(role => ({
        ...role,
        permissionList: role.permissions ?? [],
      }))
    } else {
      ElMessage.error(data.message || t('settings.loadFailed'))
    }
  } catch (error) {
    console.error('获取角色列表失败:', error)
    ElMessage.error(t('settings.loadFailed'))
  } finally {
    loading.value = false
  }
}

// 组件挂载时获取数据
onMounted(() => {
  fetchRoles()
})

// 处理分页变化
const handleSizeChange = (val) => {
  pageSize.value = val
  fetchRoles()
}

const handleCurrentChange = (val) => {
  currentPage.value = val
  fetchRoles()
}

// 处理添加角色
const handleAddRole = () => {
  // TODO: 实现添加角色功能
}

// 处理设置权限
const handleSetPermissions = (role) => {
  // TODO: 实现权限设置功能
}

// 处理编辑角色
const handleEditRole = (role) => {
  // TODO: 实现编辑角色功能
}

// 处理删除角色
const handleDeleteRole = async (role) => {
  try {
    await ElMessageBox.confirm(
      t('role.deleteConfirm'),
      t('common.confirm'),
      {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        type: 'warning'
      }
    )

    const token = localStorage.getItem('token')
    if (!token) {
      ElMessage.error(t('login.loginFailed'))
      return
    }

    const response = await fetch('/api/admin/v1/role/delete', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`,
        'Accept-Language': locale.value
      },
      body: JSON.stringify({
        id: role.id
      })
    })

    const data = await response.json()

    if (response.ok && data.code === 200) {
      ElMessage.success(t('role.deleteSuccess'))

      // 删除最后一页的最后一条数据时，跳转到上一页
      if (roles.value.length === 1 && currentPage.value > 1) {
        currentPage.value--
      }
      fetchRoles()
    } else {
      ElMessage.error(data.message || t('common.failed'))
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除角色失败:', error)
      ElMessage.error(t('common.failed'))
    }
  }
}
</script>

<style scoped>
.role-management {
  padding: 20px;
}

h1 {
  margin-bottom: 24px;
  font-weight: 500;
  font-size: 24px;
  color: #303133;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.permission-list {
  display: inline-block;
  max-width: 100%;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.no-permissions {
  color: #909399;
  font-style: italic;
}

.status-tag {
  margin-right: 8px;
}

.system-tag {
  font-size: 12px;
}
</style>
