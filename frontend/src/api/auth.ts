import request from './request'
import type { ApiResponse } from '@/models/response'
import type { UserVO } from '@/models/vo'
import type { LoginRequest, RegisterRequest, ChangePasswordRequest } from '@/models/request/auth'

export function register(data: RegisterRequest): Promise<ApiResponse<UserVO>> {
  return request.post('/api/auth/register', data)
}

export function login(data: LoginRequest): Promise<ApiResponse<{ token: string; user: UserVO }>> {
  return request.post('/api/auth/login', data)
}

export function getCurrentUser(): Promise<ApiResponse<UserVO>> {
  return request.get('/api/auth/me')
}

export function validateToken(token: string): Promise<ApiResponse<string>> {
  return request.get('/api/auth/validate', { params: { token } })
}

export function changePassword(data: ChangePasswordRequest): Promise<ApiResponse<UserVO>> {
  return request.put('/api/auth/password', data)
}

export function deleteAccount(): Promise<ApiResponse<any>> {
  return request.delete('/api/auth')
}

export function logout(): Promise<ApiResponse<boolean>> {
  return request.get('/api/auth/logout')
}
