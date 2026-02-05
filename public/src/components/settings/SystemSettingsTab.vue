<template>
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
</template>

<script setup>
import { reactive } from 'vue'
import { Plus, Loading } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'
import axios from 'axios'
import { useSystemSettings } from '../../composables/useSettings'

const { t } = useI18n()
const {
  settings,
  loading,
  saving,
  uploading,
  loadSettings,
  saveSettings,
  resetSettings,
  beforeLogoUpload
} = useSystemSettings()

// 暴露方法给父组件
defineExpose({
  loadSettings
})
</script>