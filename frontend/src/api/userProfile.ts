import request from './request'
import type { ApiResponse } from '@/models/response'
import type { UserUpdateRequest } from '@/models/request/user'
import type { UserProfileVO } from '@/models/vo'

export function uploadAvatar(file: File): Promise<ApiResponse<string>> {
  const allowedTypes = ['image/png', 'image/jpeg', 'image/gif', 'image/webp']
  const maxSize = 1 * 1024 * 1024 // 1MB

  if (!allowedTypes.includes(file.type)) {
    return Promise.reject(new Error('仅支持 PNG/JPEG/GIF/WebP 格式的头像'))
  }
  if (file.size > maxSize) {
    return Promise.reject(new Error('头像大小不能超过 1MB'))
  }

  const formData = new FormData()
  formData.append('file', file)
  return request.post('/api/user-profiles/avatar', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })
}

export function getUserProfileByUserId(userId: number): Promise<ApiResponse<UserProfileVO>> {
  return request.get(`/api/user-profiles/${userId}`)
}

export function getUserProfileByRealName(realName: string): Promise<ApiResponse<UserProfileVO>> {
  return request.get(`/api/user-profiles/name/${realName}`)
}

export function getCurrentUserProfile(): Promise<ApiResponse<UserProfileVO>> {
  return request.get('/api/user-profiles/me')
}

export function updateUserProfile(data: UserUpdateRequest): Promise<ApiResponse<UserProfileVO>> {
  return request.put('/api/user-profiles', data)
}
