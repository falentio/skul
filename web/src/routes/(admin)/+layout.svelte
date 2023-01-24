<script lang="ts">
	import Icon from "@iconify/svelte"
	import { goto, beforeNavigate } from "$app/navigation"
	import { admin, loading } from "$lib/skul"

	$: if (!$admin && !$loading) {
		goto("/dashboard/login")
	}

	beforeNavigate(({ to, cancel }) => {
		if (to.url.pathname !== "/dashboard/login" && !$admin && !$loading) {
			cancel()
			goto("/dashboard/login")
		}
	})
</script>

<div class="drawer">
	<input type="checkbox" class="drawer-toggle" id="nav-drawer">

	<div class="drawer-content flex flex-col bg-base-200">
		<div class="w-full navbar bg-base-100">
			<div class="flex-none">
				<label for="nav-drawer" class="btn btn-square btn-ghost">
					<Icon icon="mdi:menu" class="w-6 h-6" />
				</label>
			</div>
		</div>
		<div class="p-2">
			<slot />
		</div>
	</div>

	<div class="drawer-side">
		<label for="nav-drawer" class="drawer-overlay" />
		<div class="bg-base-100 w-min">
			<ul class="menu">
				<li>
					<a href="/dashboard"> Dashboard </a>
				</li>
				<li>
					<a href="/dashboard/login"> Login </a>
				</li>
			</ul>
		</div>
	</div>
</div>