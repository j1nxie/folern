import type { DiscordAuthURLResponse } from "@/lib/types/auth";
import { type FolernResponse, ToAPIURL } from "@/lib/utils";

export async function getAuthURL(): Promise<DiscordAuthURLResponse> {
	const res = await fetch(ToAPIURL("/api/auth/url"), { credentials: "include" });
	const data = await res.json() as FolernResponse<DiscordAuthURLResponse>;

	if (!res.ok) {
		throw new Error(data.error as string);
	}

	return data.data;
}
