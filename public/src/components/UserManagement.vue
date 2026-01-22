<template>
  <div class="user-management">
    <h1>{{ t('user.title') }}</h1>
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('user.title') }}</span>
          <el-button type="primary" @click="handleAddUser">{{ t('user.add') }}</el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <div class="search-bar">
        <el-input
          v-model="searchKeyword"
          :placeholder="t('user.search')"
          clearable
          @clear="handleSearch"
          @keyup.enter="handleSearch"
          style="width: 300px; margin-right: 10px;"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
        <el-button type="primary" @click="handleSearch">{{ t('common.search') }}</el-button>
        <el-button @click="handleReset">{{ t('common.reset') }}</el-button>
      </div>

      <el-table :data="users" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" :label="t('user.username')" />
        <el-table-column prop="email" :label="t('user.email')" />
        <el-table-column prop="role" :label="t('user.role')">
          <template #default="scope">
            <el-tag v-if="scope.row.role" type="primary">{{ scope.row.role.name }}</el-tag>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" :label="t('user.status')" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.status === 1 ? 'success' : 'danger'">
              {{ scope.row.status === 1 ? t('user.active') : t('user.inactive') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="ctime" :label="t('user.createTime')" width="180">
          <template #default="scope">
            {{ formatDate(scope.row.ctime) }}
          </template>
        </el-table-column>
        <el-table-column :label="t('user.operations')" width="200" fixed="right">
          <template #default="scope">
            <el-button link type="primary" size="small" @click="handleEditUser(scope.row)">{{ t('common.edit') }}</el-button>
            <el-button link type="primary" size="small" @click="handleResetPassword(scope.row)">{{ t('user.resetPassword') }}</el-button>
            <el-button
              v-if="scope.row.role_code !== 'sup_admin'"
              link
              type="danger"
              size="small"
              @click="handleDeleteUser(scope.row)"
            >
              {{ t('common.delete') }}
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

    <!-- 新增用户弹框 -->
    <el-drawer
      v-model="addDialogVisible"
      :title="t('user.add')"
      direction="rtl"
      size="500px"
    >
      <el-form
        ref="addUserFormRef"
        :model="addUserForm"
        :rules="addUserRules"
        label-width="100px"
      >
        <el-form-item :label="t('user.username')" prop="username">
          <el-input
            v-model="addUserForm.username"
            :placeholder="t('user.username')"
            maxlength="50"
            show-word-limit
          />
        </el-form-item>
        <el-form-item :label="t('user.password')" prop="password">
          <el-input
            v-model="addUserForm.password"
            type="password"
            :placeholder="t('user.newPassword')"
            maxlength="20"
            show-password
          />
        </el-form-item>
        <el-form-item :label="t('user.email')" prop="email">
          <el-input
            v-model="addUserForm.email"
            :placeholder="t('user.email')"
            maxlength="100"
          />
        </el-form-item>
        <el-form-item :label="t('user.role')" prop="role_code">
          <el-select
            v-model="addUserForm.role_code"
            :placeholder="t('user.role')"
            style="width: 100%"
          >
            <el-option
              v-for="role in roles"
              :key="role.code"
              :label="role.name"
              :value="role.code"
              :disabled="role.code === 'sup_admin'"
            />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('user.status')" prop="status">
          <el-radio-group v-model="addUserForm.status">
            <el-radio :label="1">{{ t('user.active') }}</el-radio>
            <el-radio :label="0">{{ t('user.inactive') }}</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <div style="flex: auto">
          <el-button @click="handleCancelAdd">{{ t('common.cancel') }}</el-button>
          <el-button type="primary" :loading="submitLoading" @click="handleSubmitAdd">
            {{ t('common.confirm') }}
          </el-button>
        </div>
      </template>
    </el-drawer>

    <!-- 编辑用户弹框 -->
    <el-drawer
      v-model="editDialogVisible"
      :title="t('user.edit')"
      direction="rtl"
      size="500px"
    >
      <el-form
        ref="editUserFormRef"
        :model="editUserForm"
        :rules="editUserRules"
        label-width="100px"
      >
        <el-form-item :label="t('user.username')" prop="username">
          <el-input
            v-model="editUserForm.username"
            :placeholder="t('user.username')"
            maxlength="50"
            show-word-limit
          />
        </el-form-item>
        <el-form-item :label="t('user.email')" prop="email">
          <el-input
            v-model="editUserForm.email"
            :placeholder="t('user.email')"
            maxlength="100"
          />
        </el-form-item>
        <el-form-item :label="t('user.role')" prop="role_code">
          <el-select
            v-model="editUserForm.role_code"
            :placeholder="t('user.role')"
            style="width: 100%"
          >
            <el-option
              v-for="role in roles"
              :key="role.code"
              :label="role.name"
              :value="role.code"
            />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('user.status')" prop="status">
          <el-radio-group v-model="editUserForm.status">
            <el-radio :label="1">{{ t('user.active') }}</el-radio>
            <el-radio :label="0">{{ t('user.inactive') }}</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <div style="flex: auto">
          <el-button @click="handleCancelEdit">{{ t('common.cancel') }}</el-button>
          <el-button type="primary" :loading="submitLoading" @click="handleSubmitEdit">
            {{ t('common.confirm') }}
          </el-button>
        </div>
      </template>
    </el-drawer>

  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search } from '@element-plus/icons-vue'
import axios from 'axios'
import md5 from 'js-md5'

const { t, locale } = useI18n()

// 响应式数据
const users = ref([])
const roles = ref([])
const loading = ref(false)
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const searchKeyword = ref('')

// 新增用户对话框
const addDialogVisible = ref(false)
const submitLoading = ref(false)
const addUserForm = ref({
  username: '',
  password: '',
  email: '',
  role_code: '',
  status: 1
})
const addUserFormRef = ref(null)

// 编辑用户对话框
const editDialogVisible = ref(false)
const editUserForm = ref({
  id: 0,
  username: '',
  email: '',
  role_code: '',
  status: 1
})
const editUserFormRef = ref(null)


// 表单验证规则
const addUserRules = computed(() => ({
  username: [
    { required: true, message: t('user.username') + t('common.error.required'), trigger: 'blur' },
    { min: 3, max: 50, message: t('user.username') + '长度在 3 到 50 个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: t('user.newPassword') + t('common.error.required'), trigger: 'blur' },
    { min: 6, max: 20, message: t('user.newPassword') + '长度在 6 到 20 个字符', trigger: 'blur' }
  ],
  email: [
    { type: 'email', message: t('user.email') + t('common.error.invalidFormat'), trigger: 'blur' }
  ],
  role_code: [
    { required: true, message: t('user.role') + t('common.error.required'), trigger: 'change' }
  ]
}))

const editUserRules = computed(() => ({
  username: [
    { required: true, message: t('user.username') + t('common.error.required'), trigger: 'blur' },
    { min: 3, max: 50, message: t('user.username') + '长度在 3 到 50 个字符', trigger: 'blur' }
  ],
  email: [
    { type: 'email', message: t('user.email') + t('common.error.invalidFormat'), trigger: 'blur' }
  ],
  role_code: [
    { required: true, message: t('user.role') + t('common.error.required'), trigger: 'change' }
  ]
}))


// 获取用户列表
const fetchUsers = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    if (!token) {
      ElMessage.error(t('login.loginFailed'))
      return
    }

    const response = await axios.get('/api/admin/v1/user/list', {
      headers: {
        'Authorization': `Bearer ${token}`,
        'Accept-Language': locale.value
      },
      params: {
        page: currentPage.value,
        page_size: pageSize.value,
        order_by: 'id desc',
        keyword: searchKeyword.value
      }
    })

    if (response.data.code === 200) {
      users.value = response.data.data.list || []
      total.value = response.data.data.total || 0
    } else {
      ElMessage.error(response.data.message || t('common.failed'))
    }
  } catch (error) {
    console.error('获取用户列表失败:', error)
    ElMessage.error(t('common.error.systemError'))
  } finally {
    loading.value = false
  }
}

// 获取角色列表
const fetchRoles = async () => {
  try {
    const token = localStorage.getItem('token')
    if (!token) {
      ElMessage.error(t('login.loginFailed'))
      return
    }

    const response = await axios.get('/api/admin/v1/role/all', {
      headers: {
        'Authorization': `Bearer ${token}`,
        'Accept-Language': locale.value
      }
    })

    if (response.data.code === 200) {
      roles.value = response.data.data || []
    } else {
      ElMessage.error(response.data.message || t('common.failed'))
    }
  } catch (error) {
    console.error('获取角色列表失败:', error)
    ElMessage.error(t('common.failed'))
  }
}

// 搜索
const handleSearch = () => {
  currentPage.value = 1
  fetchUsers()
}

// 重置
const handleReset = () => {
  searchKeyword.value = ''
  currentPage.value = 1
  fetchUsers()
}

// 分页大小改变
const handleSizeChange = (val) => {
  pageSize.value = val
  currentPage.value = 1
  fetchUsers()
}

// 当前页改变
const handleCurrentChange = (val) => {
  currentPage.value = val
  fetchUsers()
}

// 新增用户
const handleAddUser = () => {
  // 重置表单
  addUserForm.value = {
    username: '',
    password: '',
    email: '',
    role_code: '',
    status: 1
  }
  // 清除表单验证状态
  if (addUserFormRef.value) {
    addUserFormRef.value.clearValidate()
  }
  // 打开弹框
  addDialogVisible.value = true
}

// 处理取消添加
const handleCancelAdd = () => {
  addDialogVisible.value = false
  // 清除表单验证状态
  if (addUserFormRef.value) {
    addUserFormRef.value.clearValidate()
  }
}

// 处理提交添加
const handleSubmitAdd = async () => {
  if (!addUserFormRef.value) return

  await addUserFormRef.value.validate(async (valid) => {
    if (valid) {
      submitLoading.value = true
      try {
        const token = localStorage.getItem('token')
        if (!token) {
          ElMessage.error(t('login.loginFailed'))
          return
        }

        // MD5加密密码
        const encryptedPassword = md5(addUserForm.value.password)

        const response = await axios.post('/api/admin/v1/user/create', {
          ...addUserForm.value,
          password: encryptedPassword
        }, {
          headers: {
            'Authorization': `Bearer ${token}`,
            'Accept-Language': locale.value
          }
        })

        if (response.data.code === 200) {
          ElMessage.success(t('user.addSuccess'))
          addDialogVisible.value = false
          // 刷新列表
          fetchUsers()
        } else {
          ElMessage.error(response.data.message || t('common.failed'))
        }
      } catch (error) {
        console.error('添加用户失败:', error)
        ElMessage.error(t('common.failed'))
      } finally {
        submitLoading.value = false
      }
    }
  })
}

// 编辑用户
const handleEditUser = (row) => {
  // 填充表单数据
  editUserForm.value = {
    id: row.id,
    username: row.username,
    email: row.email || '',
    role_code: row.role_code,
    status: row.status
  }
  // 清除表单验证状态
  if (editUserFormRef.value) {
    editUserFormRef.value.clearValidate()
  }
  // 打开弹框
  editDialogVisible.value = true
}

// 处理取消编辑
const handleCancelEdit = () => {
  editDialogVisible.value = false
  // 清除表单验证状态
  if (editUserFormRef.value) {
    editUserFormRef.value.clearValidate()
  }
}

// 处理提交编辑
const handleSubmitEdit = async () => {
  if (!editUserFormRef.value) return

  await editUserFormRef.value.validate(async (valid) => {
    if (valid) {
      submitLoading.value = true
      try {
        const token = localStorage.getItem('token')
        if (!token) {
          ElMessage.error(t('login.loginFailed'))
          return
        }

        const response = await axios.post('/api/admin/v1/user/update', editUserForm.value, {
          headers: {
            'Authorization': `Bearer ${token}`,
            'Accept-Language': locale.value
          }
        })

        if (response.data.code === 200) {
          ElMessage.success(t('user.editSuccess'))
          editDialogVisible.value = false
          // 刷新列表
          fetchUsers()
        } else {
          ElMessage.error(response.data.message || t('common.failed'))
        }
      } catch (error) {
        console.error('编辑用户失败:', error)
        ElMessage.error(t('common.failed'))
      } finally {
        submitLoading.value = false
      }
    }
  })
}

// 重置密码
const handleResetPassword = async (row) => {
  try {
    await ElMessageBox.confirm(
      t('user.resetPasswordConfirm'),
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

    submitLoading.value = true
    try {
      const response = await axios.post('/api/admin/v1/user/reset_pwd', {
        id: row.id
      }, {
        headers: {
          'Authorization': `Bearer ${token}`,
          'Accept-Language': locale.value
        }
      })

      if (response.data.code === 200) {
        ElMessage.success(t('user.resetPasswordSuccess'))
      } else {
        ElMessage.error(response.data.message || t('common.failed'))
      }
    } catch (error) {
      console.error('重置密码失败:', error)
      ElMessage.error(t('common.failed'))
    } finally {
      submitLoading.value = false
    }
  } catch (error) {
    // 用户取消操作
  }
}

// 删除用户
const handleDeleteUser = async (row) => {
  try {
    await ElMessageBox.confirm(
      t('user.deleteConfirm'),
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

    const response = await axios.post('/api/admin/v1/user/delete', {
      id: row.id
    }, {
      headers: {
        'Authorization': `Bearer ${token}`,
        'Accept-Language': locale.value
      }
    })

    if (response.data.code === 200) {
      ElMessage.success(t('user.deleteSuccess'))

      // 删除最后一页的最后一条数据时，跳转到上一页
      if (users.value.length === 1 && currentPage.value > 1) {
        currentPage.value--
      }
      fetchUsers()
    } else {
      ElMessage.error(response.data.message || t('common.failed'))
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除用户失败:', error)
      ElMessage.error(t('common.failed'))
    }
  }
}

// 格式化日期
const formatDate = (dateString) => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}

// 组件挂载时获取数据
onMounted(() => {
  fetchUsers()
  fetchRoles()
})
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

.search-bar {
  margin-bottom: 20px;
  display: flex;
  align-items: center;
}

/* 允许表单标签换行 */
:deep(.el-form-item__label) {
  white-space: normal;
  word-break: break-word;
  line-height: 1.5;
}
</style>
