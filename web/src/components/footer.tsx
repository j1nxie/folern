import Link from "next/link";

import { Button } from "@/components/ui/button";
import { NavigationMenu, NavigationMenuItem, NavigationMenuLink, NavigationMenuList } from "@/components/ui/navigation-menu";

export default function Footer(): React.JSX.Element {
	return (
		<footer className="w-full border-t border-border bg-background/95">
			<div className="mx-auto w-full">
				<NavigationMenu className="mx-auto flex w-full max-w-screen-2xl items-center px-6 py-4 text-sm">
					<div className="flex-1">Placeholder version string</div>

					<NavigationMenuList className="flex items-center gap-4 xl:gap-6">
						<NavigationMenuItem>
							<NavigationMenuLink asChild>
								<Button variant="link" className="p-0" asChild>
									<Link href="https://github.com/j1nxie/folern">Source code</Link>
								</Button>
							</NavigationMenuLink>
						</NavigationMenuItem>
						<NavigationMenuItem>
							<NavigationMenuLink asChild>
								<Button variant="link" className="p-0" asChild>
									<Link href="/credits">Credits</Link>
								</Button>
							</NavigationMenuLink>
						</NavigationMenuItem>
					</NavigationMenuList>
				</NavigationMenu>
			</div>
		</footer>
	);
}
