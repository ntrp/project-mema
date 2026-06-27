<script lang="ts">
	interface Props {
		username: string;
		password: string;
		onLogin: (_event: SubmitEvent) => void | Promise<void>;
	}

	let { username = $bindable(), password = $bindable(), onLogin }: Props = $props();
</script>

<section class="panel auth-panel" aria-labelledby="login-title">
	<h2 id="login-title">Admin login</h2>
	<form onsubmit={onLogin}>
		<label>
			<span>Username</span>
			<input bind:value={username} autocomplete="username" required />
		</label>
		<label>
			<span>Password</span>
			<input bind:value={password} autocomplete="current-password" type="password" required />
		</label>
		<button type="submit">Log in</button>
	</form>
</section>
