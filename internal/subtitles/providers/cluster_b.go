package providers

import (
	"media-manager/internal/subtitles/providers/animekalesi"
	"media-manager/internal/subtitles/providers/animesubinfo"
	"media-manager/internal/subtitles/providers/bayflix"
	"media-manager/internal/subtitles/providers/greeksubs"
	"media-manager/internal/subtitles/providers/greeksubtitles"
	"media-manager/internal/subtitles/providers/hosszupuska"
	"media-manager/internal/subtitles/providers/nekur"
	"media-manager/internal/subtitles/providers/prijevodionline"
	"media-manager/internal/subtitles/providers/soustitreseu"
	"media-manager/internal/subtitles/providers/subclub"
	"media-manager/internal/subtitles/providers/subf2m"
	"media-manager/internal/subtitles/providers/subs4free"
	"media-manager/internal/subtitles/providers/subssabbz"
	"media-manager/internal/subtitles/providers/subsunacs"
	"media-manager/internal/subtitles/providers/subsynchro"
	"media-manager/internal/subtitles/providers/subtitrarinoi"
	"media-manager/internal/subtitles/providers/subtitriid"
	"media-manager/internal/subtitles/providers/subtitulamostv"
	"media-manager/internal/subtitles/providers/supersubtitles"
	"media-manager/internal/subtitles/providers/titrari"
	"media-manager/internal/subtitles/providers/tvsubtitles"
	"media-manager/internal/subtitles/providers/vladoonmooo"
	"media-manager/internal/subtitles/providers/wizdom"
	"media-manager/internal/subtitles/providers/yavkanet"
	"media-manager/internal/subtitles/providers/yifysubtitles"
)

func init() {
	Register("animekalesi", animekalesi.Adapter)
	Register("animesubinfo", animesubinfo.Adapter)
	Register("bayflix", bayflix.Adapter)
	Register("greeksubs", greeksubs.Adapter)
	Register("greeksubtitles", greeksubtitles.Adapter)
	Register("hosszupuska", hosszupuska.Adapter)
	Register("nekur", nekur.Adapter)
	Register("prijevodionline", prijevodionline.Adapter)
	Register("soustitreseu", soustitreseu.Adapter)
	Register("subclub", subclub.Adapter)
	Register("subf2m", subf2m.Adapter)
	Register("subs4free", subs4free.Adapter)
	Register("subssabbz", subssabbz.Adapter)
	Register("subsunacs", subsunacs.Adapter)
	Register("subsynchro", subsynchro.Adapter)
	Register("subtitrarinoi", subtitrarinoi.Adapter)
	Register("subtitriid", subtitriid.Adapter)
	Register("subtitulamostv", subtitulamostv.Adapter)
	Register("supersubtitles", supersubtitles.Adapter)
	Register("titrari", titrari.Adapter)
	Register("tvsubtitles", tvsubtitles.Adapter)
	Register("vladoonmooo", vladoonmooo.Adapter)
	Register("wizdom", wizdom.Adapter)
	Register("yavkanet", yavkanet.Adapter)
	Register("yifysubtitles", yifysubtitles.Adapter)
}
