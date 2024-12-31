import type { User } from "@/lib/types/user";
import { type FolernResponse, ToAPIURL } from "@/lib/utils";

export default async function getUserInfo(): Promise<User> {
	const res = await fetch(
		ToAPIURL("/api/users/me"),
		{
			credentials: "include",
		},
	);

	const data = await res.json() as FolernResponse<User>;

	if (!res.ok) {
		throw new Error(data.error as string);
	}

	return data.body;
}
