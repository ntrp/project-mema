<script lang="ts">
	import { onMount } from 'svelte';
	import * as Card from '$lib/components/ui/card';
	import ManualFulfillmentActionsTable from './ManualFulfillmentActionsTable.svelte';
	import {
		listManualFulfillmentActions,
		type ManualFulfillmentAction
	} from './manualFulfillmentActions';

	let actions = $state<ManualFulfillmentAction[]>([]);
	let loading = $state(true);
	let errorMessage = $state('');

	onMount(() => {
		void load();
	});

	async function load() {
		loading = true;
		errorMessage = '';
		try {
			actions = await listManualFulfillmentActions();
		} catch (error) {
			errorMessage = error instanceof Error ? error.message : 'Could not load manual actions';
		} finally {
			loading = false;
		}
	}
</script>

<Card.Root>
	<Card.Header>
		<Card.Title>Manual Fulfillment Actions</Card.Title>
	</Card.Header>
	<Card.Content>
		<ManualFulfillmentActionsTable {actions} {loading} {errorMessage} />
	</Card.Content>
</Card.Root>
