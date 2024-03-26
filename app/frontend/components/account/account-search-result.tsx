import { searchAccounts } from "@/lib/account";
import { cn, delay } from "@/lib/utils";
import { Session } from "next-auth";

type Props = {
	user: Session["user"];
	owner: string | undefined;
	className?: string;
};

export const AccountSearchResult = async ({
	user,
	owner,
	className,
}: Props) => {
	if (!owner) {
		return null;
	}
	await delay(1000);
	const data = await searchAccounts(owner, user);
	if (!data) {
		return <p>fail to get accounts</p>;
	}

	if (data.data.length === 0) {
		return <p>no accounts found</p>;
	}

	return (
		<div className={cn("space-y-4", className)}>
			{data.data.map((account) => (
				<div
					key={account.id}
					className="flex gap-4 p-4 border border-gray-200 rounded"
				>
					<p>{account.owner}</p>
					<p>{account.name}</p>
					<p>{account.currency}</p>
				</div>
			))}
		</div>
	);
};
