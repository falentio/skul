export class SkulError extends Error {
	res?: Response
	body?: unknown
	constructor(msg: string, res?: Response, body?: unknown) {
		super(msg)
		this.res = res
		this.body = body
	}
}

export interface UserInfo {
	id: string
	name: string
	username: string
}

export interface Examination {
	id: string
	adminID: string

	name: string
	durationMinutes: number
	questionCount: number

	admin: UserInfo
}

export abstract class Client {
	#baseUrl?: URL
	constructor(baseUrl: URL) {
		this.#baseUrl = baseUrl
	}

	async fetch<T>(
		path: string, 
		method: string = "GET",
		body?: any,
	): Promise<T> {
		if (!this.#baseUrl) {
			throw new SkulError("baseUrl not set")
		}
		const url = new URL(path, this.#baseUrl).href
		console.log(url)
		const headers = new Headers()
		headers.set("user-agent", "skul-client")
		const init: RequestInit = {
			method,
			headers,
		}
		if (body) {
			if (method !== "GET") {
				init.body = JSON.stringify(body)
				headers.set("content-type", "application/json; charser=utf-8")
			} else {
				for (const [k, v] of Object.entries(body)) {
					url.searchParams.set(k, v)
				}
			}
		}
		const req = new Request(url, init)
		const res = await fetch(req)
		if (res.status === 429) {
			throw new SkulError("Too many request", res)
		}
		if (res.status === 204) {
			return
		}
		const resBody = await res.json()
		if (!res.ok) {
			throw new SkulError(resBody.message, res, resBody)
		}

		return resBody.data
	}

	websocket(path: string): Websocket {
		const url = new URL(path, this.#baseUrl)
		if (url.protocol === "http:") {
			url.protocol = "ws:"
		} else {
			url.protocol = "wss:"
		}
		return new Websocket(url.href)
	}
}

export class SkulClient extends Client {
	admin: AdminClient
	attendanceQr: AttendanceQrClient
	student: StudentClient

	constructor(baseUrl: URL) {
		super(baseUrl)

		this.admin = new AdminClient(baseUrl)
		this.student = new StudentClient(baseUrl)
	}

	health(): Promise<string> {
		return this.fetch<string>("health")
	}
}

export class AdminClient extends Client {
	constructor(baseUrl: URL) {
		baseUrl = new URL("admin/", baseUrl)
		super(baseUrl)
	}

	login(username: string, password: string): Promise<void> {
		return this.fetch<void>("login", "POST", {
			username,
			password,
		})
	}

	info() {
		return this.fetch<UserInfo>("info")
	}
}

export class ExaminationClient extends Client {
	constructor(baseUrl: URL) {
		baseUrl = new URL("examination/", baseUrl)
		super(baseUrl)
	}

	list(page: number = 1) {
		return this.fetch("list", "GET", {
			page,
		})
	}
}

export class StudentClient extends Client {
	constructor(baseUrl: URL) {
		baseUrl = new URL("student/", baseUrl)
		super(baseUrl)
	}

	login(username: string, password: string): Promise<void> {
		return this.fetch<void>("login", "POST", {
			username,
			password,
		})
	}

	info() {
		return this.fetch<UserInfo>("info")
	}
}

