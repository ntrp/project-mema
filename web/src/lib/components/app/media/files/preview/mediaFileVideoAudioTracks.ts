import type Player from 'video.js/dist/types/player';
import type { AudioTrackOption } from '$lib/components/app/media/files/preview/mediaFilePlayback';

interface VideoJsAudioTrack {
	id: string;
	enabled: boolean;
}

interface VideoJsAudioTrackList {
	length: number;
	addTrack: (_track: unknown) => void;
	on: (_type: string, _handler: () => void) => void;
	off: (_type: string, _handler: () => void) => void;
	[index: number]: VideoJsAudioTrack;
}

export function addAudioTracks(
	videojs: (typeof import('video.js'))['default'],
	instance: Player,
	tracks: AudioTrackOption[],
	activeKey: string,
	onAudioTrackChange: (_key: string) => void
) {
	if (tracks.length === 0) return undefined;
	const Track = (
		videojs as unknown as {
			AudioTrack: new (_options: Record<string, unknown>) => unknown;
		}
	).AudioTrack;
	const list = instance.audioTracks() as unknown as VideoJsAudioTrackList;
	const enabledKey = activeKey || tracks.find((track) => track.enabled)?.key || tracks[0]?.key;
	for (const track of tracks) {
		list.addTrack(
			new Track({
				id: track.key,
				kind: track.key === enabledKey ? 'main' : 'alternative',
				label: track.label,
				language: track.language ?? '',
				enabled: track.key === enabledKey
			})
		);
	}
	const handleChange = () => {
		for (let index = 0; index < list.length; index += 1) {
			const track = list[index];
			if (track?.enabled && track.id !== activeKey) {
				onAudioTrackChange(track.id);
				return;
			}
		}
	};
	list.on('change', handleChange);
	return () => list.off('change', handleChange);
}
