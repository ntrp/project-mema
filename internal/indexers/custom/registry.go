package custom

import (
	"strings"

	"media-manager/internal/indexers/custom/alpharatio"
	"media-manager/internal/indexers/custom/anidex"
	"media-manager/internal/indexers/custom/animez"
	"media-manager/internal/indexers/custom/avistaz"
	"media-manager/internal/indexers/custom/beyondhd"
	"media-manager/internal/indexers/custom/binsearch"
	"media-manager/internal/indexers/custom/broadcasthenet"
	"media-manager/internal/indexers/custom/brokenstones"
	"media-manager/internal/indexers/custom/cgpeers"
	"media-manager/internal/indexers/custom/cinemaz"
	"media-manager/internal/indexers/custom/dicmusic"
	"media-manager/internal/indexers/custom/exoticaz"
	"media-manager/internal/indexers/custom/filelist"
	"media-manager/internal/indexers/custom/gazelle"
	"media-manager/internal/indexers/custom/greatposterwall"
	"media-manager/internal/indexers/custom/hdbits"
	"media-manager/internal/indexers/custom/headphones"
	"media-manager/internal/indexers/custom/knaben"
	"media-manager/internal/indexers/custom/mteamtp"
	"media-manager/internal/indexers/custom/nebulance"
	"media-manager/internal/indexers/custom/newznab"
	"media-manager/internal/indexers/custom/nzbindex"
	"media-manager/internal/indexers/custom/orpheus"
	"media-manager/internal/indexers/custom/passthepopcorn"
	"media-manager/internal/indexers/custom/privatehd"
	"media-manager/internal/indexers/custom/redacted"
	"media-manager/internal/indexers/custom/retroflix"
	"media-manager/internal/indexers/custom/scenehd"
	"media-manager/internal/indexers/custom/secretcinema"
	"media-manager/internal/indexers/custom/speedapp"
	"media-manager/internal/indexers/custom/subsplease"
	"media-manager/internal/indexers/custom/torrentpotato"
	"media-manager/internal/indexers/custom/torrentrss"
	"media-manager/internal/indexers/custom/torrentscsv"
	"media-manager/internal/indexers/custom/torrentsyndikat"
	"media-manager/internal/indexers/custom/torznab"
	"media-manager/internal/indexers/custom/unit3d"
	"media-manager/internal/indexers/custom/xthor"
	"media-manager/internal/indexers/engine"
)

type Registry struct {
	engines map[string]engine.Engine
}

func NewRegistry(client engine.HTTPDoer) *Registry {
	registry := &Registry{engines: map[string]engine.Engine{}}
	registry.Register("Newznab", newznab.New(client))
	registry.Register("Torznab", torznab.New(client))
	registerProwlarrCustomEngines(registry, client)
	registry.Register("AlphaRatio", alpharatio.New(client))
	registry.Register("AnimeZ", animez.New(client))
	registry.Register("Anidex", anidex.New(client))
	registry.Register("AvistaZ", avistaz.New(client))
	registry.Register("BeyondHD", beyondhd.New(client))
	registry.Register("BinSearch", binsearch.New(client))
	registry.Register("BroadcastheNet", broadcasthenet.New(client))
	registry.Register("BrokenStones", brokenstones.New(client))
	registry.Register("CGPeers", cgpeers.New(client))
	registry.Register("CinemaZ", cinemaz.New(client))
	registry.Register("DICMusic", dicmusic.New(client))
	registry.Register("ExoticaZ", exoticaz.New(client))
	registry.Register("FileList", filelist.New(client))
	registry.Register("Gazelle", gazelle.New(client))
	registry.Register("GreatPosterWall", greatposterwall.New(client))
	registry.Register("HDBits", hdbits.New(client))
	registry.Register("Headphones", headphones.New(client))
	registry.Register("Knaben", knaben.New(client))
	registry.Register("MTeamTp", mteamtp.New(client))
	registry.Register("Nebulance", nebulance.New(client))
	registry.Register("NzbIndex", nzbindex.New(client))
	registry.Register("Orpheus", orpheus.New(client))
	registry.Register("PassThePopcorn", passthepopcorn.New(client))
	registry.Register("PrivateHD", privatehd.New(client))
	registry.Register("Redacted", redacted.New(client))
	registry.Register("RetroFlix", retroflix.New(client))
	registry.Register("SceneHD", scenehd.New(client))
	registry.Register("SecretCinema", secretcinema.New(client))
	registry.Register("SpeedApp", speedapp.New(client))
	registry.Register("SubsPlease", subsplease.New(client))
	registry.Register("TorrentPotato", torrentpotato.New(client))
	registry.Register("TorrentRss", torrentrss.New(client))
	registry.Register("TorrentSyndikat", torrentsyndikat.New(client))
	registry.Register("TorrentsCSV", torrentscsv.New(client))
	registry.Register("UNIT3D", unit3d.New(client))
	registry.Register("Xthor", xthor.New(client))
	return registry
}

func (r *Registry) Register(name string, indexer engine.Engine) {
	if name == "" || indexer == nil {
		return
	}
	r.engines[engineKey(name)] = indexer
}

func (r *Registry) EngineFor(config engine.Config) (engine.Engine, bool) {
	switch strings.ToLower(strings.TrimSpace(config.DefinitionID)) {
	case "generic-newznab":
		return r.engine("Newznab")
	case "generic-torznab":
		return r.engine("Torznab")
	}
	if !strings.EqualFold(config.Implementation, "Cardigann") {
		if indexer, ok := r.engine(config.Implementation); ok {
			return indexer, true
		}
	}
	switch strings.ToLower(strings.TrimSpace(config.Protocol)) {
	case "usenet":
		return r.engine("Newznab")
	case "torrent":
		return r.engine("Torznab")
	default:
		return nil, false
	}
}

func (r *Registry) engine(name string) (engine.Engine, bool) {
	indexer, ok := r.engines[engineKey(name)]
	return indexer, ok
}

func engineKey(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	var out strings.Builder
	for _, r := range value {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
			out.WriteRune(r)
		}
	}
	return out.String()
}
