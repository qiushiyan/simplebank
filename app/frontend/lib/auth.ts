import { getServerSession } from "next-auth";
import { z } from "zod";
import { env } from "./env.mjs";
import { AuthInput, authOptions } from "./nextauth";

const UserResponseSchema = z.object({
	username: z.string(),
	email: z.string().nullable(),
	created_at: z.string(),
	password_changed_at: z.string(),
});

const UserSignupResponseSchema = z.union([
	z.object({
		data: z.object({
			user: UserResponseSchema,
			access_token: z.string(),
		}),
	}),
	z.object({
		error: z.string(),
		fields: z.record(z.string()).optional(),
	}),
]);

export const authenticate = async (input: AuthInput) => {
	if (input.username === env.ADMIN_USERNAME) {
		input.password = env.ADMIN_PASSWORD;
	} else if (input.username === env.USER_USERNAME) {
		input.password = env.USER_PASSWORD;
	}

	const response = await fetch(
		input.endpoint === "signin"
			? `${env.API_URL}/signin`
			: `${env.API_URL}/signup`,
		{
			method: "POST",
			headers: {
				"Content-Type": "application/json",
			},
			body: JSON.stringify({
				username: input.username,
				password: input.password,
				endpoint: input.endpoint,
			}),
		},
	);

	const data = await response.json();
	const parsed = UserSignupResponseSchema.safeParse(data);
	return parsed.success ? parsed.data : null;
};

export const getCurrentUser = async () => {
	const session = await getServerSession(authOptions);
	return session?.user;
};
