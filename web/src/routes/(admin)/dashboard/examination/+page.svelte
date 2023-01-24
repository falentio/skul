<script lang="ts">
	import type { Examination }	from "$lib/skul"
	import { skul }	from "$lib/skul"
	import { onMount } from "svelte"

	let exam: Examination[] = []

	onMount(() => {
		skul.examination.list().then(e => exam = e)
	})
</script>

<div class="flex flex-col space-y-4"> 
	<div class="bg-base-100 card card-compact w-full">
		<div class="card-body">
			<h2 class="card-title"> Ujian </h2>
			<div class="overflow-x-auto">
				<table class="table table-compact w-full">
					<thead>
						<tr>
							<th />
							<th> Nama </th>
							<th> Durasi </th>
							<th> Jumlah Jawaban </th>
							<th> Aksi </th>
						</tr>
					</thead>
				</table>
			</div>
		</div>
	</div>

	<div class="bg-base-100 card card-compact w-full">
		<div class="card-body">
			<h2 class="card-title"> Baru </h2>

			<form class="form-control bg-base-100">
				<label class="label"> Nama </label>
				<input type="text" name="name" class="input input-bordered input-sm" />

				<label class="label"> Durasi </label>
				<input type="number" name="durasi" class="input input-bordered input-sm" />

				<label class="label"> Jumlah Jawaban </label>
				<input type="number" name="name" class="input input-bordered input-sm" />

				<button class="btn btn-primary btn-sm mt-4" type="submit"> Buat </button>
			</form>
		</div>
	</div>
</div>
<pre>{JSON.stringify(exam, null, "\t")}</pre>


