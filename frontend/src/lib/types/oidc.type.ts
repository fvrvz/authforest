export interface OAuthClient {
	client_id: string;
	client_secret?: string;
	client_name: string;
	client_type: 'public' | 'confidential';
	redirect_uris: string[];
	scopes: string;
	grant_types: string;
	access_token_expiry_minutes: number;
	refresh_token_expiry_hours: number;
	id_token_expiry_minutes: number;
	created_at?: string;
}

export type RegisterClientRequest = Omit<
	OAuthClient,
	| 'client_id'
	| 'client_secret'
	| 'redirect_uris'
	| 'created_at'
	| 'access_token_expiry_minutes'
	| 'refresh_token_expiry_hours'
	| 'id_token_expiry_minutes'
> & {
	redirect_uris: string;
	access_token_expiry_minutes?: number;
	refresh_token_expiry_hours?: number;
	id_token_expiry_minutes?: number;
};

export type UpdateClientRequest = Partial<RegisterClientRequest>;

export interface DashboardStats {
	total_users: number;
	total_clients: number;
}
