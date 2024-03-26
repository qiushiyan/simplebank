import { searchAccounts } from "@/lib/account";
import { Session } from "next-auth";

type Props = {
	user: Session["user"];
	owner: string | undefined;
};

export const AccountSearchList = async ({ user, owner }: Props) => {
	if (!owner) {
		return null;
	}

	const data = await searchAccounts(owner, user);
	if (!data) {
		return <p>fail to get accounts</p>;
	}
	return (
		<div className="space-y-4">
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
