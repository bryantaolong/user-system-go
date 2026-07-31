import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import * as authApi from '@/api/auth'
import * as userProfileApi from '@/api/userProfile'
import * as userApi from '@/api/user'
import type { UserVO, UserProfileVO } from '@/models/vo'

export const useUserStore = defineStore(
  'user',
  () => {
    const token = ref(localStorage.getItem('token') || '')
    const user = ref<UserVO | null>(null)
    const profile = ref<UserProfileVO | null>(null)

    const isAuthenticated = computed(() => !!token.value)
    const isAdmin = computed(() => user.value?.roles?.includes('ROLE_ADMIN') || false)

    function setToken(newToken: string) {
      token.value = newToken
      localStorage.setItem('token', newToken)
    }

    function clearToken() {
      token.value = ''
      localStorage.removeItem('token')
    }

    async function login(username: string, password: string): Promise<{ success: boolean; message?: string }> {
      try {
        const res = await authApi.login({ username, password })
        if (res.code === 200 && res.data) {
          setToken(res.data.token)
          await fetchUserInfo()
          return { success: true }
        }
        return { success: false, message: res.message }
      } catch (error: any) {
        return { success: false, message: error.message || '登录失败' }
      }
    }

    async function register(data: { username: string; password: string; phone?: string; email?: string }): Promise<{ success: boolean; message?: string }> {
      try {
        const res = await authApi.register(data)
        if (res.code === 200 && res.data) {
          user.value = res.data
          return { success: true }
        }
        return { success: false, message: res.message }
      } catch (error: any) {
        return { success: false, message: error.message || '注册失败' }
      }
    }

    async function fetchUserInfo(): Promise<{ success: boolean; message?: string }> {
      try {
        const userRes = await authApi.getCurrentUser()
        if (userRes.code !== 200) {
          return { success: false, message: '获取用户信息失败' }
        }
        user.value = userRes.data

        try {
          const profileRes = await userProfileApi.getCurrentUserProfile()
          if (profileRes.code === 200) {
            profile.value = profileRes.data
          }
        } catch (profileError) {
          console.warn('获取用户资料失败:', profileError)
        }

        return { success: true }
      } catch (error: any) {
        return { success: false, message: error.message || '获取用户信息失败' }
      }
    }

    async function logout(): Promise<void> {
      try {
        await authApi.logout()
      } catch (error) {
        console.error('Logout error:', error)
      } finally {
        clearToken()
        user.value = null
        profile.value = null
      }
    }

    async function changePassword(oldPassword: string, newPassword: string): Promise<{ success: boolean; message?: string }> {
      try {
        const res = await authApi.changePassword({ oldPassword, newPassword })
        if (res.code === 200) {
          // 密码修改成功后强制登出，因为后端已清除 Redis Token
          await logout()
          return { success: true }
        }
        return { success: false, message: res.message }
      } catch (error: any) {
        return { success: false, message: error.message || '修改密码失败' }
      }
    }

    async function updateProfile(data: any): Promise<{ success: boolean; message?: string }> {
      try {
        const res = await userProfileApi.updateUserProfile(data)
        if (res.code === 200) {
          profile.value = res.data
          return { success: true }
        }
        return { success: false, message: res.message }
      } catch (error: any) {
        return { success: false, message: error.message || '更新资料失败' }
      }
    }

    async function deleteAccount(): Promise<{ success: boolean; message?: string }> {
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
      user,
      profile,
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
      deleteAccount,
    }
  },
  {
    persist: {
      key: 'user-store',
      paths: ['token'],
    },
  }
)
