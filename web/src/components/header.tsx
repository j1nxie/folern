"use client";

import Link from "next/link";
import { useState } from "react";

import { Button } from "@/components/ui/button";
import { NavigationMenu, NavigationMenuItem, NavigationMenuLink, NavigationMenuList } from "@/components/ui/navigation-menu";
import { getAuthURL } from "@/lib/api/auth/get-url";
import { logout } from "@/lib/api/auth/logout";

export default function Header(): React.JSX.Element {
	const [loggedIn, _] = useState(localStorage.getItem("LOGGED_IN"));

	return (
		<header className="sticky top-0 z-50 w-full border-b border-border bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
			<div className="mx-auto w-full">
				<NavigationMenu className="mx-auto flex w-full max-w-screen-2xl items-center px-6 py-4">
					<NavigationMenuList className="flex items-center gap-4 xl:gap-6">
						<NavigationMenuItem className="lg:mr-4">
							<NavigationMenuLink asChild>
								<Button variant="link" className="p-0" asChild>
									<Link href="/">folern</Link>
								</Button>
							</NavigationMenuLink>
						</NavigationMenuItem>

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

					<div className="flex flex-1 items-center justify-between gap-2 md:justify-end">
						{loggedIn !== "true" && (
							<Button onClick={async () => {
								const { url } = await getAuthURL();

								window.location.href = url;
							}}
							>
								Login
							</Button>
						)}

						{loggedIn === "true" && (
							<Button
								onClick={async () => {
									await logout();
									localStorage.removeItem("LOGGED_IN");
									window.location.href = "/";
								}}
								variant="destructive"
							>
								Logout
							</Button>
						)}
					</div>
				</NavigationMenu>
			</div>
		</header>
	);
}
