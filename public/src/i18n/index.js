import { createI18n } from 'vue-i18n'
import zhCN from '../locales/zh.json'
import enUS from '../locales/en.json'

// 获取浏览器语言
function getLanguage() {
  // 首先从 localStorage 获取
  const savedLang = localStorage.getItem('language')
  if (savedLang) {
    return savedLang
  }

  // 从URL路径获取语言
  const pathLang = window.location.pathname.split('/')[1]
  if (['zh', 'en'].includes(pathLang)) {
    return pathLang
  }

  // 使用浏览器语言
  const browserLang = navigator.language.toLowerCase()
  if (browserLang.startsWith('zh')) {
    return 'zh'
  }
  return 'en'
}

const i18n = createI18n({
  legacy: false, // 使用 Composition API 模式
  locale: getLanguage(),
  fallbackLocale: 'en',
  messages: {
    zh: zhCN,
    en: enUS
  }
})

// 设置语言函数
export function setLanguage(lang) {
  i18n.global.locale.value = lang
  localStorage.setItem('language', lang)
  document.querySelector('html').setAttribute('lang', lang)
}

// 获取当前语言
export function getLocale() {
  return i18n.global.locale.value
}

export default i18n
