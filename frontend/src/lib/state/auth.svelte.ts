import type { User as OidcUser } from 'oidc-client-ts';

export interface UserProfile {
	firstName: string;
	lastName: string;
	fullName: string;
	username: string;
	email: string;
}

let oidcUser = $state<OidcUser | null>(null);

const profile = $derived<UserProfile | null>(
	oidcUser?.profile
		? {
				firstName: (oidcUser.profile.given_name as string) ?? '',
				lastName: (oidcUser.profile.family_name as string) ?? '',
				fullName: (oidcUser.profile.name as string) ?? '',
				username: (oidcUser.profile.preferred_username as string) ?? '',
				email: (oidcUser.profile.email as string) ?? '',
			}
		: null
);

const userInitials = $derived(
	`${profile?.firstName?.[0] ?? ''}${profile?.lastName?.[0] ?? ''}`
);

function setUser(user: OidcUser | null) {
	oidcUser = user;
}

export const authStore = {
	get user() {
		return profile;
	},
	get accessToken() {
		return oidcUser?.access_token ?? null;
	},
	get userInitials() {
		return userInitials;
	},
	setUser,
};
