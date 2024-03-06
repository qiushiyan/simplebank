"use client";
import { cn } from "@/lib/utils";
import { usePathname, useRouter } from "next/navigation";
import { useTransition } from "react";
import { Button } from "../ui/button";

type Props = {
	id: number;
	name: string;
};

export const AccountLink = ({ id, name }: Props) => {
	const [isPending, startTransition] = useTransition();
	const router = useRouter();
	const pathname = usePathname();
	const active = pathname.includes(`/account/${id}`);

	return (
		<Button
			data-pending={isPending ? "" : undefined}
			variant="link"
			className={cn("justify-start", {
				"bg-primary/80 text-foreground": active,
				"opacity-50": isPending,
			})}
			onClick={() => {
				startTransition(async () => {
					router.push(`/account/${id}`);
				});
			}}
		>
			{name}
		</Button>
	);
};
