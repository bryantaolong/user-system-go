import request from './request'

export function exportAllUsers(fileName?: string): Promise<void> {
  return request({
    url: '/api/users/export',
    method: 'GET',
    params: { fileName: fileName || '用户数据' },
    responseType: 'blob',
  }).then((response) => {
    const filename = fileName ? `${fileName}.xlsx` : '用户数据.xlsx'
    const url = window.URL.createObjectURL(new Blob([response.data]))
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', filename)
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
  }) as unknown as Promise<void>
}
