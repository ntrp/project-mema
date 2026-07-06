<script lang="ts">
	import VideoIcon from '@lucide/svelte/icons/video';
	import MediaProfileQualitySelector from '$lib/components/settings/profiles/MediaProfileQualitySelector.svelte';
	import ProfileTargetMultiSelect from '$lib/components/settings/profiles/ProfileTargetMultiSelect.svelte';
	import * as Card from '$lib/components/ui/card';
	import type { MediaProfileForm, QualitySizeSetting } from '$lib/settings/types';
	import { hdrFormatOptions, pixelFormatOptions, videoCodecOptions } from './profileTargetOptions';

	interface Props {
		form: MediaProfileForm;
		qualities: QualitySizeSetting[];
		loadingQualities: boolean;
		qualityError: string;
		onChange: (_form: MediaProfileForm) => void;
	}

	let { form, qualities, loadingQualities, qualityError, onChange }: Props = $props();

	function patch(patchValue: Partial<MediaProfileForm>) {
		onChange({ ...form, ...patchValue });
	}

	function patchVideo(patchValue: Partial<MediaProfileForm['videoTarget']>) {
		patch({ videoTarget: { ...form.videoTarget, ...patchValue } });
	}
</script>

<Card.Root>
	<Card.Header>
		<Card.Title class="flex items-center gap-2">
			<VideoIcon aria-hidden="true" />
			<span>Video Target</span>
		</Card.Title>
	</Card.Header>
	<Card.Content class="grid gap-3 mt-3">
		<div class="grid gap-3 rounded-md bg-muted/30 p-3 text-sm md:grid-cols-3">
			<ProfileTargetMultiSelect
				id="video-target-codecs"
				label="Video codec"
				values={form.videoTarget.codecs ?? []}
				options={videoCodecOptions}
				placeholder="Any codec"
				onChange={(values) => patchVideo({ codecs: values })}
			/>
			<ProfileTargetMultiSelect
				id="video-target-hdr"
				label="HDR"
				values={form.videoTarget.hdrFormats ?? []}
				options={hdrFormatOptions}
				placeholder="Any HDR format"
				onChange={(values) => patchVideo({ hdrFormats: values })}
			/>
			<ProfileTargetMultiSelect
				id="video-target-pixel-formats"
				label="Pixel format"
				values={form.videoTarget.pixelFormats ?? []}
				options={pixelFormatOptions}
				placeholder="Any pixel format"
				onChange={(values) => patchVideo({ pixelFormats: values })}
			/>
		</div>

		<details class="group rounded-md border border-border bg-background">
			<summary
				class="flex cursor-pointer list-none items-center justify-between gap-3 px-3 py-2.5 text-sm font-bold text-muted-foreground [&::-webkit-details-marker]:hidden"
			>
				<span>Qualities</span>
				<span>{form.qualityIds.length} selected</span>
			</summary>
			<div class="grid gap-3 border-t border-border p-3">
				<MediaProfileQualitySelector
					{form}
					{qualities}
					loading={loadingQualities}
					error={qualityError}
					showHeader={false}
					onChange={(value) => onChange(value)}
				/>
			</div>
		</details>
	</Card.Content>
</Card.Root>
