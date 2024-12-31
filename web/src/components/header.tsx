"use client";

import { useAtom } from "jotai";
import Image from "next/image";
import Link from "next/link";
import { useRouter } from "next/navigation";

import { authDiscordCallback, isLoggedInAtom, logoutAtom } from "@/atoms/auth";
import { Button } from "@/components/ui/button";
import { NavigationMenu, NavigationMenuItem, NavigationMenuLink, NavigationMenuList } from "@/components/ui/navigation-menu";
import { getAuthURL } from "@/lib/api/auth/get-url";

import { ModeToggle } from "./dark-mode-toggle";

export default function Header(): React.JSX.Element {
	const router = useRouter();
	const [isLoggedIn, setIsLoggedIn] = useAtom(isLoggedInAtom);
	const [{ status: loginStatus }] = useAtom(authDiscordCallback);
	const [{ mutate: processLogout }] = useAtom(logoutAtom);

	return (
		<header className="sticky top-0 z-50 w-full border-b border-border bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
			<div className="mx-auto w-full">
				<NavigationMenu className="mx-auto flex w-full max-w-screen-2xl items-center px-6 py-4">
					<NavigationMenuList className="flex items-center gap-4 xl:gap-6">
						<Link href="/">
							<Image src="/folern.png" alt="folern" width={48} height={48} className="rounded-xl" priority />
						</Link>

						{isLoggedIn && loginStatus !== "pending" && (
							<NavigationMenuItem>
								<NavigationMenuLink asChild>
									<Button variant="link" className="p-0" asChild>
										<Link href="/profile">Profile</Link>
									</Button>
								</NavigationMenuLink>
							</NavigationMenuItem>
						)}
						<NavigationMenuItem>
							<NavigationMenuLink asChild>
								<Button variant="link" className="p-0" asChild>
									<Link href="/tools">Tools</Link>
								</Button>
							</NavigationMenuLink>
						</NavigationMenuItem>
						<NavigationMenuItem>
							<NavigationMenuLink asChild>
								<Button variant="link" className="p-0" asChild>
									<Link href="/help">Help</Link>
								</Button>
							</NavigationMenuLink>
						</NavigationMenuItem>
					</NavigationMenuList>

					<div className="flex flex-1 items-center justify-end gap-2">
						{!isLoggedIn && loginStatus !== "pending" && (
							<Button onClick={async () => {
								const { url } = await getAuthURL();

								window.location.href = url;
							}}
							>
								Login
							</Button>
						)}

						{isLoggedIn && loginStatus !== "pending" && (
							<Button
								onClick={() => {
									processLogout(undefined, {
										onSuccess: () => {
											setIsLoggedIn(false);
											router.replace("/");
										},
									});
								}}
								variant="destructive"
							>
								Logout
							</Button>
						)}

						<ModeToggle />
					</div>
				</NavigationMenu>
			</div>
		</header>
	);
}
