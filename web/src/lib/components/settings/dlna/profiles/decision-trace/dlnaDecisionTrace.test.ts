import { describe, expect, it } from 'vitest';
import {
	basename,
	buildDeliveryTraceSteps,
	buildDeliveryTraceSummary,
	buildProfileMatchView,
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
			matchSource: 'match',
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
			],
			candidates: [
				{
					profileId: 'lg-webos',
					profileName: 'LG webOS',
					score: 120,
					minimumScore: 50,
					priority: 100,
					qualified: true,
					selected: true,
					ruleTrace: []
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

		expect(buildDeliveryTraceSteps(delivery)).toEqual([
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

		expect(buildDeliveryTraceSummary(delivery)).toEqual([
			{ label: 'Profile name', value: 'LG webOS' },
			{ label: 'Delivery Protocol', value: 'file' },
			{ label: 'Delivery mode', value: 'direct' },
			{ label: 'Video codec', value: 'copy (h264)' },
			{ label: 'Audio codec', value: 'copy (aac)' }
		]);

		expect(buildProfileMatchView(profileMatch)).toMatchObject({
			profileId: 'lg-webos',
			selectionMethod: 'Automatic profile match',
			score: 120,
			candidates: [{ profileId: 'lg-webos', selected: true }]
		});
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

		expect(buildDeliveryTraceSummary(delivery)).toEqual(
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

		expect(buildDeliveryTraceSummary(delivery)).toEqual(
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

		expect(buildDeliveryTraceSummary(delivery)).toEqual(
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
