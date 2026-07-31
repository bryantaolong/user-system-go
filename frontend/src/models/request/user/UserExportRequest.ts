/**
 * @module src/models/request/user/UserExportRequest.ts
 * @description 定义用户数据导出请求的数据结构。
 */
export interface UserExportRequest {
    fields?: string[]; // 要导出的字段列表
    fileName?: string; // 导出文件名
    status?: number; // 状态过滤
}