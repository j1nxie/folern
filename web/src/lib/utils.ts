import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]): string {
	return twMerge(clsx(inputs));
}

const BASE_URL = process.env.NODE_ENV === "development" ? "https://localhost:8081" : "https://folern.rylie.moe";

export function ToAPIURL(path: string, queryParams?: Record<string, string>): string {
	const url_ = path.startsWith("/") ? path : `/${path}`;
	const url = new URL(url_, BASE_URL);

	if (queryParams) {
		for (const [key, value] of Object.entries(queryParams)) {
			if (value) {
				url.searchParams.set(key, value);
			}
		}
	}

	return url.toString();
}

export interface FolernResponse<T = unknown, E = unknown> {
	success: boolean;
	body: T;
	error: E;
	cat: string;
}
