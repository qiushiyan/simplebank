export default function NotFound() {
	return (
		<div className="grid place-items-center px-6 py-24 sm:py-32 lg:px-8">
			<div className="text-center">
				<p className="text-base font-semibold">404</p>
				<h1 className="mt-4 text-3xl font-bold tracking-tight text-primary/80 sm:text-5xl">
					Page not found or the API is down
				</h1>
				<p className="mt-6 text-base leading-7 text-gray-600">
					Make sure the API is running and you have set the correct environment
					variables.
				</p>
				<div className="mt-10 flex items-center justify-center gap-x-6">
					<a href="/" className="text-sm font-semibold ">
						Go to home <span aria-hidden="true">&rarr;</span>
					</a>
				</div>
			</div>
		</div>
	);
}
