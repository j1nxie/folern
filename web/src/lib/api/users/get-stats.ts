import type { StatsResponse } from "@/lib/types/score";
import { type FolernResponse, ToAPIURL } from "@/lib/utils";

export default async function getUserStats(): Promise<StatsResponse> {
	const res = await fetch(
		ToAPIURL("/api/users/me/stats"),
		{
			credentials: "include",
		},
	);

	const data = await res.json() as FolernResponse<StatsResponse>;

	if (!res.ok) {
		throw new Error(data.error as string);
	}

	return data.body;
}
