import api from '$lib/api/api';
import { promiseWrapper, type Result } from '$lib/api/promise-wrapper';
import type { SuccessModel } from '$lib/types/common..type';
import type {
	AdminCreateUserRequest,
	AdminUpdateUserRequest,
	User,
	UserRequest,
	UserUpdateRequest,
} from '$lib/types/user.type';

export class Users {
	static register(body: UserRequest) {
		return promiseWrapper(
			api.post('/users/register', body, { withCredentials: false })
		);
	}

	static async getList(): Promise<Result<User[]>> {
		const [err, res] = await promiseWrapper<SuccessModel<User[]>>(
			api.get('/users')
		);
		if (err) return [err, null];
		return [null, res.data];
	}

	static async get(username: string): Promise<Result<User>> {
		const [err, res] = await promiseWrapper<SuccessModel<User>>(
			api.get(`/users/${username}`)
		);
		if (err) return [err, null];
		return [null, res.data];
	}

	static async delete(username: string) {
		return promiseWrapper<void>(api.delete(`/users/${username}`));
	}

	static async update(
		username: string,
		body: UserUpdateRequest
	): Promise<Result<User>> {
		const [err, res] = await promiseWrapper<SuccessModel<User>>(
			api.patch(`/users/${username}`, body)
		);
		if (err) return [err, null];
		return [null, res.data];
	}

	static async adminCreate(
		body: AdminCreateUserRequest
	): Promise<Result<User>> {
		const [err, res] = await promiseWrapper<SuccessModel<User>>(
			api.post('/users/create', body)
		);
		if (err) return [err, null];
		return [null, res.data];
	}

	static async adminUpdate(
		username: string,
		body: AdminUpdateUserRequest
	): Promise<Result<User>> {
		const [err, res] = await promiseWrapper<SuccessModel<User>>(
			api.patch(`/users/${username}/admin`, body)
		);
		if (err) return [err, null];
		return [null, res.data];
	}

	static async requestPasswordReset(email: string) {
		return promiseWrapper(
			api.post(
				'/users/request-password-reset',
				{ email },
				{ withCredentials: false }
			)
		);
	}

	static async resetPassword(token: string, password: string) {
		return promiseWrapper(
			api.post(
				'/users/reset-password',
				{ token, password },
				{ withCredentials: false }
			)
		);
	}
}
