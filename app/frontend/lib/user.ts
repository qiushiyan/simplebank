export const testUsers = [
	{
		name: "Admin",
		explanation:
			"Admin account with full access. Can view and manage all accounts, transfers, etc.",
	},
	{
		name: "User",
		explanation: "Regular user account. Can only manage personal data",
	},
] as const;

export type TestUsername = (typeof testUsers)[number]["name"];
