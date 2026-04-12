export interface Role {
  id: string;
  name: string;
  description: string;
}

export interface CreateRoleRequest {
  name: string;
  description: string;
}

export type UpdateRoleRequest = Partial<CreateRoleRequest>;
