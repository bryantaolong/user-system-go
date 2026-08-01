export interface UserCreateRequest {
  username: string;
  password: string;
  phone?: string;
  email?: string;
  roleIds?: number[];
}
