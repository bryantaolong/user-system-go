/**
 * 分页响应
 * @template T 列表项的类型
 */
export interface PageResponse<T> {
  rows: T[]
  total: number
  pageNum: number
  pageSize: number
  pages: number
}
