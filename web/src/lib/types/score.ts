export interface Chart {
	id: string;
	song_id: number;
	level: string;
}

export interface Song {
	id: number;
	title: string;
	artist: string;
	version: string;
	genre: string;
}

export interface StatsResponse {
	all: string;
	genre: Record<string, string>;
	version: Record<string, string>;
}

// TODO: i can probably make it even more typesafe by making versions / genres enums
export interface ScoreResponse {
	genre: Record<string, ScoreDocument[]>;
	version: Record<string, ScoreDocument[]>;
}

export interface ScoreDocument {
	chart_id: string;
	song_id: number;
	score: number;
	lamp: ScoreLamp;
	over_power: string;
	chart: Chart;
	song: Song;
}

export enum ScoreLamp {
	FAILED = "FAILED",
	CLEAR = "CLEAR",
	FULL_COMBO = "FULL COMBO",
	ALL_JUSTICE = "ALL JUSTICE",
	ALL_JUSTICE_CRITICAL = "ALL JUSTICE CRITICAL",
}
