import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogFooter,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
} from "@/components/ui/dialog";
import { BaseUser } from "@/index";
import { getReceived, getSent } from "@/lib/friendship";
import { delay } from "@/lib/utils";
import { PlusIcon } from "lucide-react";
import { Spinner } from "../ui/spinner";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "../ui/tabs";

type Props = {
	user: BaseUser;
	accountId: number;
};

type FriendData = Awaited<ReturnType<typeof getReceived>>;

const SentList = ({ data }: { data: FriendData }) => {
	if (!data) {
		return <p className="text-muted-foreground">Failed to get data</p>;
	}

	return data.data.map((friendship) => (
		<div key={friendship.id} className="flex items-center space-x-4">
			<img
				src={`https://api.dicebear.com/8.x/lorelei/svg?seed=${friendship.pending}`}
				alt="avatar"
				className="rounded-full size-16"
			/>
			<span>{friendship.to_account_id}</span>
			<button className="btn btn-secondary">Cancel</button>
		</div>
	));
};

const ReceivedList = ({ data }: { data: FriendData }) => {};

export const Friendship = async ({ accountId, user }: Props) => {
	const [received, sent] = await Promise.all([
		getReceived(user, accountId),
		getSent(user, accountId),
	]);

	return (
		<Dialog>
			<DialogTrigger asChild>
				<button
					className="fixed bottom-16 right-12 rounded-full border border-primary hover:bg-primary hover:border-border hover:shadow-lg ease-linear duration-200 p-4"
					type="button"
				>
					<PlusIcon className="w-6 h-6" />
				</button>
			</DialogTrigger>
			<DialogContent className="sm:max-w-[425px]">
				<DialogHeader>
					<DialogTitle>Friends</DialogTitle>
					<DialogDescription>View friendship requests</DialogDescription>
				</DialogHeader>
				<Tabs defaultValue="received">
					<TabsList>
						<TabsTrigger value="received">Received</TabsTrigger>
						<TabsTrigger value="sent">Sent</TabsTrigger>
					</TabsList>
					<TabsContent value="received">
						{received?.data.map((friendship) => (
							<div key={friendship.id} className="flex items-center space-x-4">
								<span>{friendship.from_account_id}</span>
								<button className="btn btn-primary">Accept</button>
								<button className="btn btn-secondary">Decline</button>
							</div>
						))}
					</TabsContent>
					<TabsContent value="sent">
						<SentList data={sent} />
					</TabsContent>
				</Tabs>
			</DialogContent>
		</Dialog>
	);
};

Friendship.Skeleton = () => (
	<button
		className="fixed bottom-16 right-12 rounded-full border border-primary hover:bg-primary hover:border-border hover:shadow-lg ease-linear duration-200 p-4"
		type="button"
	>
		<Spinner className="size-6" />
	</button>
);
