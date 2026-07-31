export interface UserVO {
    id: number;
    username: string;
    phone?: string;
    email?: string;
    status: 'NORMAL' | 'LOCKED' | 'BANNED';
    roles: string;
    lastLoginAt?: string;
    lastLoginIp?: string;
    createdAt: string;
}