interface MediaPlaybackErrorSource {
	error?: { code: number } | null;
}

export function mediaPlaybackErrorMessage(video?: MediaPlaybackErrorSource) {
	const error = video?.error;
	if (!error) return browserPlaybackHint('The browser stopped playback before reporting a reason.');
	const reason = mediaErrorReason(error.code);
	return browserPlaybackHint(reason);
}

function mediaErrorReason(code: number) {
	switch (code) {
		case 1:
			return 'The browser aborted the preview request.';
		case 2:
			return 'The browser lost the preview stream while loading it.';
		case 3:
			return 'The browser could not decode the generated MP4 stream.';
		case 4:
			return 'The browser does not support this preview stream.';
		default:
			return 'The browser could not play this preview stream.';
	}
}

function browserPlaybackHint(reason: string) {
	return `${reason} Try Chrome, Brave, or Edge for the widest preview support, or use VLC to stream the original file.`;
}
