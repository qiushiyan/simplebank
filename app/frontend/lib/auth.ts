import { authOptions } from "@/app/api/auth/[...nextauth]/route";
import { TestUsername } from "@/lib/user";
import { getServerSession } from "next-auth";
import { z } from "zod";
import { env } from "./env.mjs";

const TestUsernameSchema = z.custom<TestUsername>();

const UserResponseSchema = z.object({
	username: z.string(),
	email: z.string().nullable(),
	created_at: z.string(),
	password_changed_at: z.string(),
});

const GetUserResponseSchema = z.object({
	data: z.object({
		user: UserResponseSchema,
		access_token: z.string(),
	}),
});

type GetUserInput =
	| {
			username: TestUsername;
	  }
	| {
			username: string;
			password: string;
	  };

export const getUser = async (input: GetUserInput) => {
	const req = {
		username: "",
		password: "",
	};
	if (!("password" in input)) {
		const usernameParsed = TestUsernameSchema.safeParse(input.username);
		if (usernameParsed.success) {
			if (usernameParsed.data === "Admin") {
				req.username = env.ADMIN_USERNAME;
				req.password = env.ADMIN_PASSWORD;
			} else if (usernameParsed.data === "User") {
				req.username = env.USER_USERNAME;
				req.password = env.USER_PASSWORD;
			}
		}
	} else {
		req.username = input.username;
		req.password = input.password;
	}

	const response = await fetch(`${env.API_URL}/signin`, {
		method: "POST",
		headers: {
			"Content-Type": "application/json",
		},
		body: JSON.stringify(req),
	});

	if (!response.ok) {
		return null;
	}

	const data = await response.json();
	const parsed = GetUserResponseSchema.safeParse(data);
	return parsed.success ? parsed.data : null;
};

export const getCurrentUser = async () => {
	const session = await getServerSession(authOptions);
	return session?.user;
};
