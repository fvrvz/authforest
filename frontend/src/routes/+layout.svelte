<script lang="ts">
	import ConfirmDialog from '$lib/components/common/dialog/ConfirmDialog.svelte';
	import ToastProvider from '$lib/components/common/toast/ToastProvider.svelte';
	import '../app.css';

	let { children } = $props();
</script>

<svelte:head>
	<title>AuthForest</title>
</svelte:head>

{@render children()}

<ToastProvider />
<ConfirmDialog />
