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
                @keyup.enter="handleAddSearchLocation"
                style="width: 300px;"
              >
                <template #append>
                  <el-button @click="handleAddSearchLocation">{{ t('common.search') }}</el-button>
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
      size="800px"
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

        <!-- 地图容器 -->
        <el-form-item label="">
          <div style="position: relative; width: 100%;">
            <!-- 地图搜索浮层 -->
            <div class="map-search-overlay">
              <el-input
                v-model="editSearchLocation"
                :placeholder="t('position.locationPlaceholder')"
                clearable
                @keyup.enter="onEditSearchLocation"
                style="width: 300px;"
              >
                <template #append>
                  <el-button @click="onEditSearchLocation">{{ t('common.search') }}</el-button>
                </template>
              </el-input>
            </div>

            <div id="editMapContainer" style="width: 100%; height: 400px; border: 1px solid #dcdfe6; border-radius: 4px;"></div>
            <div style="margin-top: 8px; color: #909399; font-size: 12px;">
              {{ t('position.mapTip') || '点击地图选择位置，或使用上方搜索框搜索地点' }}
            </div>
          </div>
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
            :disabled="true"
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
            :disabled="true"
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
import { Search } from '@element-plus/icons-vue'
import axios from 'axios'
import { useServiceSettings } from '@/composables/useSettings'
import { useMap } from '@/composables/useMap'
import { usePosition } from '@/composables/usePosition'

const { t, locale } = useI18n()
const { settings: serviceSettings, loadSettings: loadServiceSettings } = useServiceSettings()

// 高德地图API密钥
const amapKey = ref('')

// 获取高德地图密钥
const fetchAMapKey = async () => {
  try {
    const token = localStorage.getItem('token')
    if (!token) {
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
      }
    }
  } catch (error) {
    console.error('获取地图密钥失败:', error)
  }
}

// 使用地图 composable
const {
  initAddMap,
  handleSearchLocation,
  handleMapClick,
  initEditMap,
  handleEditSearchLocation,
  handleEditMapClick,
  destroyAddMap,
  destroyEditMap
} = useMap(serviceSettings, amapKey)

// 使用位置管理 composable
const {
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
} = usePosition(t, locale, serviceSettings)

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
    await initAddMap('', handleAddMapClick)
  }
}

// 新增地图搜索地点
const handleAddSearchLocation = () => {
  handleSearchLocation(searchLocation.value, (poi, location) => {
    addPositionForm.value.location = poi.name
    addPositionForm.value.longitude = location.lng
    addPositionForm.value.latitude = location.lat
  })
}

// 新增地图点击
const handleAddMapClick = (lnglat) => {
  handleMapClick(lnglat, (address, lng, lat) => {
    if (address) {
      addPositionForm.value.location = address
    }
    addPositionForm.value.longitude = lng
    addPositionForm.value.latitude = lat
  })
}

// 编辑位置
const handleEditPosition = async (row) => {
  // 加载服务配置（service_region）
  await loadServiceSettings()

  // 填充表单数据
  editPositionForm.value = {
    id: row.id,
    location: row.location,
    custom_name: row.custom_name || '',
    longitude: row.longitude,
    latitude: row.latitude
  }
  editSearchLocation.value = ''

  // 清除表单验证状态
  if (editPositionFormRef.value) {
    editPositionFormRef.value.clearValidate()
  }

  // 打开弹框
  editDialogVisible.value = true

  // 初始化地图
  await nextTick()
  if (amapKey.value) {
    await initEditMap(editPositionForm.value, onEditMapClick)
  }
}

// 编辑地图搜索地点
const onEditSearchLocation = () => {
  handleEditSearchLocation(editSearchLocation.value, (poi, location) => {
    editPositionForm.value.location = poi.name
    editPositionForm.value.longitude = location.lng
    editPositionForm.value.latitude = location.lat
  })
}

// 编辑地图点击
const onEditMapClick = (lnglat) => {
  handleEditMapClick(lnglat, (address, lng, lat) => {
    if (address) {
      editPositionForm.value.location = address
    }
    editPositionForm.value.longitude = lng
    editPositionForm.value.latitude = lat
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

#editMapContainer {
  border-radius: 4px;
  overflow: hidden;
}

#editMapContainer :deep(.amap-toolbar) {
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