import request from './request'
import type { ApiResponse, UserRoleOptionVO } from '@/models/response'

export function listRoles(): Promise<ApiResponse<UserRoleOptionVO[]>> {
  return request.get('/api/user-roles')
}
