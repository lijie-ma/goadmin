<template>
  <div class="position-management">
    <h1>{{ t('position.title') }}</h1>
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('position.title') }}</span>
          <el-button
            v-permission="'position_create'"
            type="primary"
            @click="handleAddPosition"
          >
            {{ t('position.add') }}
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <div class="search-bar">
        <el-input
          v-model="searchKeyword"
          :placeholder="t('position.search')"
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

      <el-table :data="positions" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="location" :label="t('position.location')" />
        <el-table-column prop="custom_name" :label="t('position.customName')" width="150">
          <template #default="scope">
            {{ scope.row.custom_name || '-' }}
          </template>
        </el-table-column>
        <el-table-column :label="t('position.coordinates')" width="180">
          <template #default="scope">
            {{ scope.row.longitude }}, {{ scope.row.latitude }}
          </template>
        </el-table-column>
        <el-table-column prop="creator" :label="t('position.creator')" width="120" />
        <el-table-column prop="ctime" :label="t('position.createTime')" width="180">
          <template #default="scope">
            {{ formatDate(scope.row.ctime) }}
          </template>
        </el-table-column>
        <el-table-column :label="t('position.operations')" width="150" fixed="right">
          <template #default="scope">
            <el-button
              v-permission="'position_update'"
              link
              type="primary"
              size="small"
              @click="handleEditPosition(scope.row)"
            >
              {{ t('common.edit') }}
            </el-button>
            <el-button
              v-permission="'position_delete'"
              link
              type="danger"
              size="small"
              @click="handleDeletePosition(scope.row)"
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

    <!-- 新增位置弹框 -->
    <el-drawer
      v-model="addDialogVisible"
      :title="t('position.add')"
      direction="rtl"
      size="800px"
    >
      <el-form
        ref="addPositionFormRef"
        :model="addPositionForm"
        :rules="addPositionRules"
        label-width="100px"
      >
        <el-form-item :label="t('position.location')" prop="location">
          <el-input
            v-model="addPositionForm.location"
            :placeholder="t('position.locationPlaceholder')"
            maxlength="128"
            show-word-limit
          />
        </el-form-item>

        <!-- 地图容器 -->
        <el-form-item label="">
          <div style="position: relative; width: 100%;">
            <!-- 地图搜索浮层 -->
            <div class="map-search-overlay">
              <el-input
                v-model="searchLocation"
                :placeholder="t('position.locationPlaceholder')"
                clearable
                @keyup.enter="handleSearchLocation"
                style="width: 300px;"
              >
                <template #append>
                  <el-button @click="handleSearchLocation">{{ t('common.search') }}</el-button>
                </template>
              </el-input>
            </div>

            <div id="addMapContainer" style="width: 100%; height: 400px; border: 1px solid #dcdfe6; border-radius: 4px;"></div>
            <div style="margin-top: 8px; color: #909399; font-size: 12px;">
              {{ t('position.mapTip') || '点击地图选择位置，或使用上方搜索框搜索地点' }}
            </div>
          </div>
        </el-form-item>

        <el-form-item :label="t('position.customName')" prop="custom_name">
          <el-input
            v-model="addPositionForm.custom_name"
            :placeholder="t('position.customNamePlaceholder')"
            maxlength="128"
            show-word-limit
          />
        </el-form-item>
        <el-form-item :label="t('position.longitude')" prop="longitude">
          <el-input-number
            v-model="addPositionForm.longitude"
            :precision="6"
            :step="0.000001"
            :min="-180"
            :max="180"
            style="width: 100%"
            :disabled="true"
          />
        </el-form-item>
        <el-form-item :label="t('position.latitude')" prop="latitude">
          <el-input-number
            v-model="addPositionForm.latitude"
            :precision="6"
            :step="0.000001"
            :min="-90"
            :max="90"
            style="width: 100%"
            :disabled="true"
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

    <!-- 编辑位置弹框 -->
    <el-drawer
      v-model="editDialogVisible"
      :title="t('position.edit')"
      direction="rtl"
      size="500px"
    >
      <el-form
        ref="editPositionFormRef"
        :model="editPositionForm"
        :rules="editPositionRules"
        label-width="100px"
      >
        <el-form-item :label="t('position.location')" prop="location">
          <el-input
            v-model="editPositionForm.location"
            :placeholder="t('position.locationPlaceholder')"
            maxlength="128"
            show-word-limit
          />
        </el-form-item>
        <el-form-item :label="t('position.customName')" prop="custom_name">
          <el-input
            v-model="editPositionForm.custom_name"
            :placeholder="t('position.customNamePlaceholder')"
            maxlength="128"
            show-word-limit
          />
        </el-form-item>
        <el-form-item :label="t('position.longitude')" prop="longitude">
          <el-input-number
            v-model="editPositionForm.longitude"
            :precision="6"
            :step="0.000001"
            :min="-180"
            :max="180"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item :label="t('position.latitude')" prop="latitude">
          <el-input-number
            v-model="editPositionForm.latitude"
            :precision="6"
            :step="0.000001"
            :min="-90"
            :max="90"
            style="width: 100%"
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
import { ref, onMounted, computed, nextTick } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search } from '@element-plus/icons-vue'
import axios from 'axios'
import { useServiceSettings } from '@/composables/useSettings'

const { t, locale } = useI18n()
const { settings: serviceSettings, loadSettings: loadServiceSettings } = useServiceSettings()

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

// 地图相关
let addMap = null
let addMarker = null
let addPlaceSearch = null
let addGeocoder = null
const mapLoaded = ref(false)
const amapKey = ref('') // 高德地图API密钥

// 编辑位置对话框
const editDialogVisible = ref(false)
const editPositionForm = ref({
  id: 0,
  location: '',
  custom_name: '',
  longitude: 0,
  latitude: 0
})
const editPositionFormRef = ref(null)

// 表单验证规则
const addPositionRules = computed(() => ({
  location: [
    { required: true, message: t('position.location') + t('common.error.required'), trigger: 'blur' }
  ],
  longitude: [
    { required: true, message: t('position.longitude') + t('common.error.required'), trigger: 'blur' }
  ],
  latitude: [
    { required: true, message: t('position.latitude') + t('common.error.required'), trigger: 'blur' }
  ]
}))

const editPositionRules = computed(() => ({
  location: [
    { required: true, message: t('position.location') + t('common.error.required'), trigger: 'blur' }
  ],
  longitude: [
    { required: true, message: t('position.longitude') + t('common.error.required'), trigger: 'blur' }
  ],
  latitude: [
    { required: true, message: t('position.latitude') + t('common.error.required'), trigger: 'blur' }
  ]
}))

// 格式化日期
const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleString(locale.value === 'zh' ? 'zh-CN' : 'en-US')
}

// 加载高德地图API
const loadAMapScript = () => {
  return new Promise((resolve, reject) => {
    if (window.AMap) {
      resolve()
      return
    }

    // 检查安全配置是否已设置
    if (!window._AMapSecurityConfig) {
      console.error('AMap security config not set')
      reject(new Error('AMap security config not set'))
      return
    }

    const script = document.createElement('script')
    script.type = 'text/javascript'
    script.src = `https://webapi.amap.com/maps?v=2.0&key=${amapKey.value}&plugin=AMap.PlaceSearch,AMap.Geocoder`
    script.onload = () => {
      mapLoaded.value = true
      resolve()
    }
    script.onerror = () => {
      reject(new Error('Failed to load AMap script'))
    }
    document.head.appendChild(script)
  })
}

// 初始化新增地图
// 初始化新增地图
const initAddMap = async (initialCityParam = '') => {
  try {
    await loadAMapScript()
    await nextTick()

    const mapContainer = document.getElementById('addMapContainer')
    if (!mapContainer) {
      console.error('Map container not found')
      return
    }

    // ✅ 确定初始城市优先级
    const initialCity =
      initialCityParam ||
      serviceSettings.region ||
      '全国'

    console.log('Initializing map with city:', initialCity)

    // 创建地图实例
    addMap = new AMap.Map('addMapContainer', {
      zoom: 11,
      viewMode: '2D',
      resizeEnable: true
    })

    // 动态加载控件插件
    AMap.plugin(['AMap.ControlBar', 'AMap.Scale'], () => {
      addMap.addControl(
        new AMap.ControlBar({
          position: { right: '10px', top: '10px' }
        })
      )
      addMap.addControl(new AMap.Scale())
    })

    // 初始化搜索与地理编码服务
    addPlaceSearch = new AMap.PlaceSearch({
      city: initialCity,
      pageSize: 5,
      pageIndex: 1,
      extensions: 'all'
    })

    addGeocoder = new AMap.Geocoder({
      city: initialCity
    })

    // 地图点击事件
    addMap.on('click', (e) => {
      handleMapClick(e.lnglat)
    })

    // ✅ 设置地图中心
    if (initialCity && initialCity !== '全国') {
      setMapCenterByCity(initialCity)
    } else {
      console.log('No specific city, using default China center')
      addMap.setZoom(4)
      addMap.setCenter([104.195397, 35.86166]) // 中国地理中心点
    }
  } catch (error) {
    console.error('Failed to initialize map:', error)
    ElMessage.error(t('position.mapLoadFailed') || '地图加载失败')
  }
}


// 根据城市设置地图中心
const setMapCenterByCity = (city) => {
  if (!addGeocoder || !city) {
    console.log('Cannot set map center: geocoder or city missing')
    return
  }

  addGeocoder.getLocation(city, (status, result) => {
    if (status === 'complete' && result.geocodes?.length > 0) {
      const location = result.geocodes[0].location
      addMap.setCenter([location.lng, location.lat])
      addMap.setZoom(11)
    } else {
      console.warn(`City not found: ${city}`)
      ElMessage.warning(t('position.cityNotFound') || `未找到城市: ${city}`)
      addMap.setZoom(4)
      addMap.setCenter([104.195397, 35.86166])
    }
  })
}


// 搜索地点
const handleSearchLocation = () => {
  if (!searchLocation.value || !addPlaceSearch) {
    return
  }

  addPlaceSearch.search(searchLocation.value, (status, result) => {
    if (status === 'complete' && result.poiList.pois.length > 0) {
      const poi = result.poiList.pois[0]
      const location = poi.location

      // 设置地图中心和标记
      addMap.setCenter([location.lng, location.lat])
      addMap.setZoom(15)

      // 更新标记
      updateMarker(location.lng, location.lat)

      // 填充表单数据
      addPositionForm.value.location = poi.name
      addPositionForm.value.longitude = location.lng
      addPositionForm.value.latitude = location.lat

      ElMessage.success(t('position.locationFound') || '已找到位置')
    } else {
      ElMessage.warning(t('position.locationNotFound') || '未找到该位置')
    }
  })
}

// 处理地图点击
const handleMapClick = (lnglat) => {
  updateMarker(lnglat.lng, lnglat.lat)

  // 反向地理编码获取地址
  if (addGeocoder) {
    addGeocoder.getAddress(lnglat, (status, result) => {
      if (status === 'complete' && result.regeocode) {
        addPositionForm.value.location = result.regeocode.formattedAddress
      }
    })
  }

  addPositionForm.value.longitude = lnglat.lng
  addPositionForm.value.latitude = lnglat.lat
}

// 更新地图标记
const updateMarker = (lng, lat) => {
  if (addMarker) {
    addMarker.setPosition([lng, lat])
  } else {
    addMarker = new AMap.Marker({
      position: [lng, lat],
      map: addMap,
      draggable: true
    })

    // 标记拖拽事件
    addMarker.on('dragend', (e) => {
      const position = e.target.getPosition()
      handleMapClick(position)
    })
  }
}

// 获取高德地图密钥
const fetchAMapKey = async () => {
  try {
    const token = localStorage.getItem('token')
    if (!token) {
      ElMessage.error(t('login.loginFailed'))
      return
    }

    const response = await axios.get('/api/admin/v1/setting/decrypted', {
      headers: {
        'Authorization': `Bearer ${token}`,
        'Accept-Language': locale.value
      },
      params: {
        name: 'map_config'
      }
    })

    if (response.data.code === 200 && response.data.data) {
      const mapConfig = JSON.parse(response.data.data)
      amapKey.value = mapConfig.map_ak
      window._AMapSecurityConfig = {
          securityJsCode: mapConfig.map_scode,
      };
    } else {
      ElMessage.warning(t('position.mapKeyNotFound') || '未配置地图密钥')
    }
  } catch (error) {
    console.error('获取地图密钥失败:', error)
    ElMessage.warning(t('position.mapKeyNotFound') || '未配置地图密钥')
  }
}

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


// 新增位置
const handleAddPosition = async () => {
  // 加载服务配置（service_region）
  await loadServiceSettings()

  // 加载地图配置（map_config）
  await fetchAMapKey()

  // 重置表单
  addPositionForm.value = {
    location: '',
    custom_name: '',
    longitude: 0,
    latitude: 0
  }
  searchLocation.value = ''

  // 清除表单验证状态
  if (addPositionFormRef.value) {
    addPositionFormRef.value.clearValidate()
  }

  // 打开弹框
  addDialogVisible.value = true

  // 初始化地图
  await nextTick()
  if (amapKey.value) {
    await initAddMap()
  } else {
    ElMessage.warning(t('position.mapKeyNotFound') || '未配置地图密钥')
  }
}

// 处理取消添加
const handleCancelAdd = () => {
  addDialogVisible.value = false
  // 清除表单验证状态
  if (addPositionFormRef.value) {
    addPositionFormRef.value.clearValidate()
  }
  // 清理地图资源
  if (addMap) {
    addMap.destroy()
    addMap = null
  }
  addMarker = null
  addPlaceSearch = null
  addGeocoder = null
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
          city: serviceSettings.region || ''
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

// 编辑位置
const handleEditPosition = (row) => {
  // 填充表单数据
  editPositionForm.value = {
    id: row.id,
    location: row.location,
    custom_name: row.custom_name || '',
    longitude: row.longitude,
    latitude: row.latitude
  }
  // 清除表单验证状态
  if (editPositionFormRef.value) {
    editPositionFormRef.value.clearValidate()
  }
  // 打开弹框
  editDialogVisible.value = true
}

// 处理取消编辑
const handleCancelEdit = () => {
  editDialogVisible.value = false
  // 清除表单验证状态
  if (editPositionFormRef.value) {
    editPositionFormRef.value.clearValidate()
  }
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
          city: serviceSettings.region || ''
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

// 页面加载时获取数据
onMounted(async () => {
  fetchPositions()
  // 预加载地图密钥
  await fetchAMapKey()
})
</script>

<style scoped>
.position-management {
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

#addMapContainer {
  border-radius: 4px;
  overflow: hidden;
}

#addMapContainer :deep(.amap-toolbar) {
  right: 10px !important;
  top: 10px !important;
}

.map-search-overlay {
  position: absolute;
  top: 10px;
  left: 10px;
  z-index: 100;
  background: white;
  padding: 8px;
  border-radius: 4px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}
</style>