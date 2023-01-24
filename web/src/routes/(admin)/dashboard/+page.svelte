<script lang="ts">
	import { skul, admin, user } from "$lib/skul"
</script>

<div class="flex flex-col"> 
	<div class="bg-base-100 card card-compact md:w-1/2 w-full">
		<div class="card-body">
			<h2 class="card-title"> Dashboard </h2>
			<div> Hai {$user?.name} </div>
		</div>
	</div>
</div>