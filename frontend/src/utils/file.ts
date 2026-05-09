/**
 * 获取完整的头像访问地址
 * @param path 后端返回的相对路径
 * @returns 完整的 URL
 */
export const getAvatarUrl = (path?: string) => {
  if (!path) return ''
  // 如果是完整 URL，直接返回
  if (path.startsWith('http')) return path
  // 否则拼接 /uploads 前缀
  return `/uploads/${path}`
}
