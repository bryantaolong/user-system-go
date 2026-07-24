/**
 * 统一 API 响应封装
 * @template T 实际业务数据的类型
 */
export interface ApiResponse<T = unknown> {
  code: number
  message: string
  data: T
}
