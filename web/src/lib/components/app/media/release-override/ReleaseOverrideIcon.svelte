<script lang="ts">
	import ArrowDownIcon from '@lucide/svelte/icons/arrow-down';
	import CircleIcon from '@lucide/svelte/icons/circle';
	import UserIcon from '@lucide/svelte/icons/user';
</script>

<span class="relative inline-flex size-4" aria-hidden="true">
	<UserIcon class="absolute left-0 top-0 size-4" />
	<span class="absolute -bottom-0.5 -right-1 inline-flex size-3.5 items-center justify-center">
		<CircleIcon class="absolute inset-0 size-3.5 fill-current stroke-current" />
		<ArrowDownIcon class="relative size-2 stroke-[3] text-amber-400" />
	</span>
</span>
