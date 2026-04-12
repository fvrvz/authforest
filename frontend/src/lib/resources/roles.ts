import api from '$lib/api/api';
import { promiseWrapper, type Result } from '$lib/api/promise-wrapper';
import type { SuccessModel } from '$lib/types/common..type';
import type {
	CreateRoleRequest,
	Role,
	UpdateRoleRequest,
} from '$lib/types/role.type';

export class Roles {
	static async getList(): Promise<Result<Role[]>> {
		const [err, res] = await promiseWrapper<SuccessModel<Role[]>>(
			api.get('/roles')
		);
		if (err) return [err, null];
		return [null, res.data];
	}

	static async get(id: string): Promise<Result<Role>> {
		const [err, res] = await promiseWrapper<SuccessModel<Role>>(
			api.get(`/roles/${id}`)
		);
		if (err) return [err, null];
		return [null, res.data];
	}

	static async create(body: CreateRoleRequest): Promise<Result<Role>> {
		const [err, res] = await promiseWrapper<SuccessModel<Role>>(
			api.post('/roles', body)
		);
		if (err) return [err, null];
		return [null, res.data];
	}

	static async update(
		id: string,
		body: UpdateRoleRequest
	): Promise<Result<Role>> {
		const [err, res] = await promiseWrapper<SuccessModel<Role>>(
			api.patch(`/roles/${id}`, body)
		);
		if (err) return [err, null];
		return [null, res.data];
	}

	static async delete(id: string) {
		return promiseWrapper<void>(api.delete(`/roles/${id}`));
	}
}
