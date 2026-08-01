import request from './request'
import type { ApiResponse, PageResponse } from '@/types/response'
import type { SysUser } from '@/types/entity'
import type { UserCreateRequest, UserUpdateRequest, ChangePasswordRequest, UserSearchRequest } from '@/types/request/user'

export function createUser(data: UserCreateRequest): Promise<ApiResponse<SysUser>> {
  return request.post('/api/users', data)
}

export function listUsers(pageNum = 1, pageSize = 10): Promise<ApiResponse<PageResponse<SysUser>>> {
  return request.get('/api/users', { params: { pageNum, pageSize } })
}

export function getUserById(userId: number): Promise<ApiResponse<SysUser>> {
  return request.get(`/api/users/${userId}`)
}

export function getUserByUsername(username: string): Promise<ApiResponse<SysUser>> {
  return request.get(`/api/users/username/${username}`)
}

export function queryUsers(data: UserSearchRequest, pageNum = 1, pageSize = 10): Promise<ApiResponse<PageResponse<SysUser>>> {
  return request.post('/api/users/search', data, { params: { pageNum, pageSize } })
}

export function updateUser(userId: number, data: UserUpdateRequest): Promise<ApiResponse<SysUser>> {
  return request.put(`/api/users/${userId}`, data)
}

export function changeUserRoles(userId: number, roleIds: number[]): Promise<ApiResponse<SysUser>> {
  return request.put(`/api/users/roles/${userId}`, { roleIds })
}

export function resetPassword(userId: number, newPassword: string): Promise<ApiResponse<SysUser>> {
  const data: ChangePasswordRequest = { newPassword }
  return request.put(`/api/users/password/${userId}`, data)
}

export function blockUser(userId: number): Promise<ApiResponse<SysUser>> {
  return request.put(`/api/users/block/${userId}`)
}

export function unblockUser(userId: number): Promise<ApiResponse<SysUser>> {
  return request.put(`/api/users/unblock/${userId}`)
}

export function deleteUser(userId: number): Promise<ApiResponse<number>> {
  return request.delete(`/api/users/${userId}`)
}
