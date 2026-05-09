export interface SysUser {
  id: number;
  username: string;
  password?: string;
  phone?: string;
  email?: string;
  status: 'NORMAL' | 'LOCKED' | 'BANNED';
  roles: string;
  lastLoginAt?: string;
  lastLoginIp?: string;
  lastLoginDevice?: string;
  passwordResetAt?: string;
  loginFailCount?: number;
  lockedAt?: string;
  deleted: number;
  version: number;
  createdAt: string;
  updatedAt: string;
  createdBy: string;
  updatedBy: string;
}
