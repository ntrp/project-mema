<script lang="ts">
	import EyeOffIcon from '@lucide/svelte/icons/eye-off';
	import StatusPill from '$lib/components/shared/StatusPill.svelte';
	import { Button } from '$lib/components/ui/button';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { resolve } from '$app/paths';
	import type { MediaSearchResult } from '$lib/settings/types';
	import MediaAddButton from './MediaAddButton.svelte';
	import MediaBadge from './MediaBadge.svelte';
	import PosterPlaceholder from './PosterPlaceholder.svelte';

	interface Props {
		result: MediaSearchResult;
		adding?: boolean;
		actionLabel: string;
		inLibrary?: boolean;
		onAdd: (_candidate: MediaSearchResult) => void;
		onBlacklist?: (_candidate: MediaSearchResult) => void;
		blacklisting?: boolean;
		showBlacklistAction?: boolean;
	}

	let {
		result,
		adding = false,
		actionLabel,
		inLibrary = false,
		onAdd,
		onBlacklist,
		blacklisting = false,
		showBlacklistAction = false
	}: Props = $props();

	function posterUrl(path?: string) {
		if (!path) {
			return undefined;
		}
		if (path.startsWith('http://') || path.startsWith('https://')) {
			return path;
		}
		return `https://image.tmdb.org/t/p/w342${path}`;
	}
</script>

<article class="group/poster min-w-0 snap-start">
	<div
		class="relative aspect-[2/3] overflow-hidden rounded-md border border-border bg-card transition-[transform,border-color,box-shadow] duration-150 group-hover/poster:z-[2] group-hover/poster:-translate-y-1.5 group-hover/poster:scale-105 group-hover/poster:border-primary/50 group-hover/poster:shadow-xl group-focus-within/poster:z-[2] group-focus-within/poster:-translate-y-1.5 group-focus-within/poster:scale-105 group-focus-within/poster:border-primary/50 group-focus-within/poster:shadow-xl"
	>
		{#if posterUrl(result.posterPath)}
			<img
				class="block h-full w-full object-cover"
				src={posterUrl(result.posterPath)}
				alt=""
				loading="lazy"
			/>
		{:else}
			<PosterPlaceholder label={result.type} />
		{/if}
		{#if result.externalProvider && result.externalId}
			<a
				class="absolute inset-0 z-[1] rounded-md"
				href={resolve('/media/[provider]/[type]/[externalId]', {
					provider: result.externalProvider,
					type: result.type,
					externalId: result.externalId
				})}
				aria-label={`Open ${result.title} details`}
			></a>
		{/if}
		<MediaBadge type={result.type} {inLibrary} />
		<div
			class="pointer-events-none absolute inset-0 z-[2] flex flex-col justify-end gap-1.5 bg-card/70 px-[13px] pt-[58px] pb-[13px] opacity-0 transition-opacity duration-150 group-hover/poster:opacity-100 group-focus-within/poster:opacity-100"
		>
			<span class="text-sm leading-none text-primary-foreground">{result.year ?? 'Unknown'}</span>
			<h3 class="m-0 text-xl leading-tight text-primary-foreground">{result.title}</h3>
			<p class="line-clamp-4 m-0 text-[13px] leading-tight text-primary-foreground">
				{result.overview ?? 'No overview available.'}
			</p>
			{#if showBlacklistAction && onBlacklist}
				<Tooltip.Root>
					<Tooltip.Trigger>
						{#snippet child({ props })}
							<Button
								{...props}
								type="button"
								variant="outline"
								size="icon-sm"
								class="pointer-events-auto absolute top-2.5 right-2.5 size-[34px] min-h-[34px] min-w-[34px] border-border bg-card/80 p-0 text-foreground backdrop-blur-md hover:border-primary/50 hover:bg-muted hover:text-primary-foreground focus-visible:border-primary/50 focus-visible:bg-muted focus-visible:text-primary-foreground"
								disabled={blacklisting}
								aria-label={`Hide ${result.title} from discover`}
								onclick={(event) => {
									event.stopPropagation();
									onBlacklist(result);
								}}
							>
								<EyeOffIcon aria-hidden="true" />
							</Button>
						{/snippet}
					</Tooltip.Trigger>
					<Tooltip.Content>Hide from discover</Tooltip.Content>
				</Tooltip.Root>
			{/if}
			{#if inLibrary}
				<StatusPill tone="success">In library</StatusPill>
			{:else}
				<MediaAddButton
					{result}
					{adding}
					label={actionLabel}
					class="pointer-events-auto mt-0.5 min-h-[30px] self-start px-3 text-[13px]"
					{onAdd}
				/>
			{/if}
		</div>
	</div>
</article>
