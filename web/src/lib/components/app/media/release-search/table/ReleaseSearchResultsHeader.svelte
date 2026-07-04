<script lang="ts">
	import ArrowDownIcon from '@lucide/svelte/icons/arrow-down';
	import ArrowLeftRightIcon from '@lucide/svelte/icons/arrow-left-right';
	import ArrowUpIcon from '@lucide/svelte/icons/arrow-up';
	import BadgeCheckIcon from '@lucide/svelte/icons/badge-check';
	import ClockIcon from '@lucide/svelte/icons/clock';
	import HardDriveIcon from '@lucide/svelte/icons/hard-drive';
	import ListChecksIcon from '@lucide/svelte/icons/list-checks';
	import ServerIcon from '@lucide/svelte/icons/server';
	import StarIcon from '@lucide/svelte/icons/star';
	import UsersIcon from '@lucide/svelte/icons/users';
	import { Button } from '$lib/components/ui/button';
	import * as Table from '$lib/components/ui/table';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import type {
		ReleaseSort,
		ReleaseSortKey
	} from '$lib/components/app/media/release-display/releaseSearchResults';

	interface Props {
		sort: ReleaseSort;
		onSort: (_key: ReleaseSortKey) => void;
	}

	let { sort, onSort }: Props = $props();

	const columns: { key: ReleaseSortKey; label: string; icon?: string }[] = [
		{ key: 'source', label: 'Protocol', icon: 'transfer' },
		{ key: 'indexer', label: 'Indexer', icon: 'server' },
		{ key: 'age', label: 'Age', icon: 'time' },
		{ key: 'title', label: 'Title' },
		{ key: 'size', label: 'Size', icon: 'size' },
		{ key: 'quality', label: 'Quality', icon: 'quality' },
		{ key: 'score', label: 'Score', icon: 'score' }
	];
	const peersColumn: { key: ReleaseSortKey; label: string; icon: string } = {
		key: 'peers',
		label: 'Peers',
		icon: 'peers'
	};
</script>

<Table.Header class="sticky top-0 z-10 bg-card">
	<Table.Row>
		{#each columns as column (column.key)}
			<Table.Head class={column.key === 'title' ? 'w-full min-w-0 max-w-0' : 'w-[1%]'}>
				<div class="flex items-center">
					<Tooltip.Root>
						<Tooltip.Trigger>
							{#snippet child({ props })}
								<Button
									{...props}
									type="button"
									variant="ghost"
									class="-ml-2 h-8 px-2"
									aria-label={`Sort by ${column.label}`}
									onclick={() => onSort(column.key)}
								>
									{#if column.icon === 'transfer'}
										<ArrowLeftRightIcon aria-hidden="true" />
										<span class="sr-only">{column.label}</span>
									{:else if column.icon === 'time'}
										<ClockIcon aria-hidden="true" />
										<span class="sr-only">{column.label}</span>
									{:else if column.icon === 'size'}
										<HardDriveIcon aria-hidden="true" />
										<span class="sr-only">{column.label}</span>
									{:else if column.icon === 'server'}
										<ServerIcon aria-hidden="true" />
										<span class="sr-only">{column.label}</span>
									{:else if column.icon === 'quality'}
										<BadgeCheckIcon aria-hidden="true" />
										<span class="sr-only">{column.label}</span>
									{:else if column.icon === 'score'}
										<StarIcon aria-hidden="true" />
										<span class="sr-only">{column.label}</span>
									{:else}
										<span>{column.label}</span>
									{/if}
									{#if sort.key === column.key}
										{#if sort.direction === 'asc'}
											<ArrowUpIcon aria-hidden="true" />
										{:else}
											<ArrowDownIcon aria-hidden="true" />
										{/if}
									{/if}
								</Button>
							{/snippet}
						</Tooltip.Trigger>
						<Tooltip.Content>{column.label}</Tooltip.Content>
					</Tooltip.Root>
					{#if column.key === 'source'}
						<Tooltip.Root>
							<Tooltip.Trigger>
								{#snippet child({ props })}
									<Button
										{...props}
										type="button"
										variant="ghost"
										class="h-8 px-2"
										aria-label={`Sort by ${peersColumn.label}`}
										onclick={() => onSort(peersColumn.key)}
									>
										<UsersIcon aria-hidden="true" />
										<span class="sr-only">{peersColumn.label}</span>
										{#if sort.key === peersColumn.key}
											{#if sort.direction === 'asc'}
												<ArrowUpIcon aria-hidden="true" />
											{:else}
												<ArrowDownIcon aria-hidden="true" />
											{/if}
										{/if}
									</Button>
								{/snippet}
							</Tooltip.Trigger>
							<Tooltip.Content>{peersColumn.label}</Tooltip.Content>
						</Tooltip.Root>
					{/if}
				</div>
			</Table.Head>
		{/each}
		<Table.Head class="w-[1%]">
			<Tooltip.Root>
				<Tooltip.Trigger>
					{#snippet child({ props })}
						<span {...props} class="-ml-2 inline-flex h-8 items-center px-2 text-muted-foreground">
							<ListChecksIcon aria-hidden="true" />
							<span class="sr-only">Match</span>
						</span>
					{/snippet}
				</Tooltip.Trigger>
				<Tooltip.Content>Match</Tooltip.Content>
			</Tooltip.Root>
		</Table.Head>
		<Table.Head />
	</Table.Row>
</Table.Header>
