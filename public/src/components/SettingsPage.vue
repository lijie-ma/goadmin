<template>
  <div class="settings-container">
    <el-card class="settings-card" v-loading="loading">
      <template #header>
        <div class="card-header">
          <h2 class="title">{{ t('settings.title') }}</h2>
        </div>
      </template>

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
          >
            <img v-if="settings.logo" :src="settings.logo" class="avatar" />
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
          <el-button type="primary" @click="saveSettings" :loading="saving">{{ t('settings.saveSettings') }}</el-button>
          <el-button @click="resetSettings">{{ t('settings.reset') }}</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'
import axios from 'axios'

const { t } = useI18n()
const loading = ref(false)
const saving = ref(false)

const settings = reactive({
  system_name: t('app.title'),
  logo: '',
  language: 'zh_CN',
  admin: 1
})

const beforeLogoUpload = (file) => {
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
  // 这里应该调用实际的上传API，这里只是演示
  settings.logo = URL.createObjectURL(file)
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
      ElMessage.success(t('settings.loadSuccess'))
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

onMounted(() => {
  loadSettings()
})
</script>

<style scoped>
.settings-container {
  padding: 20px;
}

.settings-card {
  max-width: 800px;
  margin: 0 auto;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.title {
  font-size: 18px;
  font-weight: bold;
  margin: 0;
  color: #303133;
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
</style>
