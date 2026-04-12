import type { userSchema } from '$lib/schemas/user.schema';
import type z from 'zod';
import type { BaseModel } from './common..type';
import type { Role } from './role.type';

export type User = z.infer<typeof userSchema> &
	BaseModel & { fullName?: string; roles?: Role[] };

export type UserRequest = Omit<User, keyof BaseModel | 'fullName' | 'roles'>;

export type UserUpdateRequest = Partial<
	Omit<UserRequest, 'password' | 'username'>
>;

export interface AdminCreateUserRequest {
	username: string;
	email: string;
	password: string;
	firstName: string;
	lastName: string;
	DOB: string;
	roleIds: string[];
}

export interface AdminUpdateUserRequest {
	email?: string;
	firstName?: string;
	lastName?: string;
	DOB?: string;
	roleIds?: string[];
}
