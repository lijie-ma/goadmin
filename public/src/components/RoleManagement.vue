<template>
  <div class="role-management">
    <h1>{{ t('role.management') }}</h1>
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('role.list') }}</span>
          <el-button type="primary">{{ t('role.addRole') }}</el-button>
        </div>
      </template>
      <el-table :data="roles" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" :label="t('role.name')" />
        <el-table-column prop="description" :label="t('role.description')" />
        <el-table-column prop="permissions" :label="t('role.permissionCount')" width="100">
          <template #default="scope">
            {{ scope.row.permissions }}{{ t('role.count') }}
          </template>
        </el-table-column>
        <el-table-column prop="userCount" :label="t('role.userCount')" width="100">
          <template #default="scope">
            {{ scope.row.userCount }}{{ t('role.person') }}
          </template>
        </el-table-column>
        <el-table-column prop="status" :label="t('role.status')" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.status === 'active' ? 'success' : 'info'">
              {{ scope.row.status === 'active' ? t('role.enabled') : t('role.disabled') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('role.actions')" width="200">
          <template #default>
            <el-button link type="primary" size="small">{{ t('role.setPermissions') }}</el-button>
            <el-button link type="primary" size="small">{{ t('role.edit') }}</el-button>
            <el-button link type="danger" size="small">{{ t('role.delete') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const roles = ref([
  {
    id: 1,
    name: '超级管理员',
    description: '拥有系统所有权限',
    permissions: 50,
    userCount: 1,
    status: 'active'
  },
  {
    id: 2,
    name: '管理员',
    description: '拥有大部分管理权限',
    permissions: 35,
    userCount: 5,
    status: 'active'
  },
  {
    id: 3,
    name: '普通用户',
    description: '基本用户权限',
    permissions: 10,
    userCount: 100,
    status: 'active'
  },
  {
    id: 4,
    name: '访客',
    description: '只读权限',
    permissions: 5,
    userCount: 50,
    status: 'inactive'
  }
])
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
</style>
