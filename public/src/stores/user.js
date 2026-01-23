import { defineStore } from 'pinia'
import { ref } from 'vue'
import axios from 'axios'

export const useUserStore = defineStore('user', () => {
  // 用户信息
  const userInfo = ref({
    username: '',
    email: '',
    role_code: '',
    permission_codes: []
  })

  // 权限代码
  const permissionCodes = ref([])

  // 是否已登录
  const isLoggedIn = ref(false)

  // 获取用户信息
  const fetchUserInfo = async (locale = 'zh') => {
    try {
      const token = localStorage.getItem('token')
      if (!token) {
        isLoggedIn.value = false
        return false
      }

      const response = await axios.get('/api/admin/v1/user/info', {
        headers: {
          'Authorization': `Bearer ${token}`,
          'Accept-Language': locale === 'zh' ? 'zh-CN' : 'en-US'
        }
      })

      if (response.data.code === 200) {
        userInfo.value = response.data.data
        permissionCodes.value = response.data.data.permission_codes || []
        isLoggedIn.value = true
        // 更新 localStorage 中的用户信息
        localStorage.setItem('user', JSON.stringify(response.data.data))
        return true
      } else {
        isLoggedIn.value = false
        return false
      }
    } catch (error) {
      console.error('获取用户信息失败:', error)
      isLoggedIn.value = false
      // 如果获取失败，尝试从 localStorage 读取
      const storedUser = localStorage.getItem('user')
      if (storedUser) {
        try {
          const userData = JSON.parse(storedUser)
          userInfo.value = userData
          permissionCodes.value = userData.permission_codes || []
          isLoggedIn.value = true
          return true
        } catch (error) {
          console.error('解析用户信息失败:', error)
          isLoggedIn.value = false
          return false
        }
      }
      return false
    }
  }

  // 检查是否有权限（基于路径）
  const hasPermission = (path) => {
    // 如果权限代码为空，默认显示所有菜单
    if (!permissionCodes.value || permissionCodes.value.length === 0) {
      return true
    }
    // 检查权限代码中是否包含该路径
    return permissionCodes.value.includes(path)
  }

  // 检查是否有权限代码
  const hasPermissionCode = (permissionCode) => {
    // 如果权限代码为空，默认允许访问
    if (!permissionCodes.value || permissionCodes.value.length === 0) {
      return true
    }
    // 检查用户是否有该权限代码
    return permissionCodes.value.includes(permissionCode)
  }

  // 登出
  const logout = () => {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    userInfo.value = {
      username: '',
      email: '',
      role_code: '',
      permission_codes: []
    }
    permissionCodes.value = []
    isLoggedIn.value = false
  }

  // 从 localStorage 初始化用户信息
  const initFromLocalStorage = () => {
    const storedUser = localStorage.getItem('user')
    if (storedUser) {
      try {
        const userData = JSON.parse(storedUser)
        userInfo.value = userData
        permissionCodes.value = userData.permission_codes || []
        isLoggedIn.value = !!localStorage.getItem('token')
      } catch (error) {
        console.error('解析用户信息失败:', error)
        isLoggedIn.value = false
      }
    }
  }

  return {
    userInfo,
    permissionCodes,
    isLoggedIn,
    fetchUserInfo,
    hasPermission,
    hasPermissionCode,
    logout,
    initFromLocalStorage
  }
})