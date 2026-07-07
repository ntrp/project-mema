<script lang="ts">
	import type { DLNAStatus } from '$lib/settings/types';

	interface Props {
		status?: DLNAStatus;
	}

	let { status }: Props = $props();
</script>

<div class="grid gap-4">
	<div class="grid gap-2">
		<h3 class="m-0 text-sm font-semibold">Advertised URLs</h3>
		{#if status?.advertisedUrls?.length}
			<ul class="m-0 grid gap-1 pl-4 text-sm text-muted-foreground">
				{#each status.advertisedUrls as url (url)}
					<li class="break-all">{url}</li>
				{/each}
			</ul>
		{:else}
			<p class="m-0 text-sm text-muted-foreground">No advertised URLs</p>
		{/if}
	</div>

	<div class="grid gap-2">
		<h3 class="m-0 text-sm font-semibold">Recent clients</h3>
		<div class="overflow-x-auto rounded-md border border-border">
			<table class="w-full min-w-[720px] border-collapse text-left text-sm">
				<thead class="bg-muted/50 text-xs text-muted-foreground uppercase">
					<tr>
						<th class="px-3 py-2 font-medium">Client IP</th>
						<th class="px-3 py-2 font-medium">Profile</th>
						<th class="px-3 py-2 font-medium">Last SOAP</th>
						<th class="px-3 py-2 font-medium">Last seen</th>
						<th class="px-3 py-2 font-medium">Last error</th>
					</tr>
				</thead>
				<tbody>
					{#each status?.recentClients ?? [] as client (client.ip)}
						<tr class="border-t border-border">
							<td class="px-3 py-2 font-medium">{client.ip}</td>
							<td class="px-3 py-2">{client.profileId}</td>
							<td class="px-3 py-2">{client.lastSoapAction || 'None'}</td>
							<td class="px-3 py-2">{new Date(client.lastSeen).toLocaleString()}</td>
							<td class="px-3 py-2">{client.lastError ?? 'None'}</td>
						</tr>
					{:else}
						<tr>
							<td class="px-3 py-4 text-muted-foreground" colspan="5">No recent clients</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	</div>

	<div class="grid gap-2">
		<h3 class="m-0 text-sm font-semibold">Active streams</h3>
		<p class="m-0 text-sm text-muted-foreground">
			{status?.activeStreams?.length ?? 0} streams, {status?.activeTranscodes?.length ?? 0} transcodes
		</p>
	</div>
</div>
