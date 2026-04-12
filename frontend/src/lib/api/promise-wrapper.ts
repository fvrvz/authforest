import type { AxiosError, AxiosResponse } from 'axios';

export type Result<T> = [AxiosError, null] | [null, T];

export async function promiseWrapper<T>(
	promise: Promise<AxiosResponse<T>>
): Promise<Result<T>> {
	try {
		const response = await promise;
		return [null, response.data];
	} catch (e) {
		return [e as AxiosError, null];
	}
}
