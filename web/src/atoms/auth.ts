import { atomWithStorage } from "jotai/utils";
import { atomWithMutation } from "jotai-tanstack-query";

import { logout } from "@/lib/api/auth/logout";
import { processDiscordCallback } from "@/lib/api/auth/process-discord-callback";
import { processKamaitachiCallback } from "@/lib/api/auth/process-kamaitachi-callback";
import type { AuthResponse } from "@/lib/types/auth";

export const isLoggedInAtom = atomWithStorage("auth-logged-in", false);

interface CallbackParams {
	code: string;
	state: string;
}

export const authDiscordCallback = atomWithMutation(
	() => ({
		mutationKey: ["auth-discord-callback"],
		mutationFn: async ({ code, state }: CallbackParams): Promise<AuthResponse> => {
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
			return true;
		},
	}),
);
