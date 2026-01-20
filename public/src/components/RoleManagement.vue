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

    <!-- 新增角色弹框 -->
    <el-drawer
      v-model="addDialogVisible"
      :title="t('role.addRole')"
      direction="rtl"
      size="500px"
    >
      <el-form
        ref="addRoleFormRef"
        :model="addRoleForm"
        :rules="addRoleRules"
        label-width="100px"
      >
        <el-form-item :label="t('role.name')" prop="name">
          <el-input
            v-model="addRoleForm.name"
            :placeholder="t('role.namePlaceholder')"
            maxlength="50"
            show-word-limit
          />
        </el-form-item>
        <el-form-item :label="t('role.description')" prop="description">
          <el-input
            v-model="addRoleForm.description"
            type="textarea"
            :placeholder="t('role.descPlaceholder')"
            :rows="4"
            maxlength="200"
            show-word-limit
          />
        </el-form-item>
        <el-form-item :label="t('role.status')" prop="status">
          <el-radio-group v-model="addRoleForm.status">
            <el-radio :label="1">{{ t('role.enabled') }}</el-radio>
            <el-radio :label="2">{{ t('role.disabled') }}</el-radio>
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

    <!-- 编辑角色弹框 -->
    <el-drawer
      v-model="editDialogVisible"
      :title="t('role.editRole')"
      direction="rtl"
      size="500px"
    >
      <el-form
        ref="editRoleFormRef"
        :model="editRoleForm"
        :rules="editRoleRules"
        label-width="100px"
      >
        <el-form-item :label="t('role.name')" prop="name">
          <el-input
            v-model="editRoleForm.name"
            :placeholder="t('role.namePlaceholder')"
            maxlength="50"
            show-word-limit
          />
        </el-form-item>
        <el-form-item :label="t('role.description')" prop="description">
          <el-input
            v-model="editRoleForm.description"
            type="textarea"
            :placeholder="t('role.descPlaceholder')"
            :rows="4"
            maxlength="200"
            show-word-limit
          />
        </el-form-item>
        <el-form-item :label="t('role.status')" prop="status">
          <el-radio-group v-model="editRoleForm.status">
            <el-radio :label="1">{{ t('role.enabled') }}</el-radio>
            <el-radio :label="2">{{ t('role.disabled') }}</el-radio>
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

    <!-- 权限设置弹框 -->
    <el-drawer
      v-model="permissionDialogVisible"
      :title="t('role.setPermissionsTitle', { name: currentRole?.name })"
      direction="rtl"
      size="600px"
    >
      <div v-loading="permissionLoading">
        <el-alert
          v-if="currentRole?.system_flag === 1"
          :title="t('role.systemRolePermissionTip')"
          type="warning"
          :closable="false"
          style="margin-bottom: 20px"
        />

        <div v-for="module in allPermissions" :key="module.module" class="permission-module">
          <h4 class="module-title">{{ module.module }}</h4>
          <el-checkbox-group
            v-model="selectedPermissions"
            :disabled="currentRole?.system_flag === 1"
          >
            <el-checkbox
              v-for="perm in module.permissions"
              :key="perm.code"
              :label="perm.code"
              class="permission-checkbox"
            >
              <span>{{ perm.name }}</span>
              <el-tooltip
                v-if="perm.description"
                :content="perm.description"
                placement="top"
              >
                <el-icon style="margin-left: 4px; color: #909399;">
                  <InfoFilled />
                </el-icon>
              </el-tooltip>
            </el-checkbox>
          </el-checkbox-group>
        </div>
      </div>

      <template #footer>
        <div style="flex: auto">
          <el-button @click="permissionDialogVisible = false">{{ t('common.cancel') }}</el-button>
          <el-button
            type="primary"
            :loading="submitLoading"
            :disabled="currentRole?.system_flag === 1"
            @click="handleSubmitPermissions"
          >
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
import { InfoFilled } from '@element-plus/icons-vue'

const { t, locale } = useI18n()

// 获取所有权限列表
const fetchAllPermissions = async () => {
  try {
    const token = localStorage.getItem('token')
    if (!token) {
      ElMessage.error(t('login.loginFailed'))
      return
    }

    const response = await fetch('/api/admin/v1/role/permissions/all', {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Accept-Language': locale.value
      }
    })

    const data = await response.json()

    if (response.ok && data.code === 200) {
      allPermissions.value = data.data || []
    } else {
      ElMessage.error(data.message || t('common.failed'))
    }
  } catch (error) {
    console.error('获取权限列表失败:', error)
    ElMessage.error(t('common.failed'))
  }
}

// 获取角色的权限列表
const fetchRolePermissions = async (roleCode) => {
  try {
    const token = localStorage.getItem('token')
    if (!token) {
      ElMessage.error(t('login.loginFailed'))
      return
    }

    const response = await fetch(`/api/admin/v1/role/permissions/get?code=${roleCode}`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Accept-Language': locale.value
      }
    })

    const data = await response.json()

    if (response.ok && data.code === 200) {
      selectedPermissions.value = data.data || []
    } else {
      ElMessage.error(data.message || t('common.failed'))
    }
  } catch (error) {
    console.error('获取角色权限失败:', error)
    ElMessage.error(t('common.failed'))
  }
}

const roles = ref([])
const loading = ref(false)
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

// 添加角色对话框
const addDialogVisible = ref(false)
const submitLoading = ref(false)
const addRoleForm = ref({
  name: '',
  description: '',
  status: 1
})
const addRoleFormRef = ref(null)

// 编辑角色对话框
const editDialogVisible = ref(false)
const editRoleForm = ref({
  id: 0,
  name: '',
  description: '',
  status: 1
})
const editRoleFormRef = ref(null)

// 权限设置弹框
const permissionDialogVisible = ref(false)
const permissionLoading = ref(false)
const currentRole = ref(null)
const allPermissions = ref([])
const selectedPermissions = ref([])

// 表单验证规则
const addRoleRules = computed(() => ({
  name: [
    { required: true, message: t('role.nameRequired'), trigger: 'blur' },
    { min: 2, max: 50, message: t('role.nameLength'), trigger: 'blur' }
  ],
  description: [
    { max: 200, message: t('role.descLength'), trigger: 'blur' }
  ]
}))

const editRoleRules = computed(() => ({
  name: [
    { required: true, message: t('role.nameRequired'), trigger: 'blur' },
    { min: 2, max: 50, message: t('role.nameLength'), trigger: 'blur' }
  ],
  description: [
    { max: 200, message: t('role.descLength'), trigger: 'blur' }
  ]
}))

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
  // 重置表单
  addRoleForm.value = {
    name: '',
    description: '',
    status: 1
  }
  // 清除表单验证状态
  if (addRoleFormRef.value) {
    addRoleFormRef.value.clearValidate()
  }
  // 打开弹框
  addDialogVisible.value = true
}

// 处理取消添加
const handleCancelAdd = () => {
  addDialogVisible.value = false
  // 清除表单验证状态
  if (addRoleFormRef.value) {
    addRoleFormRef.value.clearValidate()
  }
}

// 处理提交添加
const handleSubmitAdd = async () => {
  if (!addRoleFormRef.value) return

  await addRoleFormRef.value.validate(async (valid) => {
    if (valid) {
      submitLoading.value = true
      try {
        const token = localStorage.getItem('token')
        if (!token) {
          ElMessage.error(t('login.loginFailed'))
          return
        }

        const response = await fetch('/api/admin/v1/role/create', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`,
            'Accept-Language': locale.value
          },
          body: JSON.stringify(addRoleForm.value)
        })

        const data = await response.json()

        if (response.ok && data.code === 200) {
          ElMessage.success(t('role.addSuccess'))
          addDialogVisible.value = false
          // 刷新列表
          fetchRoles()
        } else {
          ElMessage.error(data.message || t('common.failed'))
        }
      } catch (error) {
        console.error('添加角色失败:', error)
        ElMessage.error(t('common.failed'))
      } finally {
        submitLoading.value = false
      }
    }
  })
}

// 处理设置权限
const handleSetPermissions = async (role) => {
  currentRole.value = role
  permissionLoading.value = true
  permissionDialogVisible.value = true

  try {
    // 获取所有权限列表
    await fetchAllPermissions()

    // 获取角色当前的权限
    await fetchRolePermissions(role.code)
  } catch (error) {
    console.error('加载权限数据失败:', error)
  } finally {
    permissionLoading.value = false
  }
}

// 处理提交权限设置
const handleSubmitPermissions = async () => {
  submitLoading.value = true
  try {
    const token = localStorage.getItem('token')
    if (!token) {
      ElMessage.error(t('login.loginFailed'))
      return
    }

    const response = await fetch('/api/admin/v1/role/permissions/assign', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`,
        'Accept-Language': locale.value
      },
      body: JSON.stringify({
        id: currentRole.value.id,
        permissions: selectedPermissions.value
      })
    })

    const data = await response.json()

    if (response.ok && data.code === 200) {
      ElMessage.success(t('role.setPermissionsSuccess'))
      permissionDialogVisible.value = false
      // 刷新角色列表以更新权限显示
      fetchRoles()
    } else {
      ElMessage.error(data.message || t('common.failed'))
    }
  } catch (error) {
    console.error('设置权限失败:', error)
    ElMessage.error(t('common.failed'))
  } finally {
    submitLoading.value = false
  }
}

// 处理编辑角色
const handleEditRole = (role) => {
  // 填充表单数据
  editRoleForm.value = {
    id: role.id,
    name: role.name,
    description: role.description || '',
    status: role.status
  }
  // 清除表单验证状态
  if (editRoleFormRef.value) {
    editRoleFormRef.value.clearValidate()
  }
  // 打开弹框
  editDialogVisible.value = true
}

// 处理取消编辑
const handleCancelEdit = () => {
  editDialogVisible.value = false
  // 清除表单验证状态
  if (editRoleFormRef.value) {
    editRoleFormRef.value.clearValidate()
  }
}

// 处理提交编辑
const handleSubmitEdit = async () => {
  if (!editRoleFormRef.value) return

  await editRoleFormRef.value.validate(async (valid) => {
    if (valid) {
      submitLoading.value = true
      try {
        const token = localStorage.getItem('token')
        if (!token) {
          ElMessage.error(t('login.loginFailed'))
          return
        }

        const response = await fetch('/api/admin/v1/role/update', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`,
            'Accept-Language': locale.value
          },
          body: JSON.stringify(editRoleForm.value)
        })

        const data = await response.json()

        if (response.ok && data.code === 200) {
          ElMessage.success(t('role.editSuccess'))
          editDialogVisible.value = false
          // 刷新列表
          fetchRoles()
        } else {
          ElMessage.error(data.message || t('common.failed'))
        }
      } catch (error) {
        console.error('编辑角色失败:', error)
        ElMessage.error(t('common.failed'))
      } finally {
        submitLoading.value = false
      }
    }
  })
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

.permission-module {
  margin-bottom: 24px;
  padding: 16px;
  border: 1px solid #ebeef5;
  border-radius: 4px;
  background-color: #fff;
}

.permission-module:last-child {
  margin-bottom: 0;
}

.module-title {
  margin: 0 0 16px;
  font-size: 16px;
  font-weight: 500;
  color: #303133;
}

.permission-checkbox {
  margin: 8px 24px 8px 0;
  display: inline-flex;
  align-items: center;
}

.permission-checkbox :deep(.el-checkbox__input) {
  align-self: flex-start;
  margin-top: 3px;
}

.permission-checkbox :deep(.el-checkbox__label) {
  display: inline-flex;
  align-items: center;
}
</style>
