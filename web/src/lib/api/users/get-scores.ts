import type { ScoreResponse } from "@/lib/types/score";
import { type FolernResponse, ToAPIURL } from "@/lib/utils";

export default async function getUserScores(): Promise<ScoreResponse> {
	const res = await fetch(
		ToAPIURL("/api/users/me/scores"),
		{
			credentials: "include",
		},
	);

	const data = await res.json() as FolernResponse<ScoreResponse>;

	if (!res.ok) {
		throw new Error(data.error as string);
	}

	return data.body;
}
