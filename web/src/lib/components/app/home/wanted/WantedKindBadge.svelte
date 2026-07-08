<script lang="ts">
	import { Badge } from '$lib/components/ui/badge';
	import type { WantedRowKind } from './wantedRows';

	interface Props {
		kind: WantedRowKind;
	}

	let { kind }: Props = $props();

	const label = $derived(
		kind === 'custom_format_upgrade' ? 'Upgrade' : kind === 'target' ? 'Target' : 'Media'
	);
	const variant = $derived(kind === 'media' ? 'destructive' : 'secondary');
</script>

<Badge {variant}>{label}</Badge>
