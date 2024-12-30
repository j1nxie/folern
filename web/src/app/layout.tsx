import "./globals.css";

import type { Metadata } from "next";
import { Geist, Geist_Mono as GeistMono } from "next/font/google";

import Footer from "@/components/footer";
import Header from "@/components/header";
import { LayoutWrapper } from "@/components/layout-wrapper";

import { Providers } from "./providers";

const geistSans = Geist({
	variable: "--font-geist-sans",
	subsets: ["latin"],
});

const geistMono = GeistMono({
	variable: "--font-geist-mono",
	subsets: ["latin"],
});

export const metadata: Metadata = {
	title: "folern",
	description: "An unofficial tool for tracking and calculating OVER POWER stats for CHUNITHM.",
};

export default function RootLayout({
	children,
}: Readonly<{
	children: React.ReactNode;
}>): React.JSX.Element {
	return (
		<html lang="en" suppressHydrationWarning>
			<body
				className={`${geistSans.variable} ${geistMono.variable} antialiased`}
			>
				<Providers>
					<div className="relative flex min-h-svh flex-col bg-background font-sans">
						<LayoutWrapper>
							<Header />

							<main className="flex grow">{children}</main>

							<Footer />
						</LayoutWrapper>
					</div>
				</Providers>
			</body>
		</html>
	);
}
