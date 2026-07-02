type StreamEnvelope<T> = {
	data: T;
};

export function parseSystemEvent<T>(event: Event) {
	try {
		return (JSON.parse((event as MessageEvent<string>).data) as StreamEnvelope<T>).data;
	} catch {
		return undefined;
	}
}
