<template>
  <div class="settings-container">
    <!-- 标签栏 -->
    <div class="tab-header">
      <div
        v-for="tab in tabs"
        :key="tab.key"
        :class="['tab-item', { active: activeTab === tab.key }]"
        @click="switchTab(tab.key)"
      >
        {{ tab.label }}
      </div>
    </div>

    <!-- 内容区域 -->
    <div class="tab-content">
      <!-- 系统设置内容 -->
      <div v-show="activeTab === 'system'" class="tab-pane">
        <el-card class="settings-card" v-loading="loading">
          <el-form :model="settings" label-width="120px">
            <!-- 基本设置 -->
            <el-divider content-position="left">{{ t('settings.basic') }}</el-divider>
            <el-form-item :label="t('settings.systemName')">
              <el-input v-model="settings.system_name" :placeholder="t('settings.systemNamePlaceholder')" />
            </el-form-item>
            <el-form-item :label="t('settings.systemLogo')">
              <el-upload
                class="avatar-uploader"
                action="#"
                :show-file-list="false"
                :before-upload="beforeLogoUpload"
                :disabled="uploading"
              >
                <div v-if="uploading" class="avatar-uploader-loading">
                  <el-icon class="is-loading"><Loading /></el-icon>
                </div>
                <img v-else-if="settings.logo" :src="settings.logo" class="avatar" />
                <el-icon v-else class="avatar-uploader-icon"><Plus /></el-icon>
              </el-upload>
            </el-form-item>

            <!-- 系统配置 -->
            <el-divider content-position="left">{{ t('settings.systemConfig') }}</el-divider>
            <el-form-item :label="t('settings.systemLanguage')">
              <el-select v-model="settings.language" :placeholder="t('settings.selectLanguage')">
                <el-option :label="t('settings.simplifiedChinese')" value="zh_CN" />
                <el-option :label="t('settings.english')" value="en_US" />
              </el-select>
            </el-form-item>
            <!-- 安全设置 -->
            <el-divider content-position="left">{{ t('settings.securitySettings') }}</el-divider>
            <el-form-item :label="t('settings.loginCaptcha')">
              <el-switch
                v-model="settings.admin"
                :active-value="1"
                :inactive-value="0"
              />
            </el-form-item>

            <!-- 保存按钮 -->
            <el-form-item>
              <el-button
                v-permission="'server_set'"
                type="primary"
                @click="saveSettings"
                :loading="saving"
              >
                {{ t('settings.saveSettings') }}
              </el-button>
              <el-button
                v-permission="'server_set'"
                @click="resetSettings"
              >
                {{ t('settings.reset') }}
              </el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </div>

      <!-- 第三方服务配置内容 -->
      <div v-show="activeTab === 'thirdParty'" class="tab-pane">
        <el-card class="settings-card" v-loading="servicesLoading">
          <el-form :model="thirdPartySettings" label-width="120px">
            <!-- 地图服务配置 -->
            <el-divider content-position="left">{{ t('settings.mapService') }}</el-divider>
            <el-form-item :label="t('settings.mapAk')">
              <el-input
                v-model="thirdPartySettings.map_ak"
                :placeholder="t('settings.mapAkPlaceholder')"
              />
            </el-form-item>

            <!-- 保存按钮 -->
            <el-form-item>
              <el-button
                v-permission="'server_set'"
                type="primary"
                @click="saveThirdPartySettings"
                :loading="serviceSaving"
              >
                {{ t('settings.saveSettings') }}
              </el-button>
              <el-button
                v-permission="'server_set'"
                @click="resetThirdPartySettings"
              >
                {{ t('settings.reset') }}
              </el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { Plus, Loading } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'
import axios from 'axios'

const { t } = useI18n()
const loading = ref(false)
const saving = ref(false)
const uploading = ref(false)
const servicesLoading = ref(false)
const serviceSaving = ref(false)

// 标签页配置 - 使用computed实现多语言适配
const tabs = computed(() => [
  { key: 'system', label: t('settings.tabSystem') },
  { key: 'thirdParty', label: t('settings.tabThirdParty') }
])

// 当前激活的标签
const activeTab = ref('system')

// 切换标签
const switchTab = (key) => {
  activeTab.value = key
}

const settings = reactive({
  system_name: t('app.title'),
  logo: '',
  language: 'zh_CN',
  admin: 1
})

const thirdPartySettings = reactive({
  map_ak: ''
})

const beforeLogoUpload = async (file) => {
  const isImage = file.type.startsWith('image/')
  const isLt2M = file.size / 1024 / 1024 < 2

  if (!isImage) {
    ElMessage.error(t('settings.uploadImageOnly'))
    return false
  }
  if (!isLt2M) {
    ElMessage.error(t('settings.uploadSizeLimit'))
    return false
  }

  // 上传文件到服务器
  uploading.value = true
  try {
    const formData = new FormData()
    formData.append('file', file)

    const response = await axios.post('/api/admin/v1/upload/file', formData, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`,
        'Content-Type': 'multipart/form-data'
      }
    })

    if (response.data.code === 200) {
      // 使用服务器返回的完整URL
      settings.logo = response.data.data.url
      ElMessage.success(t('settings.uploadSuccess'))
    } else {
      ElMessage.error(response.data.message || t('settings.uploadFailed'))
    }
  } catch (error) {
    console.error('上传失败:', error)
    ElMessage.error(t('settings.uploadFailed') + ': ' + (error.response?.data?.message || error.message))
  } finally {
    uploading.value = false
  }

  return false
}

const loadSettings = async () => {
  loading.value = true
  try {
    const response = await axios.get('/api/admin/v1/setting/get_settings', {
      headers: {
        'Authorization': `Bearer ${sessionStorage.getItem('token')}`
      }
    })
    if (response.data.code === 200) {
      Object.assign(settings, response.data.data)
      // ElMessage.success(t('settings.loadSuccess'))
    } else {
      throw new Error(response.data.message)
    }
  } catch (error) {
    ElMessage.error(t('settings.loadFailed') + ': ' + error.message)
  } finally {
    loading.value = false
  }
}

const saveSettings = async () => {
  saving.value = true
  try {
    const response = await axios.post('/api/admin/v1/setting/set_settings', settings, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })
    if (response.data.code === 200) {
      ElMessage.success(t('settings.saveSuccess'))
      localStorage.setItem('language', settings.language)
    } else {
      throw new Error(response.data.message)
    }
  } catch (error) {
    ElMessage.error(t('settings.saveFailed') + ': ' + error.message)
  } finally {
    saving.value = false
  }
}

const defaultSettings = {
  system_name: t('app.title'),
  logo: '',
  language: 'zh_CN',
  admin: 1
}

const resetSettings = () => {
  Object.assign(settings, defaultSettings)
  ElMessage.info(t('settings.resetSuccess'))
}

// 第三方服务配置相关方法
const loadThirdPartySettings = async () => {
  servicesLoading.value = true
  try {
    const response = await axios.get('/api/admin/v1/setting/get', {
      params: {
        names: 'map_ak'
      },
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })
    if (response.data.code === 200) {
      // 批量获取返回的是map格式，需要从map中取出third_party_config
      const configData = response.data.data['third_party_config'] || {}
      Object.assign(thirdPartySettings, configData)
    } else {
      throw new Error(response.data.message)
    }
  } catch (error) {
    console.error('加载第三方服务配置失败:', error)
    ElMessage.error(t('settings.loadFailed') + ': ' + error.message)
  } finally {
    servicesLoading.value = false
  }
}

const saveThirdPartySettings = async () => {
  serviceSaving.value = true
  try {
    const response = await axios.post('/api/admin/v1/setting/set', {
      name: 'map_ak',
      value: thirdPartySettings
    }, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })
    if (response.data.code === 200) {
      ElMessage.success(t('settings.saveSuccess'))
    } else {
      throw new Error(response.data.message)
    }
  } catch (error) {
    ElMessage.error(t('settings.saveFailed') + ': ' + error.message)
  } finally {
    serviceSaving.value = false
  }
}

const defaultThirdPartySettings = {
  map_ak: ''
}

const resetThirdPartySettings = () => {
  Object.assign(thirdPartySettings, defaultThirdPartySettings)
  ElMessage.info(t('settings.resetSuccess'))
}

onMounted(() => {
  loadSettings()
  loadThirdPartySettings()
})
</script>

<style scoped>
.settings-container {
  padding: 20px;
}

/* 标签栏样式 */
.tab-header {
  display: flex;
  border-bottom: 2px solid #e4e7ed;
  margin-bottom: 20px;
}

.tab-item {
  padding: 12px 24px;
  cursor: pointer;
  font-size: 14px;
  color: #606266;
  border-bottom: 2px solid transparent;
  margin-bottom: -2px;
  transition: all 0.3s;
  position: relative;
}

.tab-item:hover {
  color: #409eff;
}

.tab-item.active {
  color: #409eff;
  border-bottom-color: #409eff;
  font-weight: 500;
}

.tab-item.active::after {
  content: '';
  position: absolute;
  bottom: -2px;
  left: 0;
  right: 0;
  height: 2px;
  background-color: #409eff;
}

/* 内容区域样式 */
.tab-content {
  min-height: 400px;
}

.tab-pane {
  animation: fadeIn 0.3s ease-in-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.settings-card {
  max-width: 800px;
  margin: 0 auto;
}

.avatar-uploader {
  :deep(.el-upload) {
    border: 1px dashed var(--el-border-color);
    border-radius: 6px;
    cursor: pointer;
    position: relative;
    overflow: hidden;
    transition: var(--el-transition-duration-fast);
  }

  :deep(.el-upload:hover) {
    border-color: var(--el-color-primary);
  }
}

.avatar-uploader-icon {
  font-size: 28px;
  color: #8c939d;
  width: 100px;
  height: 100px;
  text-align: center;
  line-height: 100px;
}

.avatar {
  width: 100px;
  height: 100px;
  display: block;
}

.avatar-uploader-loading {
  width: 100px;
  height: 100px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 28px;
  color: var(--el-color-primary);
}
</style>
