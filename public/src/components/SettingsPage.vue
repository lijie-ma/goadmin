<template>
  <div class="settings-container">
    <el-card class="settings-card" v-loading="loading">
      <template #header>
        <div class="card-header">
          <h2 class="title">系统设置</h2>
        </div>
      </template>

      <el-form :model="settings" label-width="120px">
        <!-- 基本设置 -->
        <el-divider content-position="left">基本设置</el-divider>
        <el-form-item label="系统名称">
          <el-input v-model="settings.systemName" placeholder="请输入系统名称" />
        </el-form-item>
        <el-form-item label="系统Logo">
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
        <el-divider content-position="left">系统配置</el-divider>
        <el-form-item label="系统语言">
          <el-select v-model="settings.language" placeholder="请选择系统语言">
            <el-option label="简体中文" value="zh_CN" />
            <el-option label="English" value="en_US" />
          </el-select>
        </el-form-item>

        <!-- 主题设置 -->
        <el-divider content-position="left">主题设置</el-divider>
        <el-form-item label="主题色">
          <el-color-picker v-model="settings.theme.primaryColor" />
        </el-form-item>
        <el-form-item label="导航模式">
          <el-radio-group v-model="settings.theme.navMode">
            <el-radio label="sidebar">侧边栏</el-radio>
            <el-radio label="top">顶部导航</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="暗黑模式">
          <el-switch v-model="settings.theme.darkMode" />
        </el-form-item>

        <!-- 安全设置 -->
        <el-divider content-position="left">安全设置</el-divider>
        <el-form-item label="登录验证码">
          <el-switch v-model="settings.security.captchaEnabled" />
        </el-form-item>
        <el-form-item label="密码强度">
          <el-select v-model="settings.security.passwordStrength" placeholder="请选择密码强度要求">
            <el-option label="低" value="low" />
            <el-option label="中" value="medium" />
            <el-option label="高" value="high" />
          </el-select>
        </el-form-item>

        <!-- 保存按钮 -->
        <el-form-item>
          <el-button type="primary" @click="saveSettings" :loading="saving">保存设置</el-button>
          <el-button @click="resetSettings">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import axios from 'axios'

const loading = ref(false)
const saving = ref(false)

const settings = reactive({
  systemName: '后台管理系统',
  logo: '',
  theme: {
    primaryColor: '#409EFF',
    navMode: 'sidebar',
    darkMode: false
  },
  language: 'zh_CN',
  timezone: 'Asia/Shanghai',
  security: {
    captchaEnabled: true,
    passwordStrength: 'medium'
  }
})

const beforeLogoUpload = (file) => {
  const isImage = file.type.startsWith('image/')
  const isLt2M = file.size / 1024 / 1024 < 2

  if (!isImage) {
    ElMessage.error('上传文件只能是图片格式!')
    return false
  }
  if (!isLt2M) {
    ElMessage.error('上传图片大小不能超过 2MB!')
    return false
  }
  // 这里应该调用实际的上传API，这里只是演示
  settings.logo = URL.createObjectURL(file)
  return false
}

const loadSettings = async () => {
  loading.value = true
  try {
    const response = await axios.get('/api/admin/v1/setting/settings', {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })
    if (response.data.code === 200) {
      Object.assign(settings, response.data.data)
      ElMessage.success('设置加载成功')
    } else {
      throw new Error(response.data.message)
    }
  } catch (error) {
    ElMessage.error(`加载设置失败: ${error.message}`)
  } finally {
    loading.value = false
  }
}

const saveSettings = async () => {
  saving.value = true
  try {
    const response = await axios.put('/api/admin/v1/setting/settings', settings, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })
    if (response.data.code === 200) {
      ElMessage.success('设置保存成功')
    } else {
      throw new Error(response.data.message)
    }
  } catch (error) {
    ElMessage.error(`保存设置失败: ${error.message}`)
  } finally {
    saving.value = false
  }
}

const defaultSettings = {
  systemName: '后台管理系统',
  logo: '',
  theme: {
    primaryColor: '#409EFF',
    navMode: 'sidebar',
    darkMode: false
  },
  language: 'zh_CN',
  timezone: 'Asia/Shanghai',
  security: {
    captchaEnabled: true,
    passwordStrength: 'medium'
  }
}

const resetSettings = () => {
  Object.assign(settings, defaultSettings)
  ElMessage.info('已重置为默认设置')
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
