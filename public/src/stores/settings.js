import { ref } from 'vue'

// 创建响应式的系统设置状态
export const systemName = ref('管理系统')
export const needCaptcha = ref(false)
export const systemSettings = ref({})

// 获取系统设置
export async function fetchSystemSettings() {
  try {
    const response = await fetch('/api/admin/v1/setting/get_settings', {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json'
      }
    })

    if (response.status === 200) {
      const result = await response.json()
      systemSettings.value = result.data || {}

      // 更新系统名称
      if (result.data && result.data.app_name) {
        systemName.value = result.data.app_name
      }

      // 更新验证码状态
      if (result.data && typeof result.data.admin === 'number') {
        needCaptcha.value = result.data.admin === 1
      }
    }
  } catch (error) {
    console.error('获取系统设置失败:', error)
    // 出错时使用默认值
    needCaptcha.value = true // 默认开启验证码
  }
}
