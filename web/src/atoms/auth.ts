import { atomWithStorage } from "jotai/utils";
import { atomWithMutation } from "jotai-tanstack-query";

import { logout } from "@/lib/api/auth/logout";
import { processCallback } from "@/lib/api/auth/process-callback";
import type { AuthResponse } from "@/lib/types/auth";

export const isLoggedInAtom = atomWithStorage("auth-logged-in", false);

interface CallbackParams {
	code: string;
	state: string;
}

export const authCallbackAtom = atomWithMutation(
	() => ({
		mutationKey: ["auth-callback"],
		mutationFn: async ({ code, state }: CallbackParams): Promise<AuthResponse> => {
			if (!code || !state) {
				throw new Error("Invalid callback parameters.");
			}

			const result = await processCallback({ code, state });

			return result;
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
