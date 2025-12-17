<template>
  <div class="welcome-card">
    <h2>{{ t('welcome', { name: systemName }) }}</h2>
    <p>{{ t('dashboardDesc') }}</p>

    <el-card style="margin-top: 20px;">
      <h3>{{ t('systemConfig') }}</h3>
      <p>{{ t('systemName') }}：{{ systemName }}</p>
      <p>{{ t('captchaStatus') }}：{{ needCaptcha ? t('enabled') : t('disabled') }}</p>
      <p>{{ t('allConfig') }}：{{ JSON.stringify(systemSettings, null, 2) }}</p>
    </el-card>
  </div>
</template>

<script setup>
import { systemName, needCaptcha, systemSettings } from '../stores/settings'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
</script>

<style scoped>
.welcome-card {
  padding: 20px;
}

.welcome-card h2 {
  margin-bottom: 10px;
  color: #303133;
}

.welcome-card p {
  color: #606266;
  margin-bottom: 10px;
}
</style>
