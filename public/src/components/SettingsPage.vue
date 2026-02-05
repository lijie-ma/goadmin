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
        <SystemSettingsTab ref="systemTabRef" />
      </div>

      <!-- 第三方服务配置内容 -->
      <div v-show="activeTab === 'thirdParty'" class="tab-pane">
        <ThirdPartySettingsTab ref="thirdPartyTabRef" />
      </div>

      <!-- 服务配置内容 -->
      <div v-show="activeTab === 'service'" class="tab-pane">
        <ServiceSettingsTab ref="serviceTabRef" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import SystemSettingsTab from './settings/SystemSettingsTab.vue'
import ThirdPartySettingsTab from './settings/ThirdPartySettingsTab.vue'
import ServiceSettingsTab from './settings/ServiceSettingsTab.vue'
import '../assets/styles/settingsPage.css'

const { t } = useI18n()

// 标签页配置 - 使用computed实现多语言适配
const tabs = computed(() => [
  { key: 'system', label: t('settings.tabSystem') },
  { key: 'thirdParty', label: t('settings.tabThirdParty') },
  { key: 'service', label: t('settings.tabService') }
])

// 当前激活的标签
const activeTab = ref('system')

// 子组件引用
const systemTabRef = ref(null)
const thirdPartyTabRef = ref(null)
const serviceTabRef = ref(null)

// 切换标签
const switchTab = (key) => {
  activeTab.value = key
}

onMounted(() => {
  // 加载各个 tab 的数据
  if (systemTabRef.value?.loadSettings) {
    systemTabRef.value.loadSettings()
  }
  if (thirdPartyTabRef.value?.loadSettings) {
    thirdPartyTabRef.value.loadSettings()
  }
  if (serviceTabRef.value?.loadSettings) {
    serviceTabRef.value.loadSettings()
  }
})
</script>
