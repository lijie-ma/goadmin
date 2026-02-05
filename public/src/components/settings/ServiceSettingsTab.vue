<template>
  <el-card class="settings-card" v-loading="loading">
    <el-form :model="settings" label-width="120px">
      <!-- 服务配置 -->
      <el-divider content-position="left">{{ t('settings.serviceConfig') }}</el-divider>
      <el-form-item :label="t('settings.region')">
        <el-input
          v-model="settings.region"
          :placeholder="t('settings.regionPlaceholder')"
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
import { useI18n } from 'vue-i18n'
import { useServiceSettings } from '../../composables/useSettings'

const { t } = useI18n()
const {
  settings,
  loading,
  saving,
  loadSettings,
  saveSettings,
  resetSettings
} = useServiceSettings()

// 暴露方法给父组件
defineExpose({
  loadSettings
})
</script>