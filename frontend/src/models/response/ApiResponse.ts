// shared/src/models/response/ApiResponse.ts
/**
 * @module models/response/ApiResponse
 * @description 通用后端API响应数据结构。
 * @template T 表示实际业务数据data的类型
 */
export interface ApiResponse<T = any> {
    code: number;    // 业务状态码
    message: string; // 业务消息
    data: T;         // 实际的业务数据
}
