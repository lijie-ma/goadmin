import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useUserStore } from '@/stores/user'

export const useMenuStore = defineStore('menu', () => {
  const userStore = useUserStore()

  // 权限代码到路径的映射
  const permissionPathMap = {
    'user_list': '/users',
    'user_create': '/users',
    'user_update': '/users',
    'user_delete': '/users',
    'user_resetPwd': '/users',
    'role_list': '/roles',
    'role_create': '/roles',
    'role_update': '/roles',
    'role_delete': '/roles',
    'role_active': '/roles',
    'role_all': '/roles',
    'role_info': '/roles',
    'role_perm_set': '/roles',
    'role_perm_info': '/roles',
    'server_get': '/settings',
    'server_set': '/settings',
    'server_batch_get': '/settings',
    'captcha_get': '/settings',
    'captcha_set': '/settings'
  }

  // 菜单配置
  const menuConfig = ref([
    {
      path: '/dashboard',
      name: 'dashboard',
      icon: 'Monitor',
      title: 'dashboard',
      permissionCode: null // 仪表板默认显示，不需要特定权限
    },
    {
      path: '/users',
      name: 'users',
      icon: 'User',
      title: 'userManagement',
      permissionCode: 'user_list'
    },
    {
      path: '/roles',
      name: 'roles',
      icon: 'Lock',
      title: 'roleManagement',
      permissionCode: 'role_list'
    },
    {
      path: '/settings',
      name: 'settings',
      icon: 'Setting',
      title: 'systemSettings',
      permissionCode: 'server_get'
    }
  ])

  // 根据权限代码获取对应的路径
  const getPathByPermissionCode = (permissionCode) => {
    return permissionPathMap[permissionCode]
  }

  // 根据路径获取需要的权限代码
  const getPermissionCodeByPath = (path) => {
    const menu = menuConfig.value.find(m => m.path === path)
    return menu ? menu.permissionCode : null
  }

  // 根据权限过滤菜单
  const filteredMenus = computed(() => {
    return menuConfig.value.filter(menu => {
      // 如果菜单没有配置权限代码，默认显示
      if (!menu.permissionCode) {
        return true
      }
      // 检查用户是否有该权限代码
      return userStore.hasPermissionCode(menu.permissionCode)
    })
  })

  return {
    menuConfig,
    filteredMenus,
    permissionPathMap,
    getPathByPermissionCode,
    getPermissionCodeByPath
  }
})