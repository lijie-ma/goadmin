<template>
  <div id="app">
    <router-view />
    <!-- 语言切换按钮 -->
    <div class="language-switch" v-if="$route.path.includes('/login')">
      <el-dropdown @command="handleLanguageChange">
        <span class="el-dropdown-link">
          {{ currentLanguage === 'zh' ? '中文' : 'English' }}
          <el-icon><arrow-down /></el-icon>
        </span>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="zh">中文</el-dropdown-item>
            <el-dropdown-item command="en">English</el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </div>
</template>

<script>
import { useI18n } from 'vue-i18n'
import { useRouter, useRoute } from 'vue-router'
import { ArrowDown } from '@element-plus/icons-vue'
import { computed } from 'vue'

export default {
  name: 'App',
  components: {
    ArrowDown
  },
  setup() {
    const i18n = useI18n()
    const router = useRouter()
    const route = useRoute()

    // 获取当前语言
    const currentLanguage = computed(() => {
      return i18n.locale.value
    })

    // 切换语言
    const handleLanguageChange = (lang) => {
      // 更新本地存储
      localStorage.setItem('language', lang)

      // 更新 i18n locale
      i18n.locale.value = lang

      // 更新 Element Plus locale
      const elementLocale = lang === 'zh' ? 'zh-cn' : 'en'
      app.config.globalProperties.$ELEMENT.locale = elementLocale

      // 更新路由
      const currentPath = route.fullPath
      const newPath = currentPath.replace(/^\/(zh|en)/, `/${lang}`)
      router.push(newPath)

      // 更新 HTML lang 属性
      document.querySelector('html').setAttribute('lang', lang)
    }

    return {
      currentLanguage,
      handleLanguageChange
    }
  }
}
</script>

<style>
.language-switch {
  position: fixed;
  top: 20px;
  right: 20px;
  z-index: 1000;
}

.el-dropdown-link {
  cursor: pointer;
  color: var(--el-color-primary);
  display: flex;
  align-items: center;
  padding: 8px 16px;
  border-radius: 4px;
  background-color: white;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.el-dropdown-link:hover {
  background-color: var(--el-color-primary-light-9);
}
</style>
