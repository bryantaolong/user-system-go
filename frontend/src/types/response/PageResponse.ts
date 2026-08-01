// shared/src/models/response/PageResponse.ts

export interface PageResponse<T> {
    rows: T[];
    total: number;
    pageNum: number;
    pageSize: number;
    pages: number;
}
