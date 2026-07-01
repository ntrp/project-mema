<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Card } from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';

	interface Props {
		username: string;
		password: string;
		onLogin: (_event: SubmitEvent) => void | Promise<void>;
	}

	let { username = $bindable(), password = $bindable(), onLogin }: Props = $props();
</script>

<Card class="w-full max-w-[420px] p-5" aria-labelledby="login-title">
	<h2 id="login-title" class="m-0 text-3xl font-semibold text-foreground">Admin login</h2>
	<form class="grid gap-4" onsubmit={onLogin}>
		<div class="grid gap-2">
			<Label for="login-username">Username</Label>
			<Input id="login-username" bind:value={username} autocomplete="username" required />
		</div>
		<div class="grid gap-2">
			<Label for="login-password">Password</Label>
			<Input
				id="login-password"
				bind:value={password}
				autocomplete="current-password"
				type="password"
				required
			/>
		</div>
		<Button type="submit">Log in</Button>
	</form>
</Card>
