import { atomWithQuery } from "jotai-tanstack-query";

import getUserInfo from "@/lib/api/users/get-info";
import type { User } from "@/lib/types/user";

import { isLoggedInAtom } from "./auth";

export const userDataAtom = atomWithQuery(
	get => ({
		queryKey: ["user-data"],
		queryFn: (): Promise<User> | null => {
			const isLoggedIn = get(isLoggedInAtom);

			if (!isLoggedIn) {
				return null;
			}

			return getUserInfo();
		},
		enabled: get(isLoggedInAtom),
	}),
);
