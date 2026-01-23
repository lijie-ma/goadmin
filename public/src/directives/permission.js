import { useUserStore } from '@/stores/user'

/**
 * 权限指令
 * 用法：v-permission="'user_create'" 或 v-permission="['user_create', 'user_update']"
 * 当用户没有指定权限时，元素会被移除
 */
export default {
  mounted(el, binding) {
    const { value } = binding
    const userStore = useUserStore()

    if (value) {
      const hasPermission = checkPermission(value, userStore)
      if (!hasPermission) {
        // 移除元素
        el.parentNode && el.parentNode.removeChild(el)
      }
    } else {
      throw new Error('需要指定权限代码！例如：v-permission="\'user_create\'"')
    }
  },

  updated(el, binding) {
    const { value } = binding
    const userStore = useUserStore()

    if (value) {
      const hasPermission = checkPermission(value, userStore)
      if (!hasPermission) {
        // 移除元素
        el.parentNode && el.parentNode.removeChild(el)
      }
    }
  }
}

/**
 * 检查权限
 * @param {string|string[]} permission - 权限代码或权限代码数组
 * @param {Object} userStore - 用户store
 * @returns {boolean} 是否有权限
 */
function checkPermission(permission, userStore) {
  if (typeof permission === 'string') {
    return userStore.hasPermissionCode(permission)
  } else if (Array.isArray(permission)) {
    // 数组模式：只要有一个权限即可
    return permission.some(p => userStore.hasPermissionCode(p))
  }
  return false
}