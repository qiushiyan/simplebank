import NextAuth, { DefaultSession } from "next-auth";

type BaseUser = {
	name: string;
	access_token: string;
	email?: string;
};

declare module "next-auth" {
	interface Session extends DefaultSession {
		user: BaseUser;
	}

	interface User extends BaseUser {}
}
