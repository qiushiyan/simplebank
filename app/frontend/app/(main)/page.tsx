import { getCurrentUser } from "@/lib/auth";

export default async function () {
	const user = await getCurrentUser();
	if (!user) {
		return null;
	}

	return (
		<div>
			<p>hello {user.name}</p>
		</div>
	);
}
