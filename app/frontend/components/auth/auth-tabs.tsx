import { Button } from "../ui/button";

import { Input } from "../ui/input";
import { Label } from "../ui/label";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "../ui/tabs";
import { SignIn } from "./signin";
import { SignUp } from "./signup";

type Props = {
	tab: "signin" | "signup";
};

export const AuthTabs = ({ tab }: Props) => {
	return (
		<section>
			<Tabs defaultValue={tab}>
				<TabsList className="grid w-full grid-cols-2">
					<TabsTrigger value="signin">Sign in</TabsTrigger>
					<TabsTrigger value="signup">Sign up</TabsTrigger>
				</TabsList>
				<TabsContent value="signin" className="">
					<SignIn />
				</TabsContent>
				<TabsContent value="signup">
					<SignUp />
				</TabsContent>
			</Tabs>
		</section>
	);
};
