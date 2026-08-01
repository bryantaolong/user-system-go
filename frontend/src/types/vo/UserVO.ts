export interface UserVO {
    id: number;
    username: string;
    phone?: string;
    email?: string;
    status: 'NORMAL' | 'LOCKED' | 'BANNED' | 'UNKNOWN';
    roles: string;
    lastLoginAt?: string;
    lastLoginIp?: string;
    createdAt: string;
}