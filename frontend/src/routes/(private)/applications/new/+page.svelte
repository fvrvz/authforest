<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { OIDC } from '$lib/resources/oidc';
	import { registerClientSchema } from '$lib/schemas/oidc.schema';
	import { toastService } from '$lib/services/toast.service.svelte';
	import type { OAuthClient } from '$lib/types/oidc.type';
	import { Button, Helper, Input, Label, Select } from 'flowbite-svelte';
	import { defaults, superForm } from 'sveltekit-superforms';
	import { zod4, zod4Client } from 'sveltekit-superforms/adapters';

	let createdClient = $state<OAuthClient | null>(null);

	const superform = superForm(
		defaults(
			{
				client_name: '',
				client_type: 'public' as const,
				redirect_uris: '',
				scopes: 'openid profile email',
				grant_types: 'authorization_code',
			},
			zod4(registerClientSchema)
		),
		{
			validators: zod4Client(registerClientSchema),
			SPA: true,
			onUpdate: async ({ form }) => {
				if (!form.valid || !$tainted) return;

				const [err, data] = await OIDC.registerClient(form.data);
				if (err) {
					toastService.error('Failed to register application');
					return;
				}

				createdClient = data;
				toastService.success('Application registered successfully');
			},
		}
	);

	const { form, tainted, enhance, formId, submitting, errors } = superform;

	const clientTypes = [
		{ value: 'public', name: 'Public (SPA, Mobile App)' },
		{ value: 'confidential', name: 'Confidential (Server-side App)' },
	];
</script>

<div class="mx-auto max-w-2xl space-y-6">
	<div>
		<h1 class="text-3xl font-bold dark:text-white">Register Application</h1>
		<p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
			Create a new OAuth 2.0 / OIDC client application.
		</p>
	</div>

	{#if createdClient}
		<div
			class="rounded-xl border border-green-200 bg-green-50 p-6 dark:border-green-700 dark:bg-green-900/20"
		>
			<h2 class="mb-4 text-xl font-semibold text-green-800 dark:text-green-300">
				Application Registered Successfully
			</h2>
			<div class="space-y-3 text-sm">
				<div>
					<span class="font-medium text-gray-700 dark:text-gray-300"
						>Client ID:</span
					>
					<code
						class="ml-2 rounded bg-gray-100 px-2 py-1 text-xs dark:bg-gray-800"
					>
						{createdClient.client_id}
					</code>
				</div>
				{#if createdClient.client_secret}
					<div>
						<span class="font-medium text-gray-700 dark:text-gray-300"
							>Client Secret:</span
						>
						<code
							class="ml-2 rounded bg-gray-100 px-2 py-1 text-xs dark:bg-gray-800"
						>
							{createdClient.client_secret}
						</code>
					</div>
					<p class="text-xs font-medium text-red-600 dark:text-red-400">
						Save the client secret now. You won't be able to see it again.
					</p>
				{/if}
				<div>
					<span class="font-medium text-gray-700 dark:text-gray-300"
						>Application Name:</span
					>
					<span class="ml-2">{createdClient.client_name}</span>
				</div>
				<div>
					<span class="font-medium text-gray-700 dark:text-gray-300">Type:</span
					>
					<span class="ml-2">{createdClient.client_type}</span>
				</div>
			</div>
			<div class="mt-6 flex gap-3">
				<Button
					color="alternative"
					class="cursor-pointer"
					onclick={() => {
						createdClient = null;
						$form.client_name = '';
						$form.redirect_uris = '';
					}}
				>
					Register Another
				</Button>
				<Button
					class="cursor-pointer"
					onclick={() => goto(resolve('/applications'))}
				>
					View All Applications
				</Button>
			</div>
		</div>
	{:else}
		<form
			use:enhance
			class="space-y-5 rounded-xl border border-gray-200 bg-white p-6 shadow-sm dark:border-gray-600 dark:bg-gray-800"
			id={$formId}
		>
			<div class="space-y-2">
				<Label for="client_name">Application Name</Label>
				<Input
					type="text"
					bind:value={$form.client_name}
					id="client_name"
					placeholder="e.g. My Angular SPA"
					required
				/>
				{#if $errors.client_name}
					<Helper color="red">{$errors.client_name}</Helper>
				{/if}
			</div>

			<div class="space-y-2">
				<Label for="client_type">Client Type</Label>
				<Select
					id="client_type"
					items={clientTypes}
					bind:value={$form.client_type}
				/>
				<Helper>
					{$form.client_type === 'public'
						? 'For browser-based apps (SPAs) and mobile apps. Uses PKCE for security.'
						: 'For server-side apps that can securely store a client secret.'}
				</Helper>
			</div>

			<div class="space-y-2">
				<Label for="redirect_uris">Redirect URIs</Label>
				<Input
					type="text"
					bind:value={$form.redirect_uris}
					id="redirect_uris"
					placeholder="http://localhost:4200/callback"
					required
				/>
				<Helper>Comma-separated list of allowed callback URLs.</Helper>
				{#if $errors.redirect_uris}
					<Helper color="red">{$errors.redirect_uris}</Helper>
				{/if}
			</div>

			<div class="space-y-2">
				<Label for="scopes">Scopes</Label>
				<Input
					type="text"
					bind:value={$form.scopes}
					id="scopes"
					placeholder="openid profile email offline_access"
				/>
				<Helper
					>Space-separated scopes. Available: openid, profile, email,
					offline_access</Helper
				>
			</div>

			<div class="space-y-2">
				<Label for="grant_types">Grant Types</Label>
				<Input
					type="text"
					bind:value={$form.grant_types}
					id="grant_types"
					placeholder="authorization_code"
				/>
				<Helper
					>Space-separated. Available: authorization_code, refresh_token</Helper
				>
			</div>

			<h3 class="pt-2 text-lg font-semibold dark:text-white">
				Token Expiry Configuration
			</h3>

			<div class="grid gap-4 sm:grid-cols-3">
				<div class="space-y-2">
					<Label for="access_token_expiry_minutes">Access Token (minutes)</Label
					>
					<Input
						type="number"
						bind:value={$form.access_token_expiry_minutes}
						id="access_token_expiry_minutes"
						placeholder="15"
						min="1"
					/>
				</div>

				<div class="space-y-2">
					<Label for="refresh_token_expiry_hours">Refresh Token (hours)</Label>
					<Input
						type="number"
						bind:value={$form.refresh_token_expiry_hours}
						id="refresh_token_expiry_hours"
						placeholder="2"
						min="1"
					/>
				</div>

				<div class="space-y-2">
					<Label for="id_token_expiry_minutes">ID Token (minutes)</Label>
					<Input
						type="number"
						bind:value={$form.id_token_expiry_minutes}
						id="id_token_expiry_minutes"
						placeholder="15"
						min="1"
					/>
				</div>
			</div>
			<Helper
				>Leave blank to use server defaults (Access: 15 min, Refresh: 2 hr, ID:
				15 min)</Helper
			>

			<div class="flex items-center gap-3 pt-2">
				<Button
					type="submit"
					form={$formId}
					class="cursor-pointer"
					loading={$submitting}
				>
					Register Application
				</Button>
				<a
					href={resolve('/applications')}
					class="text-sm text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200"
				>
					Cancel
				</a>
			</div>
		</form>
	{/if}
</div>
