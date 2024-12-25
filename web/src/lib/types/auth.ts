import type { User } from "./user";

export interface DiscordAuthURLResponse {
	url: string;
}

export interface DiscordAuthResponse {
	token: string;
	user: User;
}
