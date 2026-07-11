import { describe, expect, it } from 'vitest';
import {
	basename,
	buildTraceSteps,
	buildTraceSummary,
	filterTraceSteps
} from './dlnaDecisionTrace';
import type { DLNADeliveryTraceResponse, DLNAProfileMatchTraceResponse } from '$lib/settings/types';

describe('DLNA decision trace helpers', () => {
	it('formats filenames and step rows', () => {
		expect(basename('/library/Movies/Movie.mkv')).toBe('Movie.mkv');
		const profileMatch: DLNAProfileMatchTraceResponse = {
			profileId: 'lg-webos',
			profileName: 'LG webOS',
			sourceProfileId: 'lg-webos',
			matchSource: 'headers',
			matchReason: 'matched user agent',
			winningRule: 'userAgent:LG',
			fallbackPath: 'generic',
			score: 120,
			candidateProfileIds: ['lg-webos'],
			headersSummary: ['User-Agent: LG TV'],
			ruleTrace: [
				{
					profileId: 'lg-webos',
					profileName: 'LG webOS',
					field: 'userAgent',
					value: 'LG TV',
					rule: 'contains LG',
					score: 120,
					result: 'pass'
				}
			]
		};
		const delivery: DLNADeliveryTraceResponse = {
			profileId: 'lg-webos',
			profileName: 'LG webOS',
			mediaFileName: 'Movie.mkv',
			objectId: 'movie-1',
			resourceId: 'resource-1',
			streamMode: 'direct',
			deliveryProtocol: 'file',
			mode: 'direct',
			videoCodec: 'copy',
			audioCodec: 'copy',
			reasonCodes: ['direct-play'],
			capabilityTrace: [
				{
					field: 'videoCodec',
					value: 'h264',
					rule: 'supported',
					result: 'pass'
				},
				{
					field: 'audioCodec',
					value: 'aac',
					rule: 'supported',
					result: 'pass'
				}
			]
		};

		expect(buildTraceSteps(profileMatch, delivery)).toEqual([
			{
				id: 'profile-selected',
				stage: 'Profile matching',
				field: 'selected profile',
				rule: 'matched user agent',
				value: 'LG webOS (lg-webos)',
				score: 120,
				result: 'pass'
			},
			{
				id: 'profile-0',
				stage: 'Profile matching',
				field: 'userAgent',
				rule: 'LG webOS (lg-webos): contains LG',
				value: 'LG TV',
				score: 120,
				result: 'pass'
			},
			{
				id: 'delivery-0',
				stage: 'Delivery decision',
				field: 'videoCodec',
				rule: 'supported',
				value: 'h264',
				result: 'pass'
			},
			{
				id: 'delivery-1',
				stage: 'Delivery decision',
				field: 'audioCodec',
				rule: 'supported',
				value: 'aac',
				result: 'pass'
			}
		]);

		expect(buildTraceSummary(profileMatch, delivery)).toEqual([
			{ label: 'Profile name', value: 'LG webOS' },
			{ label: 'Delivery Protocol', value: 'file' },
			{ label: 'Delivery mode', value: 'direct' },
			{ label: 'Video codec', value: 'copy (h264)' },
			{ label: 'Audio codec', value: 'copy (aac)' },
			{ label: 'Match score', value: '120' },
			{ label: 'Match reason', value: 'matched user agent' },
			{ label: 'Winning rule', value: 'userAgent:LG' }
		]);
	});

	it('shows current codecs for direct delivery when the output codecs are omitted', () => {
		const delivery: DLNADeliveryTraceResponse = {
			profileId: 'generic',
			profileName: 'Generic DLNA',
			mediaFileName: 'Movie.mkv',
			objectId: 'movie-1',
			resourceId: 'resource-1',
			streamMode: 'direct',
			deliveryProtocol: 'file',
			mode: 'direct',
			videoCodec: '',
			audioCodec: '',
			reasonCodes: ['direct-play'],
			capabilityTrace: [
				{ field: 'videoCodec', value: 'h264', rule: 'h264', result: 'pass' },
				{ field: 'audioCodec', value: 'aac', rule: 'aac', result: 'pass' }
			]
		};

		expect(buildTraceSummary(undefined, delivery)).toEqual(
			expect.arrayContaining([
				{ label: 'Video codec', value: 'copy (h264)' },
				{ label: 'Audio codec', value: 'copy (aac)' }
			])
		);
	});

	it('formats remux and transcoded details from the capability trace', () => {
		const delivery: DLNADeliveryTraceResponse = {
			profileId: 'generic',
			profileName: 'Generic DLNA',
			mediaFileName: 'Movie.mkv',
			objectId: 'movie-1',
			resourceId: 'resource-1',
			streamMode: 'remux',
			deliveryProtocol: 'file',
			mode: 'remux',
			videoCodec: 'copy',
			audioCodec: 'aac',
			reasonCodes: ['container_not_supported'],
			capabilityTrace: [
				{ field: 'container', value: 'mkv', rule: 'mp4', result: 'fail' },
				{ field: 'videoCodec', value: 'hevc', rule: 'hevc', result: 'pass' },
				{ field: 'audioCodec', value: 'dts', rule: 'aac', result: 'fail' }
			]
		};

		expect(buildTraceSummary(undefined, delivery)).toEqual(
			expect.arrayContaining([
				{ label: 'Delivery mode', value: 'remux (mkv -> mpegts)' },
				{ label: 'Video codec', value: 'copy (hevc)' },
				{ label: 'Audio codec', value: 'transcoding (dts -> aac)' }
			])
		);
	});

	it('does not show empty codec details', () => {
		const delivery: DLNADeliveryTraceResponse = {
			profileId: 'generic',
			profileName: 'Generic DLNA',
			mediaFileName: 'Movie.mkv',
			objectId: 'movie-1',
			resourceId: 'resource-1',
			streamMode: 'direct',
			deliveryProtocol: 'file',
			mode: 'direct',
			videoCodec: 'copy',
			audioCodec: 'aac',
			reasonCodes: [],
			capabilityTrace: [
				{ field: 'videoCodec', value: '', rule: 'supported', result: 'pass' },
				{ field: 'audioCodec', value: '   ', rule: 'unsupported', result: 'fail' }
			]
		};

		expect(buildTraceSummary(undefined, delivery)).toEqual(
			expect.arrayContaining([
				{ label: 'Video codec', value: 'copy' },
				{ label: 'Audio codec', value: 'transcoding (aac)' }
			])
		);
	});

	it('filters failed trace rows when requested', () => {
		const steps = [
			{
				id: 'a',
				stage: 'Profile matching',
				field: 'userAgent',
				rule: 'contains',
				value: 'LG',
				result: 'pass'
			},
			{
				id: 'b',
				stage: 'Delivery decision',
				field: 'codec',
				rule: 'supported',
				value: 'vp9',
				result: 'fail'
			}
		];

		expect(filterTraceSteps(steps, true)).toEqual([steps[0]]);
		expect(filterTraceSteps(steps, false)).toEqual(steps);
	});
});
