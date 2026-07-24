import request from '@/utils/request'
import type { ApiResponse } from '@/models/response'
import type { UserRoleOptionVO } from '@/models/vo'

/**
 * 用户角色管理API
 */

/**
 * 获取全部角色下拉选项
 */
export function listRoles(): Promise<ApiResponse<UserRoleOptionVO[]>> {
  return request.get('/api/user-roles')
}
