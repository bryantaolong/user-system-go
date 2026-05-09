import request from '@/utils/request'
import type { Result } from '@/models/response'

/**
 * 系统日志监控 API
 */

/**
 * 获取最新系统日志
 * @param lines 返回的最大行数，默认 200
 * @param file 日志文件名，可选
 */
export function listLatestLogs(lines = 200, file?: string): Promise<Result<string[]>> {
  return request.get('/api/admin/logs', {
    params: { lines, file }
  })
}

/**
 * 获取可用的日志文件列表
 */
export function listLogFiles(): Promise<Result<string[]>> {
  return request.get('/api/admin/logs/files')
}
