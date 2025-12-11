<template>
  <div class="user-management">
    <h1>用户管理</h1>
    <el-card>
      <template #header>
        <div class="card-header">
          <span>用户列表</span>
          <el-button type="primary">新增用户</el-button>
        </div>
      </template>
      <el-table :data="users" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" label="用户名" />
        <el-table-column prop="email" label="邮箱" />
        <el-table-column prop="role" label="角色" />
        <el-table-column prop="status" label="状态">
          <template #default="scope">
            <el-tag :type="scope.row.status === 'active' ? 'success' : 'danger'">
              {{ scope.row.status === 'active' ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180">
          <template #default>
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

const users = ref([
  {
    id: 1,
    username: 'admin',
    email: 'admin@example.com',
    role: '超级管理员',
    status: 'active'
  },
  {
    id: 2,
    username: 'user1',
    email: 'user1@example.com',
    role: '普通用户',
    status: 'active'
  },
  {
    id: 3,
    username: 'user2',
    email: 'user2@example.com',
    role: '普通用户',
    status: 'inactive'
  }
])
</script>

<style scoped>
.user-management {
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
