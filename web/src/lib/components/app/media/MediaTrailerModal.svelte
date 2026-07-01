<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog';

	interface Props {
		title: string;
		url: string;
		onClose: () => void;
	}

	let { title, url, onClose }: Props = $props();
	let open = $state(true);
	const embedUrl = $derived(toEmbedUrl(url));

	function handleOpenChange(nextOpen: boolean) {
		open = nextOpen;
		if (!nextOpen) {
			onClose();
		}
	}

	function toEmbedUrl(value: string) {
		try {
			const parsed = new globalThis.URL(value);
			const host = parsed.hostname.replace(/^www\./, '');
			if (host === 'youtube.com' || host === 'm.youtube.com') {
				const key = parsed.searchParams.get('v');
				return key ? `https://www.youtube.com/embed/${key}` : value;
			}
			if (host === 'youtu.be') {
				const key = parsed.pathname.slice(1);
				return key ? `https://www.youtube.com/embed/${key}` : value;
			}
		} catch {
			return value;
		}
		return value;
	}
</script>

<Dialog.Root bind:open onOpenChange={handleOpenChange}>
	<Dialog.Content class="w-[min(1120px,calc(100vw-32px))] gap-4 p-4 sm:max-w-none">
		<Dialog.Header class="pr-10">
			<Dialog.Title>{title}</Dialog.Title>
		</Dialog.Header>
		<div class="aspect-video overflow-hidden rounded-md bg-black">
			<iframe
				class="block size-full"
				src={embedUrl}
				{title}
				allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share"
				allowfullscreen
			></iframe>
		</div>
	</Dialog.Content>
</Dialog.Root>
