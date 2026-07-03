package metadata

import (
	"strconv"
	"strings"
)

func tmdbCreatorPeople(creators []tmdbName) []Person {
	people := make([]Person, 0, len(creators))
	for _, creator := range creators {
		name := strings.TrimSpace(creator.Name)
		if name == "" {
			continue
		}
		people = append(people, Person{
			ExternalProvider: optionalString("tmdb"),
			ExternalID:       optionalString(strconv.FormatInt(creator.ID, 10)),
			Name:             name,
			Role:             optionalString("Creator"),
			ProfilePath:      optionalString(creator.ProfilePath),
		})
	}
	return people
}

func tmdbCrewPeople(crew []tmdbCrewMember) []Person {
	people := []Person{}
	seen := map[string]bool{}
	for _, member := range crew {
		name := strings.TrimSpace(member.Name)
		roles := tmdbCrewLabels(member)
		if name == "" || len(roles) == 0 {
			continue
		}
		for _, role := range roles {
			key := strconv.FormatInt(member.ID, 10) + ":" + role
			if seen[key] {
				continue
			}
			seen[key] = true
			people = append(people, Person{
				ExternalProvider: optionalString("tmdb"),
				ExternalID:       optionalString(strconv.FormatInt(member.ID, 10)),
				Name:             name,
				Role:             optionalString(role),
				ProfilePath:      optionalString(member.ProfilePath),
			})
			if len(people) >= 80 {
				return people
			}
		}
	}
	return people
}

func tmdbCrewFacts(crew []tmdbCrewMember) []Fact {
	mapped := map[string][]string{
		"Director": {},
		"Writer":   {},
		"Editor":   {},
		"Producer": {},
	}
	for _, member := range crew {
		name := strings.TrimSpace(member.Name)
		if name == "" {
			continue
		}
		for _, label := range tmdbCrewLabels(member) {
			mapped[label] = appendUniqueLimit(mapped[label], name, 12)
		}
	}
	facts := []Fact{}
	for _, label := range []string{"Director", "Writer", "Editor", "Producer"} {
		if len(mapped[label]) > 0 {
			facts = append(facts, Fact{Label: label, Value: strings.Join(mapped[label], ", ")})
		}
	}
	return facts
}

func tmdbCrewLabels(member tmdbCrewMember) []string {
	switch member.Job {
	case "Director":
		return []string{"Director"}
	case "Writer", "Screenplay", "Story", "Teleplay":
		return []string{"Writer"}
	case "Editor":
		return []string{"Editor"}
	case "Producer", "Executive Producer":
		return []string{"Producer"}
	}
	switch member.Department {
	case "Writing":
		return []string{"Writer"}
	case "Editing":
		return []string{"Editor"}
	case "Production":
		return []string{"Producer"}
	default:
		return nil
	}
}
