<script lang="ts">
	import { Badge } from 'flowbite-svelte';
</script>

<div class="space-y-6">
	<div>
		<h1 class="text-3xl font-bold dark:text-white">Settings</h1>
		<p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
			Identity Provider configuration and endpoint reference.
		</p>
	</div>

	<div
		class="rounded-xl border border-gray-200 bg-white p-6 dark:border-gray-600 dark:bg-gray-800"
	>
		<h2 class="mb-4 text-lg font-semibold dark:text-white">OIDC Endpoints</h2>
		<div class="overflow-x-auto">
			<table class="w-full text-left text-sm">
				<thead
					class="border-b border-gray-200 text-xs uppercase text-gray-500 dark:border-gray-700 dark:text-gray-400"
				>
					<tr>
						<th class="pb-3 pr-4">Endpoint</th>
						<th class="pb-3 pr-4">Method</th>
						<th class="pb-3">Path</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-gray-100 dark:divide-gray-700">
					<tr>
						<td class="py-3 pr-4 font-medium dark:text-white">Discovery</td>
						<td class="py-3 pr-4"><Badge color="green">GET</Badge></td>
						<td class="py-3"
							><code class="text-xs">/.well-known/openid-configuration</code
							></td
						>
					</tr>
					<tr>
						<td class="py-3 pr-4 font-medium dark:text-white">JWKS</td>
						<td class="py-3 pr-4"><Badge color="green">GET</Badge></td>
						<td class="py-3"
							><code class="text-xs">/.well-known/jwks.json</code></td
						>
					</tr>
					<tr>
						<td class="py-3 pr-4 font-medium dark:text-white">Authorize</td>
						<td class="py-3 pr-4"><Badge color="green">GET</Badge></td>
						<td class="py-3"><code class="text-xs">/oauth2/authorize</code></td>
					</tr>
					<tr>
						<td class="py-3 pr-4 font-medium dark:text-white">Token</td>
						<td class="py-3 pr-4"><Badge color="blue">POST</Badge></td>
						<td class="py-3"><code class="text-xs">/oauth2/token</code></td>
					</tr>
					<tr>
						<td class="py-3 pr-4 font-medium dark:text-white">UserInfo</td>
						<td class="py-3 pr-4"><Badge color="green">GET</Badge></td>
						<td class="py-3"><code class="text-xs">/oauth2/userinfo</code></td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>

	<div
		class="rounded-xl border border-gray-200 bg-white p-6 dark:border-gray-600 dark:bg-gray-800"
	>
		<h2 class="mb-4 text-lg font-semibold dark:text-white">Supported Scopes</h2>
		<div class="space-y-4">
			<div class="flex items-start gap-3">
				<Badge color="gray" class="mt-0.5 shrink-0">openid</Badge>
				<div>
					<p class="text-sm text-gray-700 dark:text-gray-300">
						Required. Returns sub, iss, aud, exp, iat, auth_time, nonce claims
						in the ID token.
					</p>
				</div>
			</div>
			<div class="flex items-start gap-3">
				<Badge color="gray" class="mt-0.5 shrink-0">profile</Badge>
				<div>
					<p class="text-sm text-gray-700 dark:text-gray-300">
						Returns name, given_name, family_name, preferred_username,
						updated_at.
					</p>
				</div>
			</div>
			<div class="flex items-start gap-3">
				<Badge color="gray" class="mt-0.5 shrink-0">email</Badge>
				<div>
					<p class="text-sm text-gray-700 dark:text-gray-300">
						Returns email and email_verified claims.
					</p>
				</div>
			</div>
			<div class="flex items-start gap-3">
				<Badge color="gray" class="mt-0.5 shrink-0">offline_access</Badge>
				<div>
					<p class="text-sm text-gray-700 dark:text-gray-300">
						Issues a refresh token for long-lived sessions.
					</p>
				</div>
			</div>
		</div>
	</div>

	<div
		class="rounded-xl border border-gray-200 bg-white p-6 dark:border-gray-600 dark:bg-gray-800"
	>
		<h2 class="mb-4 text-lg font-semibold dark:text-white">
			Security Configuration
		</h2>
		<div class="grid gap-4 sm:grid-cols-2">
			<div class="rounded-lg border border-gray-100 p-4 dark:border-gray-700">
				<p class="text-sm font-medium text-gray-500 dark:text-gray-400">
					Token Signing Algorithm
				</p>
				<p class="mt-1 text-lg font-semibold dark:text-white">RS256</p>
			</div>
			<div class="rounded-lg border border-gray-100 p-4 dark:border-gray-700">
				<p class="text-sm font-medium text-gray-500 dark:text-gray-400">
					PKCE Method
				</p>
				<p class="mt-1 text-lg font-semibold dark:text-white">S256</p>
			</div>
			<div class="rounded-lg border border-gray-100 p-4 dark:border-gray-700">
				<p class="text-sm font-medium text-gray-500 dark:text-gray-400">
					Grant Type
				</p>
				<p class="mt-1 text-lg font-semibold dark:text-white">
					Authorization Code
				</p>
			</div>
			<div class="rounded-lg border border-gray-100 p-4 dark:border-gray-700">
				<p class="text-sm font-medium text-gray-500 dark:text-gray-400">
					Refresh Token
				</p>
				<p class="mt-1 text-lg font-semibold dark:text-white">
					Rotation Enabled
				</p>
			</div>
		</div>
	</div>
</div>
