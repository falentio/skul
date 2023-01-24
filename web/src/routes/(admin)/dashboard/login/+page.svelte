<script lang="ts">
	import { skul, refresh } from "$lib/skul"
	import { goto } from "$app/navigation"

	let username = ""
	let password = ""
	let error: Error | null  = null
	let loading = false
	$: disabled = loading || username === "" || password === "" || error !== null
	$: if (username || password) {
		error = null
	}

	async function login() {
		try {
			loading = true
			error = null

			await skul.admin.login(username, password)
			await refresh()

			username = ""
			password = ""

			goto("/dashboard")
		} catch (e) {
			if (e instanceof Error) {
				error = e
				return
			}
			error = new Error("unknown error")
			console.error(e)
		} finally {
			loading = false
		}
	}
</script>

<form 
	class="form-control bg-base-100 max-w-sm mx-auto p-4"
	on:submit|preventDefault={login}
>
	<label class="label"> Username </label>
	<input type="text" name="username" class="input input-bordered input-sm" bind:value={username} />

	<label class="label"> Password </label>
	<input type="password" name="password" class="input input-bordered input-sm" bind:value={password} />

	{#if error}
		<div class="alert alert-error mt-4">
			{error.message}
		</div>
	{/if}

	<button type="submit" class="mx-auto btn btn-primary btn-sm mt-4" {disabled}> Login </button>
</form>