import { ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'
import axios from 'axios'

// 系统设置
export function useSystemSettings() {
  const { t } = useI18n()
  const loading = ref(false)
  const saving = ref(false)
  const uploading = ref(false)

  const settings = reactive({
    system_name: t('app.title'),
    logo: '',
    language: 'zh_CN',
    admin: 1
  })

  const defaultSettings = {
    system_name: t('app.title'),
    logo: '',
    language: 'zh_CN',
    admin: 1
  }

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

  const resetSettings = () => {
    Object.assign(settings, defaultSettings)
    ElMessage.info(t('settings.resetSuccess'))
  }

  return {
    settings,
    loading,
    saving,
    uploading,
    loadSettings,
    saveSettings,
    resetSettings,
    beforeLogoUpload
  }
}

// 第三方服务设置
export function useThirdPartySettings() {
  const { t } = useI18n()
  const loading = ref(false)
  const saving = ref(false)

  const settings = reactive({
    map_ak: '',
    map_scode: ''
  })

  const defaultSettings = {
    map_ak: '',
    map_scode: ''
  }

  const loadSettings = async () => {
    loading.value = true
    try {
      const response = await axios.get('/api/admin/v1/setting/get', {
        params: {
          names: 'map_config'
        },
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
      })
      if (response.data.code === 200) {
        const configData = response.data.data || {}
        if (configData.map_config) {
          try {
            const mapConfig = JSON.parse(configData.map_config)
            settings.map_ak = mapConfig.map_ak || ''
            settings.map_scode = mapConfig.map_scode || ''
          } catch (e) {
            console.error('解析 map_config 失败:', e)
          }
        }
      } else {
        throw new Error(response.data.message)
      }
    } catch (error) {
      console.error('加载第三方服务配置失败:', error)
      ElMessage.error(t('settings.loadFailed') + ': ' + error.message)
    } finally {
      loading.value = false
    }
  }

  const saveSettings = async () => {
    saving.value = true
    try {
      const mapConfig = {
        map_ak: settings.map_ak,
        map_scode: settings.map_scode
      }

      const response = await axios.post('/api/admin/v1/setting/set', {
        name: 'map_config',
        value: JSON.stringify(mapConfig)
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
      saving.value = false
    }
  }

  const resetSettings = () => {
    Object.assign(settings, defaultSettings)
    ElMessage.info(t('settings.resetSuccess'))
  }

  return {
    settings,
    loading,
    saving,
    loadSettings,
    saveSettings,
    resetSettings
  }
}

// 服务配置设置
export function useServiceSettings() {
  const { t } = useI18n()
  const loading = ref(false)
  const saving = ref(false)

  const settings = reactive({
    region: ''
  })

  const defaultSettings = {
    region: ''
  }

  const loadSettings = async () => {
    loading.value = true
    try {
      const response = await axios.get('/api/admin/v1/setting/get', {
        params: {
          names: 'service_region'
        },
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
      })
      if (response.data.code === 200) {
        const configData = response.data.data || {}
        if (configData.service_region) {
          try {
            const serviceConfig = JSON.parse(configData.service_region)
            settings.region = serviceConfig.region || ''
          } catch (e) {
            console.error('解析 service_region 失败:', e)
          }
        }
      } else {
        throw new Error(response.data.message)
      }
    } catch (error) {
      console.error('加载服务配置失败:', error)
      ElMessage.error(t('settings.loadFailed') + ': ' + error.message)
    } finally {
      loading.value = false
    }
  }

  const saveSettings = async () => {
    saving.value = true
    try {
      const serviceConfig = {
        region: settings.region
      }

      const response = await axios.post('/api/admin/v1/setting/set', {
        name: 'service_region',
        value: JSON.stringify(serviceConfig)
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
      saving.value = false
    }
  }

  const resetSettings = () => {
    Object.assign(settings, defaultSettings)
    ElMessage.info(t('settings.resetSuccess'))
  }

  return {
    settings,
    loading,
    saving,
    loadSettings,
    saveSettings,
    resetSettings
  }
}