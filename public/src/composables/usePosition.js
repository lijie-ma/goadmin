import { ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import axios from 'axios'

export function usePosition(t, locale, serviceSettings) {
  // 响应式数据
  const positions = ref([])
  const loading = ref(false)
  const currentPage = ref(1)
  const pageSize = ref(10)
  const total = ref(0)
  const searchKeyword = ref('')

  // 新增位置对话框
  const addDialogVisible = ref(false)
  const submitLoading = ref(false)
  const searchLocation = ref('')
  const addPositionForm = ref({
    location: '',
    custom_name: '',
    longitude: 0,
    latitude: 0
  })
  const addPositionFormRef = ref(null)

  // 编辑位置对话框
  const editDialogVisible = ref(false)
  const editSearchLocation = ref('')
  const editPositionForm = ref({
    id: 0,
    location: '',
    custom_name: '',
    longitude: 0,
    latitude: 0
  })
  const editPositionFormRef = ref(null)

  // 获取位置列表
  const fetchPositions = async () => {
    loading.value = true
    try {
      const token = localStorage.getItem('token')
      if (!token) {
        ElMessage.error(t('login.loginFailed'))
        return
      }

      const response = await axios.get('/api/admin/v1/position/list', {
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
        positions.value = response.data.data.list || []
        total.value = response.data.data.total || 0
      } else {
        ElMessage.error(response.data.message || t('common.failed'))
      }
    } catch (error) {
      console.error('获取位置列表失败:', error)
      ElMessage.error(t('common.error.systemError'))
    } finally {
      loading.value = false
    }
  }

  // 搜索
  const handleSearch = () => {
    currentPage.value = 1
    fetchPositions()
  }

  // 重置
  const handleReset = () => {
    searchKeyword.value = ''
    currentPage.value = 1
    fetchPositions()
  }

  // 分页大小改变
  const handleSizeChange = (val) => {
    pageSize.value = val
    currentPage.value = 1
    fetchPositions()
  }

  // 当前页改变
  const handleCurrentChange = (val) => {
    currentPage.value = val
    fetchPositions()
  }

  // 处理提交添加
  const handleSubmitAdd = async () => {
    if (!addPositionFormRef.value) return

    await addPositionFormRef.value.validate(async (valid) => {
      if (valid) {
        submitLoading.value = true
        try {
          const token = localStorage.getItem('token')
          if (!token) {
            ElMessage.error(t('login.loginFailed'))
            return
          }

          // 构建提交数据，自动添加 city 字段
          const submitData = {
            ...addPositionForm.value,
            city: serviceSettings?.region || ''
          }

          const response = await axios.post('/api/admin/v1/position/create', submitData, {
            headers: {
              'Authorization': `Bearer ${token}`,
              'Accept-Language': locale.value
            }
          })

          if (response.data.code === 200) {
            ElMessage.success(t('position.addSuccess'))
            addDialogVisible.value = false
            // 刷新列表
            fetchPositions()
          } else {
            ElMessage.error(response.data.message || t('common.failed'))
          }
        } catch (error) {
          console.error('添加位置失败:', error)
          ElMessage.error(t('common.failed'))
        } finally {
          submitLoading.value = false
        }
      }
    })
  }

  // 处理取消添加
  const handleCancelAdd = (destroyAddMap) => {
    addDialogVisible.value = false
    // 清除表单验证状态
    if (addPositionFormRef.value) {
      addPositionFormRef.value.clearValidate()
    }
    // 清理地图资源
    destroyAddMap()
  }

  // 处理提交编辑
  const handleSubmitEdit = async () => {
    if (!editPositionFormRef.value) return

    await editPositionFormRef.value.validate(async (valid) => {
      if (valid) {
        submitLoading.value = true
        try {
          const token = localStorage.getItem('token')
          if (!token) {
            ElMessage.error(t('login.loginFailed'))
            return
          }

          // 构建提交数据，自动添加 city 字段
          const submitData = {
            ...editPositionForm.value,
            city: serviceSettings?.region || ''
          }

          const response = await axios.post('/api/admin/v1/position/update', submitData, {
            headers: {
              'Authorization': `Bearer ${token}`,
              'Accept-Language': locale.value
            }
          })

          if (response.data.code === 200) {
            ElMessage.success(t('position.editSuccess'))
            editDialogVisible.value = false
            // 刷新列表
            fetchPositions()
          } else {
            ElMessage.error(response.data.message || t('common.failed'))
          }
        } catch (error) {
          console.error('更新位置失败:', error)
          ElMessage.error(t('common.failed'))
        } finally {
          submitLoading.value = false
        }
      }
    })
  }

  // 处理取消编辑
  const handleCancelEdit = (destroyEditMap) => {
    editDialogVisible.value = false
    // 清除表单验证状态
    if (editPositionFormRef.value) {
      editPositionFormRef.value.clearValidate()
    }
    // 清理地图资源
    destroyEditMap()
  }

  // 删除位置
  const handleDeletePosition = (row) => {
    ElMessageBox.confirm(
      t('position.deleteConfirm'),
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

        const response = await axios.post('/api/admin/v1/position/delete', {
          id: row.id
        }, {
          headers: {
            'Authorization': `Bearer ${token}`,
            'Accept-Language': locale.value
          }
        })

        if (response.data.code === 200) {
          ElMessage.success(t('position.deleteSuccess'))
          // 刷新列表
          fetchPositions()
        } else {
          ElMessage.error(response.data.message || t('common.failed'))
        }
      } catch (error) {
        console.error('删除位置失败:', error)
        ElMessage.error(t('common.failed'))
      }
    }).catch(() => {
      // 用户取消删除
    })
  }

  return {
    // 响应式数据
    positions,
    loading,
    currentPage,
    pageSize,
    total,
    searchKeyword,
    addDialogVisible,
    submitLoading,
    searchLocation,
    addPositionForm,
    addPositionFormRef,
    editDialogVisible,
    editSearchLocation,
    editPositionForm,
    editPositionFormRef,
    // 方法
    fetchPositions,
    handleSearch,
    handleReset,
    handleSizeChange,
    handleCurrentChange,
    handleSubmitAdd,
    handleCancelAdd,
    handleSubmitEdit,
    handleCancelEdit,
    handleDeletePosition
  }
}