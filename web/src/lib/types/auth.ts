import type { User } from "./user";

export interface AuthURLResponse {
	url: string;
}

export interface AuthResponse {
	token: string;
	user: User;
}
