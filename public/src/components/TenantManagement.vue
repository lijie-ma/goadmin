<template>
  <div class="tenant-management">
    <h1>{{ t('tenant.title') }}</h1>
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('tenant.title') }}</span>
          <el-button
            v-permission="'tenant_create'"
            type="primary"
            @click="handleAddTenant"
          >
            {{ t('tenant.add') }}
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <div class="search-bar">
        <el-input
          v-model="searchKeyword"
          :placeholder="t('tenant.search')"
          clearable
          @clear="handleSearch"
          @keyup.enter="handleSearch"
          style="width: 300px; margin-right: 10px;"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
        <el-select
          v-model="statusFilter"
          :placeholder="t('tenant.status')"
          clearable
          style="width: 150px; margin-right: 10px;"
          @change="handleSearch"
        >
          <el-option :label="t('tenant.enabled')" :value="1" />
          <el-option :label="t('tenant.disabled')" :value="2" />
        </el-select>
        <el-button type="primary" @click="handleSearch">{{ t('common.search') }}</el-button>
        <el-button @click="handleReset">{{ t('common.reset') }}</el-button>
      </div>

      <el-table :data="tenants" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" :label="t('tenant.name')" width="150" />
        <el-table-column prop="code" :label="t('tenant.code')" width="120" />
        <el-table-column prop="contact_email" :label="t('tenant.contactEmail')" width="180">
          <template #default="scope">
            {{ scope.row.contact_email || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="contact_phone" :label="t('tenant.contactPhone')" width="130">
          <template #default="scope">
            {{ scope.row.contact_phone || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="status" :label="t('tenant.status')" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.status === 1 ? 'success' : 'danger'">
              {{ scope.row.status === 1 ? t('tenant.enabled') : t('tenant.disabled') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="ctime" :label="t('tenant.createTime')" width="180">
          <template #default="scope">
            {{ formatDate(scope.row.ctime) }}
          </template>
        </el-table-column>
        <el-table-column :label="t('tenant.operations')" width="150" fixed="right">
          <template #default="scope">
            <el-button
              v-permission="'tenant_update'"
              link
              type="primary"
              size="small"
              @click="handleEditTenant(scope.row)"
            >
              {{ t('common.edit') }}
            </el-button>
            <el-button
              v-permission="'tenant_delete'"
              link
              type="danger"
              size="small"
              @click="handleDeleteTenant(scope.row)"
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

    <!-- 新增租户弹框 -->
    <el-drawer
      v-model="addDialogVisible"
      :title="t('tenant.add')"
      direction="rtl"
      size="500px"
    >
      <el-form
        ref="addTenantFormRef"
        :model="addTenantForm"
        :rules="addTenantRules"
        label-width="100px"
      >
        <el-form-item :label="t('tenant.name')" prop="name">
          <el-input
            v-model="addTenantForm.name"
            :placeholder="t('tenant.namePlaceholder')"
            maxlength="128"
            show-word-limit
          />
        </el-form-item>
        <el-form-item :label="t('tenant.code')" prop="code">
          <el-input
            v-model="addTenantForm.code"
            :placeholder="t('tenant.codePlaceholder')"
            maxlength="64"
            show-word-limit
          />
        </el-form-item>
        <el-form-item :label="t('tenant.contactEmail')" prop="contact_email">
          <el-input
            v-model="addTenantForm.contact_email"
            :placeholder="t('tenant.contactEmailPlaceholder')"
            maxlength="128"
          />
        </el-form-item>
        <el-form-item :label="t('tenant.contactPhone')" prop="contact_phone">
          <el-input
            v-model="addTenantForm.contact_phone"
            :placeholder="t('tenant.contactPhonePlaceholder')"
            maxlength="32"
          />
        </el-form-item>
        <el-form-item :label="t('tenant.status')" prop="status">
          <el-radio-group v-model="addTenantForm.status">
            <el-radio :label="1">{{ t('tenant.enabled') }}</el-radio>
            <el-radio :label="2">{{ t('tenant.disabled') }}</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item :label="t('tenant.config')" prop="config">
          <el-input
            v-model="addTenantForm.config"
            type="textarea"
            :rows="4"
            :placeholder="t('tenant.configPlaceholder')"
          />
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

    <!-- 编辑租户弹框 -->
    <el-drawer
      v-model="editDialogVisible"
      :title="t('tenant.edit')"
      direction="rtl"
      size="500px"
    >
      <el-form
        ref="editTenantFormRef"
        :model="editTenantForm"
        :rules="editTenantRules"
        label-width="100px"
      >
        <el-form-item :label="t('tenant.name')" prop="name">
          <el-input
            v-model="editTenantForm.name"
            :placeholder="t('tenant.namePlaceholder')"
            maxlength="128"
            show-word-limit
          />
        </el-form-item>
        <el-form-item :label="t('tenant.code')" prop="code">
          <el-input
            v-model="editTenantForm.code"
            :placeholder="t('tenant.codePlaceholder')"
            maxlength="64"
            show-word-limit
          />
        </el-form-item>
        <el-form-item :label="t('tenant.contactEmail')" prop="contact_email">
          <el-input
            v-model="editTenantForm.contact_email"
            :placeholder="t('tenant.contactEmailPlaceholder')"
            maxlength="128"
          />
        </el-form-item>
        <el-form-item :label="t('tenant.contactPhone')" prop="contact_phone">
          <el-input
            v-model="editTenantForm.contact_phone"
            :placeholder="t('tenant.contactPhonePlaceholder')"
            maxlength="32"
          />
        </el-form-item>
        <el-form-item :label="t('tenant.status')" prop="status">
          <el-radio-group v-model="editTenantForm.status">
            <el-radio :label="1">{{ t('tenant.enabled') }}</el-radio>
            <el-radio :label="2">{{ t('tenant.disabled') }}</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item :label="t('tenant.config')" prop="config">
          <el-input
            v-model="editTenantForm.config"
            type="textarea"
            :rows="4"
            :placeholder="t('tenant.configPlaceholder')"
          />
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
import { Search } from '@element-plus/icons-vue'
import { useTenant } from '@/composables/useTenant'

const { t, locale } = useI18n()

// 使用租户管理 composable
const {
  tenants,
  loading,
  currentPage,
  pageSize,
  total,
  searchKeyword,
  statusFilter,
  addDialogVisible,
  submitLoading,
  addTenantForm,
  addTenantFormRef,
  editDialogVisible,
  editTenantForm,
  editTenantFormRef,
  fetchTenants,
  handleSearch,
  handleReset,
  handleSizeChange,
  handleCurrentChange,
  handleSubmitAdd,
  handleCancelAdd,
  handleSubmitEdit,
  handleCancelEdit,
  handleDeleteTenant
} = useTenant(t, locale)

// 表单验证规则
const addTenantRules = computed(() => ({
  name: [
    { required: true, message: t('tenant.name') + t('common.error.required'), trigger: 'blur' },
    { max: 128, message: t('tenant.nameLengthLimit'), trigger: 'blur' }
  ],
  code: [
    { required: true, message: t('tenant.code') + t('common.error.required'), trigger: 'blur' },
    { max: 64, message: t('tenant.codeLengthLimit'), trigger: 'blur' }
  ],
  contact_email: [
    { type: 'email', message: t('tenant.emailFormatError'), trigger: 'blur' }
  ],
  contact_phone: [
    { max: 32, message: t('tenant.phoneLengthLimit'), trigger: 'blur' }
  ]
}))

const editTenantRules = computed(() => ({
  name: [
    { required: true, message: t('tenant.name') + t('common.error.required'), trigger: 'blur' },
    { max: 128, message: t('tenant.nameLengthLimit'), trigger: 'blur' }
  ],
  code: [
    { required: true, message: t('tenant.code') + t('common.error.required'), trigger: 'blur' },
    { max: 64, message: t('tenant.codeLengthLimit'), trigger: 'blur' }
  ],
  contact_email: [
    { type: 'email', message: t('tenant.emailFormatError'), trigger: 'blur' }
  ],
  contact_phone: [
    { max: 32, message: t('tenant.phoneLengthLimit'), trigger: 'blur' }
  ]
}))

// 格式化日期
const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleString(locale.value === 'zh' ? 'zh-CN' : 'en-US')
}

// 新增租户
const handleAddTenant = () => {
  // 重置表单
  addTenantForm.value = {
    name: '',
    code: '',
    contact_email: '',
    contact_phone: '',
    status: 1,
    config: ''
  }

  // 清除表单验证状态
  if (addTenantFormRef.value) {
    addTenantFormRef.value.clearValidate()
  }

  // 打开弹框
  addDialogVisible.value = true
}

// 编辑租户
const handleEditTenant = (row) => {
  // 填充表单数据
  editTenantForm.value = {
    id: row.id,
    name: row.name,
    code: row.code,
    contact_email: row.contact_email || '',
    contact_phone: row.contact_phone || '',
    status: row.status,
    config: row.config || ''
  }

  // 清除表单验证状态
  if (editTenantFormRef.value) {
    editTenantFormRef.value.clearValidate()
  }

  // 打开弹框
  editDialogVisible.value = true
}

// 页面加载时获取数据
onMounted(() => {
  fetchTenants()
})
</script>

<style scoped>
.tenant-management {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.search-bar {
  display: flex;
  align-items: center;
  margin-bottom: 20px;
}
</style>