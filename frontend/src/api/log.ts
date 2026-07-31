import request from './request'
import type { ApiResponse } from '@/models/response'

export function listLatestLogs(lines = 200, file?: string): Promise<ApiResponse<string[]>> {
  return request.get('/api/admin/logs', { params: { lines, file } })
}

export function listLogFiles(): Promise<ApiResponse<string[]>> {
  return request.get('/api/admin/logs/files')
}
