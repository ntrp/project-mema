#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
MEDIA_DIR="$ROOT_DIR/.data/media/test-movie"
DURATION="3"

need() {
	if ! command -v "$1" >/dev/null 2>&1; then
		echo "Missing required tool: $1" >&2
		exit 1
	fi
}

clean_name() {
	printf '%s' "$1" | sed 's/\.mkv$//'
}

bitrate_bps() {
	case "$1" in
		*k) printf '%s000' "${1%k}" ;;
		*) printf '%s' "$1" ;;
	esac
}

write_srt() {
	local path="$1"
	local label="$2"
	cat >"$path" <<SRT
1
00:00:00,000 --> 00:00:01,000
Mock subtitle: $label

2
00:00:02,000 --> 00:00:03,000
Mock subtitle second line: $label
SRT
}

write_chapters() {
	local path="$1"
	cat >"$path" <<'META'
;FFMETADATA1
[CHAPTER]
TIMEBASE=1/1000
START=0
END=1000
title=Opening
[CHAPTER]
TIMEBASE=1/1000
START=1000
END=2000
title=Middle
[CHAPTER]
TIMEBASE=1/1000
START=2000
END=3000
title=Credits
META
}

poster() {
	local path="$1"
	ffmpeg -hide_banner -loglevel error -y \
		-f lavfi -i "color=c=white:s=300x450:d=0.1" \
		-frames:v 1 "$path"
}

movie_one_audio() {
	local out="$1" audio_lang="$2" audio_codec="$3" audio_bitrate="$4" channels="$5" sub_lang="${6:-}" chapters="${7:-no}"
	local dir base srt="" meta="" sub_index="" chapter_index=""
	dir="$(dirname "$out")"
	base="$(clean_name "$(basename "$out")")"
	mkdir -p "$dir"
	local args=(-hide_banner -loglevel error -y
		-f lavfi -i "color=c=white:s=320x180:r=24:d=$DURATION"
		-f lavfi -i "anoisesrc=color=white:duration=$DURATION:sample_rate=48000:amplitude=0.02")
	if [[ -n "$sub_lang" ]]; then
		srt="$dir/$base.embedded.$sub_lang.srt"
		write_srt "$srt" "$base embedded $sub_lang"
		args+=(-i "$srt")
		sub_index="2"
	fi
	if [[ "$chapters" == "yes" ]]; then
		meta="$dir/$base.chapters.ffmetadata"
		write_chapters "$meta"
		args+=(-i "$meta")
		chapter_index="$([[ -n "$sub_index" ]] && echo 3 || echo 2)"
	fi
	args+=(-map 0:v:0 -map 1:a:0)
	if [[ -n "$sub_index" ]]; then args+=(-map "$sub_index:s:0"); fi
	if [[ -n "$chapter_index" ]]; then args+=(-map_metadata "$chapter_index" -map_chapters "$chapter_index"); fi
	args+=(-t "$DURATION" -c:v libx264 -preset ultrafast -pix_fmt yuv420p
		-c:a "$audio_codec" -b:a "$audio_bitrate" -ac "$channels")
	if [[ -n "$sub_index" ]]; then args+=(-c:s srt); fi
	args+=(-metadata:s:v:0 "title=Mock white video")
	args+=(-metadata:s:a:0 "language=$audio_lang")
	args+=(-metadata:s:a:0 "title=Mock white noise $audio_lang")
	args+=(-metadata:s:a:0 "BPS=$(bitrate_bps "$audio_bitrate")")
	if [[ -n "$sub_index" ]]; then
		args+=(-metadata:s:s:0 "language=$sub_lang")
		args+=(-metadata:s:s:0 "title=Mock subtitle $sub_lang")
	fi
	args+=("$out")
	ffmpeg "${args[@]}"
	rm -f "$srt" "$meta"
}

movie_two_audio() {
	local out="$1"
	local dir base srt
	dir="$(dirname "$out")"
	base="$(clean_name "$(basename "$out")")"
	mkdir -p "$dir"
	srt="$dir/$base.embedded.eng.srt"
	write_srt "$srt" "$base embedded eng"
	ffmpeg -hide_banner -loglevel error -y \
		-f lavfi -i "color=c=white:s=320x180:r=24:d=$DURATION" \
		-f lavfi -i "anoisesrc=color=white:duration=$DURATION:sample_rate=48000:amplitude=0.02" \
		-f lavfi -i "anoisesrc=color=white:duration=$DURATION:sample_rate=48000:amplitude=0.02" \
		-i "$srt" \
		-map 0:v:0 -map 1:a:0 -map 2:a:0 -map 3:s:0 \
		-t "$DURATION" -c:v libx264 -preset ultrafast -pix_fmt yuv420p \
		-c:a:0 aac -b:a:0 256k -ac:a:0 2 \
		-c:a:1 aac -b:a:1 192k -ac:a:1 2 \
		-c:s srt \
		-metadata:s:v:0 "title=Mock white video" \
		-metadata:s:a:0 language=eng \
		-metadata:s:a:0 "title=Mock white noise eng" \
		-metadata:s:a:0 BPS=256000 \
		-metadata:s:a:1 language=spa \
		-metadata:s:a:1 "title=Mock white noise spa" \
		-metadata:s:a:1 BPS=192000 \
		-metadata:s:s:0 language=eng \
		-metadata:s:s:0 "title=Mock subtitle eng" \
		"$out"
	rm -f "$srt"
}

verify_media_file() {
	local file="$1"
	local streams
	streams="$(ffprobe -v error -show_entries stream=codec_type -of csv=p=0 "$file")"
	if ! grep -qx 'video' <<<"$streams"; then
		echo "Generated media is missing a video stream: $file" >&2
		exit 1
	fi
	if ! grep -qx 'audio' <<<"$streams"; then
		echo "Generated media is missing an audio stream: $file" >&2
		exit 1
	fi
	streams="$(ffprobe -v error -show_streams -show_format -of json "$(cd "$(dirname "$file")" && pwd)/$(basename "$file")")"
	if ! grep -q '"codec_type": "video"' <<<"$streams"; then
		echo "App-style media probe is missing a video stream: $file" >&2
		exit 1
	fi
	if ! grep -q '"codec_type": "audio"' <<<"$streams"; then
		echo "App-style media probe is missing an audio stream: $file" >&2
		exit 1
	fi
}

verify_media_tree() {
	local found=0
	while IFS= read -r -d '' file; do
		found=1
		verify_media_file "$file"
	done < <(find "$MEDIA_DIR" -type f -name '*.mkv' -print0)
	if [[ "$found" -eq 0 ]]; then
		echo "No MKV files were generated under $MEDIA_DIR" >&2
		exit 1
	fi
}

external_srt() {
	local movie="$1" lang="$2"
	local dir base
	dir="$(dirname "$movie")"
	base="$(clean_name "$(basename "$movie")")"
	write_srt "$dir/$base.$lang.srt" "$base external $lang"
}

mock_other_files() {
	local movie="$1"
	local dir base
	dir="$(dirname "$movie")"
	base="$(clean_name "$(basename "$movie")")"
	poster "$dir/$base.poster.jpg"
	printf 'Mock metadata for %s\n' "$base" >"$dir/$base.nfo"
}

write_profile_refs() {
	local key="$1"
	shift
	printf '%s\n' "$@" >"$MEDIA_DIR/$key/MEMA_PROFILES.txt"
}

scenario() {
	local key="$1" movie="$2" profile="${3:-$1}"
	mkdir -p "$MEDIA_DIR/$key"
	printf '%s\n' "$movie" >"$MEDIA_DIR/$key/TMDB_MOVIE.txt"
	write_profile_refs "$key" "$profile"
}

scenario_movies() {
	local key="$1"
	shift
	mkdir -p "$MEDIA_DIR/$key"
	printf '%s\n' "$@" >"$MEDIA_DIR/$key/TMDB_MOVIES.txt"
}

verify_profile_refs() {
	local dir found=0
	while IFS= read -r -d '' dir; do
		found=1
		if [[ ! -s "$dir/MEMA_PROFILES.txt" ]]; then
			echo "Generated fixture is missing MEMA_PROFILES.txt: $dir" >&2
			exit 1
		fi
	done < <(find "$MEDIA_DIR" -mindepth 1 -maxdepth 1 -type d -print0)
	if [[ "$found" -eq 0 ]]; then
		echo "No fixture directories were generated under $MEDIA_DIR" >&2
		exit 1
	fi
}

need ffmpeg
need ffprobe
rm -rf "$MEDIA_DIR"
mkdir -p "$MEDIA_DIR"

scenario "01-ok-embedded" "Finding Nemo (2003), TMDB 12"
movie_one_audio "$MEDIA_DIR/01-ok-embedded/Finding.Nemo.2003.tmdb-12.1080p.WEB-DL.AAC2.0.EN.mkv" eng aac 256k 2 eng yes
mock_other_files "$MEDIA_DIR/01-ok-embedded/Finding.Nemo.2003.tmdb-12.1080p.WEB-DL.AAC2.0.EN.mkv"

scenario "02-missing-italian-audio" "Amelie (2001), TMDB 194"
movie_one_audio "$MEDIA_DIR/02-missing-italian-audio/Amelie.2001.tmdb-194.1080p.WEB-DL.AAC2.0.EN.mkv" eng aac 256k 2 eng no
external_srt "$MEDIA_DIR/02-missing-italian-audio/Amelie.2001.tmdb-194.1080p.WEB-DL.AAC2.0.EN.mkv" ita

scenario "03-wrong-audio-codec" "The Matrix (1999), TMDB 603"
movie_one_audio "$MEDIA_DIR/03-wrong-audio-codec/The.Matrix.1999.tmdb-603.1080p.WEB-DL.AC3.5.1.EN.mkv" eng ac3 640k 6 eng no

scenario "04-wrong-audio-channels" "Spirited Away (2001), TMDB 129"
movie_one_audio "$MEDIA_DIR/04-wrong-audio-channels/Spirited.Away.2001.tmdb-129.1080p.WEB-DL.AAC1.0.JA.mkv" jpn aac 192k 1 eng no

scenario "05-low-audio-bitrate" "Mad Max: Fury Road (2015), TMDB 76341"
movie_one_audio "$MEDIA_DIR/05-low-audio-bitrate/Mad.Max.Fury.Road.2015.tmdb-76341.1080p.WEB-DL.AAC2.0.EN.mkv" eng aac 96k 2 eng no

scenario "06-unwanted-audio" "Ratatouille (2007), TMDB 2062"
movie_two_audio "$MEDIA_DIR/06-unwanted-audio/Ratatouille.2007.tmdb-2062.1080p.WEB-DL.AAC2.0.EN-ES.mkv"

scenario "07-embedded-subtitle-needed" "Paddington (2014), TMDB 116149"
movie_one_audio "$MEDIA_DIR/07-embedded-subtitle-needed/Paddington.2014.tmdb-116149.1080p.WEB-DL.AAC2.0.EN.mkv" eng aac 256k 2 "" no
external_srt "$MEDIA_DIR/07-embedded-subtitle-needed/Paddington.2014.tmdb-116149.1080p.WEB-DL.AAC2.0.EN.mkv" eng

scenario "08-external-subtitle-mode" "Arrival (2016), TMDB 329865"
movie_one_audio "$MEDIA_DIR/08-external-subtitle-mode/Arrival.2016.tmdb-329865.1080p.WEB-DL.AAC2.0.EN.mkv" eng aac 256k 2 "" no
external_srt "$MEDIA_DIR/08-external-subtitle-mode/Arrival.2016.tmdb-329865.1080p.WEB-DL.AAC2.0.EN.mkv" eng

scenario "09-mixed-existing-external" "Inside Out (2015), TMDB 150540"
movie_one_audio "$MEDIA_DIR/09-mixed-existing-external/Inside.Out.2015.tmdb-150540.1080p.WEB-DL.AAC2.0.EN.mkv" eng aac 256k 2 eng no
external_srt "$MEDIA_DIR/09-mixed-existing-external/Inside.Out.2015.tmdb-150540.1080p.WEB-DL.AAC2.0.EN.mkv" ita

scenario "10-unwanted-subtitle" "Parasite (2019), TMDB 496243"
movie_one_audio "$MEDIA_DIR/10-unwanted-subtitle/Parasite.2019.tmdb-496243.1080p.WEB-DL.AAC2.0.KO.mkv" kor aac 256k 2 spa no
external_srt "$MEDIA_DIR/10-unwanted-subtitle/Parasite.2019.tmdb-496243.1080p.WEB-DL.AAC2.0.KO.mkv" eng

scenario "11-chapter-delete-summary" "Inception (2010), TMDB 27205"
movie_one_audio "$MEDIA_DIR/11-chapter-delete-summary/Inception.2010.tmdb-27205.1080p.WEB-DL.AAC2.0.EN.Chapters.mkv" eng aac 256k 2 eng yes

scenario "12-other-files-actions" "WALL-E (2008), TMDB 10681"
movie_one_audio "$MEDIA_DIR/12-other-files-actions/WALL-E.2008.tmdb-10681.1080p.WEB-DL.AAC2.0.EN.mkv" eng aac 256k 2 "" no
external_srt "$MEDIA_DIR/12-other-files-actions/WALL-E.2008.tmdb-10681.1080p.WEB-DL.AAC2.0.EN.mkv" eng
mock_other_files "$MEDIA_DIR/12-other-files-actions/WALL-E.2008.tmdb-10681.1080p.WEB-DL.AAC2.0.EN.mkv"

scenario_movies "13-three-movies-one-folder" \
	"The Grand Budapest Hotel (2014), TMDB 120467" \
	"Interstellar (2014), TMDB 157336" \
	"The Truman Show (1998), TMDB 37165"
write_profile_refs "13-three-movies-one-folder" \
	"13-three-movies-one-folder-grand-budapest" \
	"13-three-movies-one-folder-interstellar" \
	"13-three-movies-one-folder-truman-show"
movie_one_audio "$MEDIA_DIR/13-three-movies-one-folder/The.Grand.Budapest.Hotel.2014.tmdb-120467.1080p.WEB-DL.AAC2.0.EN.mkv" eng aac 256k 2 eng no
movie_one_audio "$MEDIA_DIR/13-three-movies-one-folder/Interstellar.2014.tmdb-157336.1080p.WEB-DL.AAC2.0.EN.mkv" eng aac 256k 2 eng no
movie_one_audio "$MEDIA_DIR/13-three-movies-one-folder/The.Truman.Show.1998.tmdb-37165.1080p.WEB-DL.AAC2.0.EN.mkv" eng aac 256k 2 eng no

cat >"$MEDIA_DIR/README.md" <<'README'
# Test Movie Fixtures

Generated by `scripts/restore-test-media.sh`.

Each fixture folder has `MEMA_PROFILES.txt` with the dev-seeded profile/media names.
The dev seed pre-imports those media items with the matching profile attached and monitoring off.

Use the seeded case media/profile combinations to exercise audio status, subtitle modes,
external subtitle import, unwanted subtitle/audio marking, chapter actions, other-file actions,
and multiple movies in one folder.
All videos are white test clips, all audio is generated white noise, and all subtitles/posters are mock data.
README

verify_media_tree
verify_profile_refs
echo "Restored test media under $MEDIA_DIR"
