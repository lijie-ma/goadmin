import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useUserStore } from '@/stores/user'
import { Document } from '@element-plus/icons-vue'

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
    'server_get': '/settings/system',
    'server_set': '/settings/system',
    'server_batch_get': '/settings/system',
    'captcha_get': '/settings/system',
    'captcha_set': '/settings/system',
    'position_list': '/settings/position',
    'position_create': '/settings/position',
    'position_update': '/settings/position',
    'position_delete': '/settings/position',
    'tenant_list': '/settings/tenant',
    'tenant_create': '/settings/tenant',
    'tenant_update': '/settings/tenant',
    'tenant_delete': '/settings/tenant',
    'tenant_info': '/settings/tenant',
    'operate_log': '/operate-logs'
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
      permissionCode: null,
      children: [
        {
          path: '/settings/system',
          name: 'systemSettings',
          icon: 'Setting',
          title: 'systemSettings',
          permissionCode: 'server_get'
        },
        {
          path: '/settings/position',
          name: 'positionManagement',
          icon: 'Location',
          title: 'positionManagement',
          permissionCode: 'position_list'
        },
        {
          path: '/settings/tenant',
          name: 'tenantManagement',
          icon: 'OfficeBuilding',
          title: 'tenantManagement',
          permissionCode: 'tenant_list'
        }
      ]
    },
    {
      path: '/operate-logs',
      name: 'operateLogs',
      icon: 'Document',
      title: 'operateLogs',
      permissionCode: 'operate_log'
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
    return menuConfig.value.map(menu => {
      // 如果有子菜单，递归过滤子菜单
      if (menu.children && menu.children.length > 0) {
        const filteredChildren = menu.children.filter(child => {
          // 如果子菜单没有配置权限代码，默认显示
          if (!child.permissionCode) {
            return true
          }
          // 检查用户是否有该权限代码
          return userStore.hasPermissionCode(child.permissionCode)
        })
        // 如果有可见的子菜单，返回带子菜单的父菜单
        if (filteredChildren.length > 0) {
          return {
            ...menu,
            children: filteredChildren
          }
        }
        // 如果没有可见的子菜单，不显示父菜单
        return null
      }
      // 如果菜单没有配置权限代码，默认显示
      if (!menu.permissionCode) {
        return menu
      }
      // 检查用户是否有该权限代码
      return userStore.hasPermissionCode(menu.permissionCode) ? menu : null
    }).filter(menu => menu !== null)
  })

  return {
    menuConfig,
    filteredMenus,
    permissionPathMap,
    getPathByPermissionCode,
    getPermissionCodeByPath
  }
})