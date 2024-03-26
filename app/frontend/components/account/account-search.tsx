import { Session } from "next-auth";
import { AccountSearchInput } from "./account-search-input";
import { AccountSearchResult } from "./account-search-result";

type Props = {
	user: Session["user"];
	owner: string | undefined;
};

export const AccountSearch = ({ owner, user }: Props) => {
	return (
		<div className="group">
			<AccountSearchInput owner={owner} />
			<AccountSearchResult
				owner={owner}
				user={user}
				className="group-has-[[data-pending]]:animate-pulse"
			/>
		</div>
	);
};
