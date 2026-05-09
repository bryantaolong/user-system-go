import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { UserVO } from '@/models/vo'
import type { UserProfileVO } from '@/models/vo'
import * as authApi  from '@/api/auth/auth'
import * as userProfileApi from '@/api/user/userProfile'

export const useUserStore = defineStore('user', () => {
  const token = ref<string>(localStorage.getItem('token') || '')
  const userInfo = ref<UserVO | null>(null)
  const userProfile = ref<UserProfileVO | null>(null)

  const isAuthenticated = computed(() => !!token.value)
  const isAdmin = computed(() => {
    return userInfo.value?.roles.includes('ROLE_ADMIN') || false
  })

  /**
   * 设置Token
   */
  const setToken = (newToken: string) => {
    token.value = newToken
    localStorage.setItem('token', newToken)
  }

  /**
   * 清除Token
   */
  const clearToken = () => {
    token.value = ''
    localStorage.removeItem('token')
  }

  /**
   * 登录
   */
  const login = async (username: string, password: string) => {
    try {
      const res = await authApi.login({ username, password })
      if (res.code === 200 && res.data) {
        setToken(res.data)
        await fetchUserInfo()
        return { success: true }
      }
      return { success: false, message: res.message }
    } catch (error: any) {
      return { success: false, message: error.message || '登录失败' }
    }
  }

  /**
   * 注册
   */
  const register = async (data: { username: string; password: string; phone?: string; email?: string }) => {
    try {
      const res = await authApi.register(data)
      if (res.code === 200 && res.data) {
        userInfo.value = res.data
        return { success: true }
      }
      return { success: false, message: res.message }
    } catch (error: any) {
      return { success: false, message: error.message || '注册失败' }
    }
  }

  /**
   * 获取用户信息
   * 注意：两个 API 独立调用，即使 profile 获取失败也认为登录有效
   */
  const fetchUserInfo = async () => {
    try {
      const userRes = await authApi.getCurrentUser()

      if (userRes.code !== 200) {
        return { success: false, message: '获取用户信息失败' }
      }

      userInfo.value = userRes.data

      // UserProfile 独立获取，失败不影响登录状态
      try {
        const profileRes = await userProfileApi.getCurrentUserProfile()
        if (profileRes.code === 200) {
          userProfile.value = profileRes.data
        }
      } catch (profileError) {
        console.warn('获取用户资料失败，可能用户资料尚未创建:', profileError)
      }

      return { success: true }
    } catch (error: any) {
      return { success: false, message: error.message || '获取用户信息失败' }
    }
  }

  /**
   * 退出登录
   */
  const logout = async () => {
    try {
      await authApi.logout()
    } catch (error) {
      console.error('Logout error:', error)
    } finally {
      clearToken()
      userInfo.value = null
      userProfile.value = null
    }
  }

  /**
   * 修改密码
   */
  const changePassword = async (oldPassword: string, newPassword: string) => {
    try {
      const res = await authApi.changePassword({ oldPassword, newPassword })
      if (res.code === 200) {
        return { success: true }
      }
      return { success: false, message: res.message }
    } catch (error: any) {
      return { success: false, message: error.message || '修改密码失败' }
    }
  }

  /**
   * 更新用户资料
   */
  const updateProfile = async (data: any) => {
    try {
      const res = await userProfileApi.updateUserProfile(data)
      if (res.code === 200) {
        userProfile.value = res.data
        return { success: true }
      }
      return { success: false, message: res.message }
    } catch (error: any) {
      return { success: false, message: error.message || '更新资料失败' }
    }
  }

  /**
   * 注销账号
   */
  const deleteAccount = async () => {
    try {
      const res = await authApi.deleteAccount()
      if (res.code === 200) {
        await logout()
        return { success: true }
      }
      return { success: false, message: res.message }
    } catch (error: any) {
      return { success: false, message: error.message || '注销账号失败' }
    }
  }

  return {
    token,
    userInfo,
    userProfile,
    isAuthenticated,
    isAdmin,
    setToken,
    clearToken,
    login,
    register,
    fetchUserInfo,
    logout,
    changePassword,
    updateProfile,
    deleteAccount
  }
}, {
  persist: {
    key: 'user-store',
    storage: localStorage
  }
})
