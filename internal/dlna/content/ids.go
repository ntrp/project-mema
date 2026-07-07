package content

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"

	"github.com/google/uuid"
)

const RootID = "0"

const idPrefix = "cd1_"

type Ref struct {
	Kind string `json:"k"`
	Key  string `json:"v,omitempty"`
	Aux  string `json:"a,omitempty"`
}

func RootContainerRef(key string) Ref {
	return Ref{Kind: "root", Key: key}
}

func MediaItemRef(id uuid.UUID) Ref {
	return Ref{Kind: "media", Key: id.String()}
}

func SeasonRef(id uuid.UUID) Ref {
	return Ref{Kind: "season", Key: id.String()}
}

func EpisodeRef(id uuid.UUID) Ref {
	return Ref{Kind: "episode", Key: id.String()}
}

func FileRef(mediaID uuid.UUID, path string) Ref {
	return Ref{Kind: "file", Key: mediaID.String(), Aux: filePathHash(path)}
}

func GroupRef(kind string, value string) Ref {
	return Ref{Kind: kind, Key: value}
}

func EncodeID(ref Ref) string {
	payload, _ := json.Marshal(ref)
	return idPrefix + base64.RawURLEncoding.EncodeToString(payload)
}

func DecodeID(id string) (Ref, error) {
	if id == RootID {
		return Ref{Kind: "root"}, nil
	}
	if !strings.HasPrefix(id, idPrefix) {
		return Ref{}, errors.New("invalid content directory object id")
	}
	payload, err := base64.RawURLEncoding.DecodeString(strings.TrimPrefix(id, idPrefix))
	if err != nil {
		return Ref{}, errors.New("invalid content directory object id")
	}
	var ref Ref
	if err := json.Unmarshal(payload, &ref); err != nil {
		return Ref{}, errors.New("invalid content directory object id")
	}
	if strings.TrimSpace(ref.Kind) == "" {
		return Ref{}, errors.New("invalid content directory object id")
	}
	return ref, nil
}

func filePathHash(path string) string {
	sum := sha256.Sum256([]byte(path))
	return hex.EncodeToString(sum[:16])
}
