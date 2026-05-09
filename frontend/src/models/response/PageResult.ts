// src/models/response/PageResult.ts

export interface PageResult<T> {
    rows: T[];
    total: number;
    pageNum: number;
    pageSize: number;
    pages: number;
}