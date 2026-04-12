<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import MultiInput from '$lib/components/common/multi-input/MultiInput.svelte';
	import { OIDC } from '$lib/resources/oidc';
	import { updateClientSchema } from '$lib/schemas/oidc.schema';
	import { toastService } from '$lib/services/toast.service.svelte';
	import { Button, Helper, Input, Label, Spinner } from 'flowbite-svelte';
	import { onMount } from 'svelte';
	import { defaults, superForm } from 'sveltekit-superforms';
	import { zod4, zod4Client } from 'sveltekit-superforms/adapters';

	const clientId = page.params.clientId!;
	let loading = $state(true);

	const superform = superForm(defaults(zod4(updateClientSchema)), {
		validators: zod4Client(updateClientSchema),
		SPA: true,
		onUpdate: async ({ form: f }) => {
			if (!f.valid || !$tainted) return;

			const payload: Record<string, unknown> = {};
			for (const [key, value] of Object.entries(f.data)) {
				if (value !== undefined && value !== '' && value !== null) {
					payload[key] = value;
				}
			}

			const [err] = await OIDC.updateClient(clientId, payload);
			if (err) {
				toastService.error('Failed to update application');
				return;
			}

			toastService.success('Application updated successfully');
			goto(resolve('/applications'));
		},
	});

	const { form, enhance, formId, tainted, submitting } = superform;

	onMount(async () => {
		const [err, client] = await OIDC.getClient(clientId);
		if (err) {
			toastService.error('Failed to load application');
			loading = false;
			return;
		}
		$form.client_name = client.client_name ?? '';
		$form.redirect_uris = client.redirect_uris?.length
			? [...client.redirect_uris]
			: [''];
		$form.scopes = client.scopes ?? '';
		$form.grant_types = client.grant_types ?? '';
		$form.access_token_expiry_minutes =
			client.access_token_expiry_minutes ?? undefined;
		$form.refresh_token_expiry_hours =
			client.refresh_token_expiry_hours ?? undefined;
		$form.id_token_expiry_minutes = client.id_token_expiry_minutes ?? undefined;
		loading = false;
	});
</script>

{#if loading}
	<div class="flex justify-center py-12">
		<Spinner size="10" />
	</div>
{:else}
	<div class="mx-auto max-w-2xl space-y-6">
		<div>
			<h1 class="text-3xl font-bold dark:text-white">Edit Application</h1>
			<p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
				Update the configuration for this OAuth 2.0 / OIDC client.
			</p>
		</div>

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
				/>
			</div>

			<div class="space-y-2">
				<Label for="redirect_uris">Redirect URIs</Label>
				<MultiInput
					bind:values={$form.redirect_uris as string[]}
					id="redirect_uris"
					placeholder="http://localhost:4200/callback"
				/>
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

			<div class="flex flex-col gap-3 pt-2 sm:flex-row sm:items-center">
				<Button
					type="submit"
					form={$formId}
					class="cursor-pointer"
					loading={$submitting}
				>
					Save Changes
				</Button>
				<a
					href={resolve('/applications')}
					class="text-sm text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200"
				>
					Cancel
				</a>
			</div>
		</form>
	</div>
{/if}
