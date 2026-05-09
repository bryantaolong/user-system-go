import request from '@/utils/request'
import type { Result } from '@/models/response'
import type { UserUpdateRequest } from '@/models/request/user'
import type { UserProfileVO } from '@/models/vo'

/**
 * 用户资料 API
 */

/**
 * 上传当前用户头像
 */
export function uploadAvatar(file: File): Promise<Result<string>> {
  const formData = new FormData()
  formData.append('file', file)
  return request.post('/api/user-profiles/avatar', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

/**
 * 根据用户主键查询用户资料（公开访问，用于展示用户信息）
 */
export function getUserProfileByUserId(userId: number): Promise<Result<UserProfileVO>> {
  return request.get(`/api/user-profiles/${userId}`)
}

/**
 * 根据真实姓名获取用户资料
 */
export function getUserProfileByRealName(realName: string): Promise<Result<UserProfileVO>> {
  return request.get(`/api/user-profiles/name/${realName}`)
}

/**
 * 获取当前登录用户的资料
 */
export function getCurrentUserProfile(): Promise<Result<UserProfileVO>> {
  return request.get('/api/user-profiles/me')
}

/**
 * 更新当前用户资料
 */
export function updateUserProfile(data: UserUpdateRequest): Promise<Result<UserProfileVO>> {
  return request.put('/api/user-profiles', data)
}
