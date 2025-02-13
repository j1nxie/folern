import type { DiscordAuthResponse } from "@/lib/types/auth";
import { type FolernResponse, ToAPIURL } from "@/lib/utils";

export async function processDiscordCallback({ code, state }: { code: string; state: string }): Promise<DiscordAuthResponse> {
	const res = await fetch(
		`${ToAPIURL("/api/auth/callback")}?code=${code}&state=${state}`,
		{
			credentials: "include",
		},
	);
	const data = await res.json() as FolernResponse<DiscordAuthResponse>;

	if (!res.ok) {
		throw new Error(data.error as string);
	}

	return data.body;
}
