import { atomWithQuery } from "jotai-tanstack-query";

import getUserInfo from "@/lib/api/users/get-info";
import getUserScores from "@/lib/api/users/get-scores";
import getUserStats from "@/lib/api/users/get-stats";
import type { ScoreResponse, StatsResponse } from "@/lib/types/score";
import type { User } from "@/lib/types/user";

import { isLoggedInAtom } from "./auth";

export const userDataAtom = atomWithQuery(
	get => ({
		queryKey: ["user-data"],
		queryFn: async (): Promise<User | null> => {
			const isLoggedIn = get(isLoggedInAtom);

			if (!isLoggedIn) {
				return null;
			}

			const result = await getUserInfo();

			return result;
		},
		enabled: get(isLoggedInAtom),
	}),
);

export const userScoresAtom = atomWithQuery(
	get => ({
		queryKey: ["user-scores"],
		queryFn: async (): Promise<ScoreResponse | null> => {
			const isLoggedIn = get(isLoggedInAtom);

			if (!isLoggedIn) {
				return null;
			}

			const result = await getUserScores();

			return result;
		},
		enabled: get(isLoggedInAtom),
	}),
);

export const userStatsAtom = atomWithQuery(
	get => ({
		queryKey: ["user-stats"],
		queryFn: async (): Promise<StatsResponse | null> => {
			const isLoggedIn = get(isLoggedInAtom);

			if (!isLoggedIn) {
				return null;
			}

			const result = await getUserStats();

			return result;
		},
		enabled: get(isLoggedInAtom),
	}),
);
