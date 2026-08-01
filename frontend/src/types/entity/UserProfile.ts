export interface UserProfile {
  userId: number;
  realName?: string;
  gender?: number;
  birthday?: string;
  avatar?: string;
  deleted: number;
  version: number;
  createdAt: string;
  updatedAt: string;
  createdBy?: string;
  updatedBy?: string;
}
