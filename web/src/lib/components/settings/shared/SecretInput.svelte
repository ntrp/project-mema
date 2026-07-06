<script lang="ts">
	import EyeIcon from '@lucide/svelte/icons/eye';
	import EyeOffIcon from '@lucide/svelte/icons/eye-off';
	import type { HTMLInputAttributes } from 'svelte/elements';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { cn } from '$lib/utils';

	type Props = Omit<HTMLInputAttributes, 'type' | 'value' | 'files'> & {
		value?: string;
		onValueChange?: (_value: string) => void;
	};

	let {
		value = $bindable(''),
		onValueChange,
		class: className,
		disabled,
		...restProps
	}: Props = $props();

	let visible = $state(false);
	const label = $derived(visible ? 'Hide secret' : 'Show secret');
	const Icon = $derived(visible ? EyeOffIcon : EyeIcon);

	function updateValue(event: Event & { currentTarget: HTMLInputElement }) {
		onValueChange?.(event.currentTarget.value);
	}
</script>

<div class="relative">
	<Input
		{...restProps}
		{disabled}
		class={cn('pr-10', className)}
		type={visible ? 'text' : 'password'}
		bind:value
		oninput={updateValue}
	/>
	<div class="absolute inset-y-0 right-1 flex items-center">
		<Tooltip.Provider>
			<Tooltip.Root>
				<Tooltip.Trigger>
					{#snippet child({ props })}
						<Button
							{...props}
							type="button"
							variant="ghost"
							size="icon-xs"
							aria-label={label}
							{disabled}
							onclick={() => (visible = !visible)}
						>
							<Icon aria-hidden="true" />
						</Button>
					{/snippet}
				</Tooltip.Trigger>
				<Tooltip.Content>{label}</Tooltip.Content>
			</Tooltip.Root>
		</Tooltip.Provider>
	</div>
</div>
