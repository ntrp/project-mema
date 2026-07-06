export const videoCodecOptions = [
	{ value: 'h264', label: 'H.264' },
	{ value: 'hevc', label: 'H.265 / HEVC' },
	{ value: 'av1', label: 'AV1' },
	{ value: 'vp9', label: 'VP9' },
	{ value: 'mpeg2video', label: 'MPEG-2' },
	{ value: 'vc1', label: 'VC-1' },
	{ value: 'prores', label: 'ProRes' }
];

export const hdrFormatOptions = [
	{ value: 'sdr', label: 'SDR' },
	{ value: 'hdr10', label: 'HDR10' },
	{ value: 'hdr10plus', label: 'HDR10+' },
	{ value: 'dolby-vision', label: 'Dolby Vision' },
	{ value: 'hlg', label: 'HLG' }
];

export const pixelFormatOptions = [
	{ value: 'yuv420p', label: 'YUV 4:2:0 8-bit' },
	{ value: 'yuv420p10le', label: 'YUV 4:2:0 10-bit' },
	{ value: 'yuv422p10le', label: 'YUV 4:2:2 10-bit' },
	{ value: 'yuv444p10le', label: 'YUV 4:4:4 10-bit' },
	{ value: 'yuv420p12le', label: 'YUV 4:2:0 12-bit' }
];

export const audioCodecOptions = [
	{ value: 'aac', label: 'AAC' },
	{ value: 'ac3', label: 'AC-3' },
	{ value: 'eac3', label: 'E-AC-3' },
	{ value: 'dts', label: 'DTS' },
	{ value: 'truehd', label: 'TrueHD' },
	{ value: 'flac', label: 'FLAC' },
	{ value: 'opus', label: 'Opus' },
	{ value: 'mp3', label: 'MP3' }
];

export const audioChannelOptions = [
	{ value: '1.0', label: '1.0 Mono' },
	{ value: '2.0', label: '2.0 Stereo' },
	{ value: '3.0', label: '3.0' },
	{ value: '5.1', label: '5.1' },
	{ value: '6.1', label: '6.1' },
	{ value: '7.1', label: '7.1' },
	{ value: 'atmos', label: 'Atmos' }
];

export const subtitleFormatOptions = [
	{ value: 'srt', label: 'SRT' },
	{ value: 'ass', label: 'ASS' },
	{ value: 'ssa', label: 'SSA' },
	{ value: 'vtt', label: 'WebVTT' },
	{ value: 'pgs', label: 'PGS' },
	{ value: 'subrip', label: 'SubRip' }
];
