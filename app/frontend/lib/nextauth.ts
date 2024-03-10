import { AuthOptions, getServerSession } from "next-auth";
import Credentials from "next-auth/providers/credentials";
import { z } from "zod";
import { getUser } from "./auth";
import { env } from "./env.mjs";
import { TestUsername } from "./user";

const QuerySchema = z.object({
	username: z.string(),
	password: z.string().optional(),
});

export const authOptions: AuthOptions = {
	providers: [
		Credentials({
			name: "Credentials",
			credentials: {},
			authorize: async (credentials, req) => {
				const query = QuerySchema.safeParse(req.query);
				if (!query.success) {
					return null;
				}

				const { username, password } = query.data;
				let response: Awaited<ReturnType<typeof getUser>>;
				if (!password) {
					response = await getUser({ username: username as TestUsername });
				} else {
					response = await getUser({ username, password });
				}
				if (!response) {
					return null;
				}

				return {
					id: "",
					name: response.data.user.username,
					email: response.data.user.email || undefined,
					access_token: response.data.access_token,
				};
			},
		}),
	],
	callbacks: {
		jwt: async ({ token, user }) => {
			if (user) {
				token.access_token = user.access_token;
			}
			return token;
		},
		session: async ({ session, token, user }) => {
			if (session.user) {
				session.user.access_token = token.access_token as string;
			}
			return session;
		},
	},
	pages: {
		signIn: "/auth?tab=signin",
		newUser: "/auth?tab=signup",
		error: "/auth",
	},
};
