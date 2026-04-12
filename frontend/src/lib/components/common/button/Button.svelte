<script lang="ts">
	import type { ButtonProps } from '$lib/types/button.type';

	const { children, class: classProp, ...rest }: ButtonProps = $props();
</script>

<button
	{...rest}
	class={`rounded bg-blue-300 px-4 py-2  hover:bg-blue-400 ${classProp}`}
>
	{@render children()}
</button>
