import { createEnv } from "@t3-oss/env-nextjs";
import { z } from "zod";

export const env = createEnv({
	server: {
		NEXTAUTH_URL: z.string(),
		NEXTAUTH_SECRET: z.string(),
		API_URL: z.string(),
		ADMIN_USERNAME: z.string(),
		ADMIN_PASSWORD: z.string(),
		USER_USERNAME: z.string(),
		USER_PASSWORD: z.string(),
	},
	runtimeEnv: {
		NEXTAUTH_URL: process.env.NEXTAUTH_URL,
		NEXTAUTH_SECRET: process.env.NEXTAUTH_SECRET,
		API_URL: process.env.API_URL,
		ADMIN_USERNAME: process.env.ADMIN_USERNAME,
		ADMIN_PASSWORD: process.env.ADMIN_PASSWORD,
		USER_USERNAME: process.env.USER_USERNAME,
		USER_PASSWORD: process.env.USER_PASSWORD,
	},
});
