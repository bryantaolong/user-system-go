import request from '@/utils/request'
import type { ApiResponse } from '@/models/response'
import type { UserVO } from '@/models/vo'
import type { LoginRequest } from '@/models/request/auth'
import type { RegisterRequest } from '@/models/request/auth'
import type { ChangePasswordRequest } from '@/models/request/user'

/**
 * 用户认证API
 */

/**
 * 用户注册
 */
export function register(data: RegisterRequest): Promise<ApiResponse<UserVO>> {
    return request.post('/api/auth/register', data)
}

/**
 * 用户登录
 */
export function login(data: LoginRequest): Promise<ApiResponse<string>> {
    return request.post('/api/auth/login', data)
}

/**
 * 获取当前用户信息
 */
export function getCurrentUser(): Promise<ApiResponse<UserVO>> {
    return request.get('/api/auth/me')
}

/**
 * 验证Token
 */
export function validate(token: string): Promise<ApiResponse<string>> {
    return request.get('/api/auth/validate', {params: {token}})
}

/**
 * 修改密码
 */
export function changePassword(data: ChangePasswordRequest): Promise<ApiResponse<UserVO>> {
    return request.put('/api/auth/password', data)
}

/**
 * 注销账号
 */
export function deleteAccount(): Promise<ApiResponse<UserVO>> {
    return request.delete('/api/auth')
}

/**
 * 退出登录
 */
export function logout(): Promise<boolean> {
    return request.get('/api/auth/logout')
}
