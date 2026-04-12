import z from 'zod/v4';

export const registerClientSchema = z.object({
	client_name: z.string().trim().min(1, 'Client name is required'),
	client_type: z.enum(['public', 'confidential']),
	redirect_uris: z
		.string()
		.trim()
		.min(1, 'At least one redirect URI is required'),
	scopes: z.string().optional().default('openid profile email'),
	grant_types: z.string().optional().default('authorization_code'),
	access_token_expiry_minutes: z.coerce.number().int().positive().optional(),
	refresh_token_expiry_hours: z.coerce.number().int().positive().optional(),
	id_token_expiry_minutes: z.coerce.number().int().positive().optional(),
});

export const updateClientSchema = z.object({
	client_name: z.string().trim().optional(),
	redirect_uris: z.string().trim().optional(),
	scopes: z.string().optional(),
	grant_types: z.string().optional(),
	access_token_expiry_minutes: z.coerce.number().int().positive().optional(),
	refresh_token_expiry_hours: z.coerce.number().int().positive().optional(),
	id_token_expiry_minutes: z.coerce.number().int().positive().optional(),
});
