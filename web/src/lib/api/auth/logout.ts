import { type FolernResponse, ToAPIURL } from "@/lib/utils";

export async function logout(): Promise<void> {
	const res = await fetch(ToAPIURL("/api/auth/logout"), { credentials: "include" });
	const data = await res.json() as FolernResponse<string>;

	if (!res.ok) {
		throw new Error(data.error as string);
	}
}
