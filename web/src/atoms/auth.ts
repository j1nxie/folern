import { atomWithStorage } from "jotai/utils";
import { atomWithMutation, atomWithQuery } from "jotai-tanstack-query";

import { logout } from "@/lib/api/auth/logout";
import { processDiscordCallback } from "@/lib/api/auth/process-discord-callback";
import { processKamaitachiCallback } from "@/lib/api/auth/process-kamaitachi-callback";
import getUserInfo from "@/lib/api/users/get-info";
import type { DiscordAuthResponse } from "@/lib/types/auth";
import type { User } from "@/lib/types/user";

export const isLoggedInAtom = atomWithStorage("auth-logged-in", false);

interface CallbackParams {
	code: string;
	state: string;
}

export const authVerificationAtom = atomWithQuery(
	get => ({
		queryKey: ["auth-verification"],
		queryFn: async (): Promise<User> => {
			try {
				const response = await getUserInfo();
				return response;
			} catch (error) {
				localStorage.removeItem("auth-logged-in");
				throw error;
			}
		},
		refetchInterval: 5 * 60 * 1000,
		enabled: get(isLoggedInAtom),
	}),
);

export const authDiscordCallback = atomWithMutation(
	() => ({
		mutationKey: ["auth-discord-callback"],
		mutationFn: async ({ code, state }: CallbackParams): Promise<DiscordAuthResponse> => {
			if (!code || !state) {
				throw new Error("Invalid callback parameters.");
			}

			const result = await processDiscordCallback({ code, state });

			return result;
		},
	}),
);

export const authKamaitachiCallback = atomWithMutation(
	() => ({
		mutationKey: ["auth-kamaitachi-callback"],
		mutationFn: async ({ code }: { code: string }): Promise<void> => {
			if (!code) {
				throw new Error("Invalid callback parameters.");
			}

			await processKamaitachiCallback({ code });
		},
	}),
);

export const logoutAtom = atomWithMutation(
	() => ({
		mutationKey: ["auth-logout"],
		mutationFn: async (): Promise<boolean> => {
			await logout();
			localStorage.removeItem("auth-logged-in");
			return true;
		},
	}),
);
