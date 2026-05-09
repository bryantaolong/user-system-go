import request from '@/utils/request'
import type { Result } from '@/models/response'
import type { PageResult } from '@/models/response'
import type { SysUser } from '@/models/entity'
import type { UserCreateRequest } from '@/models/request/user'
import type { UserUpdateRequest } from '@/models/request/user'
import type { ChangePasswordRequest } from '@/models/request/user'
import type { UserSearchRequest } from '@/models/request/user'

/**
 * 用户管理API
 */

/**
 * 创建用户（管理员）
 */
export function createUser(data: UserCreateRequest): Promise<Result<SysUser>> {
  return request.post('/api/users', data)
}

/**
 * 获取用户列表（分页）
 */
export function listUsers(pageNum = 1, pageSize = 10): Promise<Result<PageResult<SysUser>>> {
  return request.get('/api/users', {
    params: { pageNum, pageSize }
  })
}

/**
 * 根据ID获取用户信息
 */
export function getUserById(userId: number): Promise<Result<SysUser>> {
  return request.get(`/api/users/${userId}`)
}

/**
 * 根据用户名获取用户信息
 */
export function getUserByUsername(username: string): Promise<Result<SysUser>> {
  return request.get(`/api/users/username/${username}`)
}

/**
 * 搜索用户
 */
export function queryUsers(data: UserSearchRequest, pageNum = 1, pageSize = 10): Promise<Result<PageResult<SysUser>>> {
  return request.post('/api/users/search', data, {
    params: { pageNum, pageSize }
  })
}

/**
 * 更新用户信息
 */
export function updateUser(userId: number, data: UserUpdateRequest): Promise<Result<SysUser>> {
  return request.put(`/api/users/${userId}`, data)
}

/**
 * 修改用户角色（管理员）
 */
export function changeUserRoles(userId: number, roleIds: number[]): Promise<Result<SysUser>> {
  return request.put(`/api/users/roles/${userId}`, { roleIds })
}

/**
 * 重置用户密码（管理员）
 */
export function resetPassword(userId: number, newPassword: string): Promise<Result<SysUser>> {
  const data: ChangePasswordRequest = { newPassword }
  return request.put(`/api/users/password/${userId}`, data)
}

/**
 * 封禁用户
 */
export function blockUser(userId: number): Promise<Result<SysUser>> {
  return request.put(`/api/users/block/${userId}`)
}

/**
 * 解封用户
 */
export function unblockUser(userId: number): Promise<Result<SysUser>> {
  return request.put(`/api/users/unblock/${userId}`)
}

/**
 * 删除用户（逻辑删除）
 */
export function deleteUser(userId: number): Promise<Result<number>> {
  return request.delete(`/api/users/${userId}`)
}

