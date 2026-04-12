<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { DarkMode } from 'flowbite-svelte';
	import {
		ArrowLeft,
		BookOpen,
		ChevronRight,
		Database,
		FileCode2,
		KeyRound,
		Menu,
		Server,
		Shield,
		TreePine,
		Users,
		X,
	} from 'lucide-svelte';

	let { children } = $props();
	let mobileOpen = $state(false);

	const sections = [
		{
			title: 'Getting Started',
			items: [
				{
					label: 'Overview',
					href: '/docs',
					icon: BookOpen,
				},
				{
					label: 'Quick Start',
					href: '/docs/getting-started',
					icon: BookOpen,
				},
				{
					label: 'Docker Setup',
					href: '/docs/docker-setup',
					icon: Server,
				},
			],
		},
		{
			title: 'Admin Panel',
			items: [
				{
					label: 'Register an Application',
					href: '/docs/register-app',
					icon: FileCode2,
				},
				{
					label: 'Manage Users',
					href: '/docs/manage-users',
					icon: Users,
				},
				{
					label: 'Manage Roles',
					href: '/docs/manage-roles',
					icon: Shield,
				},
			],
		},
		{
			title: 'OIDC Integration',
			items: [
				{
					label: 'Overview',
					href: '/docs/oidc-integration',
					icon: KeyRound,
				},
				{
					label: 'JavaScript / TypeScript',
					href: '/docs/oidc-javascript',
					icon: FileCode2,
				},
				{
					label: 'React',
					href: '/docs/oidc-react',
					icon: FileCode2,
				},
				{
					label: 'NextAuth / Auth.js',
					href: '/docs/oidc-nextauth',
					icon: FileCode2,
				},
				{
					label: 'Python (Authlib)',
					href: '/docs/oidc-python',
					icon: FileCode2,
				},
				{
					label: 'Go (go-oidc)',
					href: '/docs/oidc-go',
					icon: FileCode2,
				},
			],
		},
		{
			title: 'Reference',
			items: [
				{
					label: 'API Reference',
					href: '/docs/api-reference',
					icon: Database,
				},
			],
		},
	];

	function isActive(href: string): boolean {
		const current = page.url?.pathname ?? '';
		const resolved = resolve(href);
		return current === resolved;
	}
</script>

<div class="flex min-h-dvh flex-col bg-white dark:bg-gray-900">
	<!-- Header -->
	<header
		class="sticky top-0 z-20 border-b border-gray-200 bg-white dark:border-gray-700 dark:bg-gray-900"
	>
		<div
			class="mx-auto flex h-14 max-w-screen-2xl items-center justify-between px-4 lg:px-8"
		>
			<div class="flex items-center gap-4">
				<button
					type="button"
					class="cursor-pointer rounded-md p-1.5 text-gray-500 hover:bg-gray-100 lg:hidden dark:text-gray-400 dark:hover:bg-gray-800"
					onclick={() => (mobileOpen = !mobileOpen)}
				>
					{#if mobileOpen}
						<X class="size-5" />
					{:else}
						<Menu class="size-5" />
					{/if}
				</button>

				<a
					href={resolve('/')}
					class="flex items-center gap-1.5 text-sm text-gray-500 transition-colors hover:text-primary-600 dark:text-gray-400 dark:hover:text-primary-400"
				>
					<ArrowLeft class="size-4" />
					Home
				</a>
				<span class="text-gray-300 dark:text-gray-600">|</span>
				<a href={resolve('/docs')} class="flex items-center gap-2">
					<TreePine
						class="size-5 text-primary-600 dark:text-primary-400"
						strokeWidth={1.8}
					/>
					<span class="font-semibold dark:text-white">Docs</span>
				</a>
			</div>
			<DarkMode class="rounded-full" />
		</div>
	</header>

	<div class="mx-auto flex w-full max-w-screen-2xl flex-1">
		<!-- Sidebar -->
		<!-- svelte-ignore a11y_no_static_element_interactions -->
		<!-- svelte-ignore a11y_click_events_have_key_events -->
		<div
			class="fixed inset-0 z-10 bg-black/50 transition-opacity lg:hidden {mobileOpen
				? 'opacity-100'
				: 'pointer-events-none opacity-0'}"
			onclick={() => (mobileOpen = false)}
		></div>
		<aside
			class="fixed top-14 bottom-0 z-10 w-72 shrink-0 overflow-y-auto border-r border-gray-200 bg-white transition-transform lg:sticky lg:translate-x-0 lg:bg-transparent dark:border-gray-700 dark:bg-gray-900 {mobileOpen
				? 'translate-x-0'
				: '-translate-x-full'}"
		>
			<nav class="p-4">
				{#each sections as section (section.title)}
					<div class="mb-6">
						<h3
							class="mb-2 px-3 text-xs font-semibold uppercase tracking-wider text-gray-400 dark:text-gray-500"
						>
							{section.title}
						</h3>
						<ul class="space-y-0.5">
							{#each section.items as item (item.href)}
								{@const active = isActive(item.href)}
								<li>
									<a
										href={resolve(item.href)}
										onclick={() => (mobileOpen = false)}
										class="flex items-center gap-2.5 rounded-lg px-3 py-2 text-sm transition-colors {active
											? 'bg-primary-50 font-medium text-primary-700 dark:bg-primary-950 dark:text-primary-400'
											: 'text-gray-600 hover:bg-gray-50 hover:text-gray-900 dark:text-gray-400 dark:hover:bg-gray-800 dark:hover:text-gray-200'}"
									>
										<item.icon class="size-4 shrink-0" />
										{item.label}
										{#if active}
											<ChevronRight class="ml-auto size-3.5 text-primary-400" />
										{/if}
									</a>
								</li>
							{/each}
						</ul>
					</div>
				{/each}
			</nav>
		</aside>

		<!-- Main Content -->
		<main class="min-w-0 flex-1 px-6 py-8 lg:px-12 lg:py-10">
			<article
				class="prose prose-gray dark:prose-invert mx-auto max-w-3xl
					prose-headings:scroll-mt-20
					prose-h1:text-3xl prose-h1:font-extrabold prose-h1:tracking-tight
					prose-h2:text-2xl prose-h2:font-bold prose-h2:border-b prose-h2:border-gray-200 prose-h2:pb-2 dark:prose-h2:border-gray-700
					prose-h3:text-lg prose-h3:font-semibold
					prose-code:before:content-none prose-code:after:content-none
					prose-code:rounded prose-code:px-1.5 prose-code:py-0.5 prose-code:text-sm prose-code:font-normal
					prose-pre:rounded-xl prose-pre:shadow-sm prose-pre:ring-1 prose-pre:ring-gray-200 dark:prose-pre:shadow-lg dark:prose-pre:ring-gray-700
					prose-a:text-primary-600 prose-a:no-underline hover:prose-a:underline dark:prose-a:text-primary-400
					prose-table:block prose-table:overflow-x-auto prose-table:rounded-lg prose-table:border prose-table:border-gray-200 dark:prose-table:border-gray-700
					prose-thead:bg-gray-50 dark:prose-thead:bg-gray-800
					prose-th:px-4 prose-th:py-3 prose-th:text-xs prose-th:font-semibold prose-th:uppercase prose-th:tracking-wider
					prose-td:px-4 prose-td:py-3
					prose-tr:border-b prose-tr:border-gray-200 dark:prose-tr:border-gray-700
					prose-img:rounded-xl prose-img:shadow-md
					prose-blockquote:border-primary-500 prose-blockquote:bg-primary-50/50 prose-blockquote:rounded-r-lg prose-blockquote:py-1 dark:prose-blockquote:bg-primary-950/30"
			>
				{@render children()}
			</article>

			<!-- Footer -->
			<div
				class="mx-auto mt-16 max-w-3xl border-t border-gray-200 pt-8 dark:border-gray-700"
			>
				<p class="text-sm text-gray-400 dark:text-gray-500">
					AuthForest &mdash; Open-source OIDC Identity Provider
				</p>
			</div>
		</main>
	</div>
</div>
