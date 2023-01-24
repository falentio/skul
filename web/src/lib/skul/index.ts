import { browser } from "$app/environment"
import { writable, derived } from "svelte/store"
import type { UserInfo } from "./skul"
import { SkulClient } from "./skul"

const href = new URL(browser ? window.location.href : "http://localhost:8080")
export const loading = writable<boolean>(true)
export const baseUrl = new URL("/api/", href)
export const skul = new SkulClient(baseUrl)
export const user = writable<UserInfo | null>(null)

export const admin = derived(user, ($user) => {
	return !!$user?.id.startsWith("adm")
})

export const student = derived(user, ($user) => {
	return !!$user?.id.startsWith("stu")
})

export async function refresh(): Promise<void> {
	const pendings = [
		skul.admin.info().then(user.set),
		skul.student.info().then(user.set),
	]
	await Promise.allSettled(pendings)
	loading.set(false)
}

if (browser) {
	refresh()
}
