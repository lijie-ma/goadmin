<template>
  <div class="role-management">
    <h1>角色管理</h1>
    <el-card>
      <template #header>
        <div class="card-header">
          <span>角色列表</span>
          <el-button type="primary">新增角色</el-button>
        </div>
      </template>
      <el-table :data="roles" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="角色名称" />
        <el-table-column prop="description" label="描述" />
        <el-table-column prop="permissions" label="权限数" width="100">
          <template #default="scope">
            {{ scope.row.permissions }}个
          </template>
        </el-table-column>
        <el-table-column prop="userCount" label="用户数" width="100">
          <template #default="scope">
            {{ scope.row.userCount }}人
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.status === 'active' ? 'success' : 'info'">
              {{ scope.row.status === 'active' ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200">
          <template #default>
            <el-button link type="primary" size="small">权限设置</el-button>
            <el-button link type="primary" size="small">编辑</el-button>
            <el-button link type="danger" size="small">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'

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
