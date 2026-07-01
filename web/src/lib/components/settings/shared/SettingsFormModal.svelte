<script lang="ts">
	import type { Snippet } from 'svelte';

	interface Props {
		title: string;
		onClose: () => void;
		children: Snippet;
		modalClass?: string;
	}

	let { title, onClose, children, modalClass = '' }: Props = $props();
</script>

<div class="modal-backdrop" role="presentation" onclick={onClose}>
	<div
		class={`modal-shell settings-form-modal ${modalClass}`.trim()}
		role="dialog"
		aria-modal="true"
		aria-labelledby="settings-form-modal-title"
		tabindex="-1"
		onclick={(event) => event.stopPropagation()}
		onkeydown={(event) => event.stopPropagation()}
	>
		<div class="section-heading">
			<h2 id="settings-form-modal-title">{title}</h2>
			<button type="button" class="secondary" onclick={onClose}>Close</button>
		</div>
		{@render children()}
	</div>
</div>
