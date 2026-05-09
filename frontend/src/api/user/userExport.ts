// src/api/userExport.ts
import request from '@/utils/request'

/**
 * 用户数据导出 API
 */

/**
 * 导出所有用户数据（包含所有字段）
 * @param fileName 导出文件名
 * @param status 状态过滤（可选）
 */
export function exportAllUsers(fileName?: string, status?: number | null): Promise<void> {
  return request({
    url: '/api/users/export',
    method: 'GET',
    params: {
      fileName: fileName || '用户数据',
      status: status ?? undefined
    },
    responseType: 'blob'
  }).then(response => {
    const filename = fileName ? `${fileName}.xlsx` : '用户数据.xlsx'
    const url = window.URL.createObjectURL(new Blob([response.data]))
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', filename)
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
  }) as Promise<void>
}

