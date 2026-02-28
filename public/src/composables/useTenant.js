import { ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import axios from 'axios'

export function useTenant(t, locale) {
  // 响应式数据
  const tenants = ref([])
  const loading = ref(false)
  const currentPage = ref(1)
  const pageSize = ref(10)
  const total = ref(0)
  const searchKeyword = ref('')
  const statusFilter = ref(null)

  // 新增租户对话框
  const addDialogVisible = ref(false)
  const submitLoading = ref(false)
  const addTenantForm = ref({
    name: '',
    code: '',
    contact_email: '',
    contact_phone: '',
    status: 1,
    config: ''
  })
  const addTenantFormRef = ref(null)

  // 编辑租户对话框
  const editDialogVisible = ref(false)
  const editTenantForm = ref({
    id: 0,
    name: '',
    code: '',
    contact_email: '',
    contact_phone: '',
    status: 1,
    config: ''
  })
  const editTenantFormRef = ref(null)

  // 获取租户列表
  const fetchTenants = async () => {
    loading.value = true
    try {
      const token = localStorage.getItem('token')
      if (!token) {
        ElMessage.error(t('login.loginFailed'))
        return
      }

      const params = {
        page: currentPage.value,
        page_size: pageSize.value,
        order_by: 'id desc',
        keyword: searchKeyword.value
      }

      if (statusFilter.value !== null) {
        params.status = statusFilter.value
      }

      const response = await axios.get('/api/admin/v1/tenant/list', {
        headers: {
          'Authorization': `Bearer ${token}`,
          'Accept-Language': locale.value
        },
        params
      })

      if (response.data.code === 200) {
        tenants.value = response.data.data.list || []
        total.value = response.data.data.total || 0
      } else {
        ElMessage.error(response.data.message || t('common.failed'))
      }
    } catch (error) {
      console.error('获取租户列表失败:', error)
      ElMessage.error(t('common.error.systemError'))
    } finally {
      loading.value = false
    }
  }

  // 搜索
  const handleSearch = () => {
    currentPage.value = 1
    fetchTenants()
  }

  // 重置
  const handleReset = () => {
    searchKeyword.value = ''
    statusFilter.value = null
    currentPage.value = 1
    fetchTenants()
  }

  // 分页大小改变
  const handleSizeChange = (val) => {
    pageSize.value = val
    currentPage.value = 1
    fetchTenants()
  }

  // 当前页改变
  const handleCurrentChange = (val) => {
    currentPage.value = val
    fetchTenants()
  }

  // 处理提交添加
  const handleSubmitAdd = async () => {
    if (!addTenantFormRef.value) return

    await addTenantFormRef.value.validate(async (valid) => {
      if (valid) {
        submitLoading.value = true
        try {
          const token = localStorage.getItem('token')
          if (!token) {
            ElMessage.error(t('login.loginFailed'))
            return
          }

          const response = await axios.post('/api/admin/v1/tenant/create', addTenantForm.value, {
            headers: {
              'Authorization': `Bearer ${token}`,
              'Accept-Language': locale.value
            }
          })

          if (response.data.code === 200) {
            ElMessage.success(t('tenant.addSuccess'))
            addDialogVisible.value = false
            // 刷新列表
            fetchTenants()
          } else {
            ElMessage.error(response.data.message || t('common.failed'))
          }
        } catch (error) {
          console.error('添加租户失败:', error)
          ElMessage.error(t('common.failed'))
        } finally {
          submitLoading.value = false
        }
      }
    })
  }

  // 处理取消添加
  const handleCancelAdd = () => {
    addDialogVisible.value = false
    // 清除表单验证状态
    if (addTenantFormRef.value) {
      addTenantFormRef.value.clearValidate()
    }
  }

  // 处理提交编辑
  const handleSubmitEdit = async () => {
    if (!editTenantFormRef.value) return

    await editTenantFormRef.value.validate(async (valid) => {
      if (valid) {
        submitLoading.value = true
        try {
          const token = localStorage.getItem('token')
          if (!token) {
            ElMessage.error(t('login.loginFailed'))
            return
          }

          const response = await axios.post('/api/admin/v1/tenant/update', editTenantForm.value, {
            headers: {
              'Authorization': `Bearer ${token}`,
              'Accept-Language': locale.value
            }
          })

          if (response.data.code === 200) {
            ElMessage.success(t('tenant.editSuccess'))
            editDialogVisible.value = false
            // 刷新列表
            fetchTenants()
          } else {
            ElMessage.error(response.data.message || t('common.failed'))
          }
        } catch (error) {
          console.error('更新租户失败:', error)
          ElMessage.error(t('common.failed'))
        } finally {
          submitLoading.value = false
        }
      }
    })
  }

  // 处理取消编辑
  const handleCancelEdit = () => {
    editDialogVisible.value = false
    // 清除表单验证状态
    if (editTenantFormRef.value) {
      editTenantFormRef.value.clearValidate()
    }
  }

  // 删除租户
  const handleDeleteTenant = (row) => {
    ElMessageBox.confirm(
      t('tenant.deleteConfirm'),
      t('common.confirm'),
      {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        type: 'warning'
      }
    ).then(async () => {
      try {
        const token = localStorage.getItem('token')
        if (!token) {
          ElMessage.error(t('login.loginFailed'))
          return
        }

        const response = await axios.post('/api/admin/v1/tenant/delete', {
          id: row.id
        }, {
          headers: {
            'Authorization': `Bearer ${token}`,
            'Accept-Language': locale.value
          }
        })

        if (response.data.code === 200) {
          ElMessage.success(t('tenant.deleteSuccess'))
          // 刷新列表
          fetchTenants()
        } else {
          ElMessage.error(response.data.message || t('common.failed'))
        }
      } catch (error) {
        console.error('删除租户失败:', error)
        ElMessage.error(t('common.failed'))
      }
    }).catch(() => {
      // 用户取消删除
    })
  }

  return {
    // 响应式数据
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
    // 方法
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
  }
}