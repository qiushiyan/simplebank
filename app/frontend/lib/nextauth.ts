import { AuthOptions } from "next-auth";
import Credentials from "next-auth/providers/credentials";
import { z } from "zod";
import { authenticate } from "./auth";

const QuerySchema = z.object({
	username: z.string(),
	password: z.string().optional(),
	email: z.string().optional(),
	endpoint: z.union([z.literal("signup"), z.literal("signin")]),
});

export type AuthInput = z.infer<typeof QuerySchema>;

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

				const response = await authenticate(query.data);
				if (!response) {
					return null;
				}

				if ("error" in response) {
					if (response.fields) {
						throw new Error(
							JSON.stringify({
								message: response.error,
								details: response.fields,
							}),
						);
					}

					throw new Error(
						JSON.stringify({
							message: response.error,
						}),
					);
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
