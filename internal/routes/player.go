package routes

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"
	"sync"

	"github.com/shakerrdev/cs-stalker/internal/api"
)

const (
	recentMatchesShown = 12
	teammatesShown     = 10
)

type playerPage struct {
	Summary      api.Summary
	Level        int
	BanState     string
	BanLines     []string
	Playtime     *api.Playtime
	PlaytimeNote string
	Leetify      *api.Leetify
	LeetifyNote  string
	Teammates    []teammate
	MapRanks     []api.LeetifyMapRank
	Matches      []api.LeetifyMatch
	OlderMatches []api.LeetifyMatch
}

type teammate struct {
	Steam64ID string
	Name      string
	Avatar    string
	Matches   int
}

func gatherPlayer(summary api.Summary) playerPage {
	steam64id := summary.Steam64id
	page := playerPage{Summary: summary}
	var bans *api.Bans

	var wg sync.WaitGroup
	wg.Add(4)
	go func() {
		defer wg.Done()
		level, err := api.SteamLevel(steam64id)
		switch {
		case err == nil:
			page.Level = level
		case !errors.Is(err, api.ErrHidden):
			log.Printf("player %s: %v", steam64id, err)
		}
	}()
	go func() {
		defer wg.Done()
		found, err := api.PlayerBans(steam64id)
		if err != nil {
			log.Printf("player %s: %v", steam64id, err)
			return
		}
		bans = &found
	}()
	go func() {
		defer wg.Done()
		playtime, err := api.CSPlaytime(steam64id)
		switch {
		case err == nil:
			page.Playtime = &playtime
		case errors.Is(err, api.ErrHidden):
			page.PlaytimeNote = "Their Steam game details are private."
		case errors.Is(err, api.ErrNotFound):
			page.PlaytimeNote = "No Counter-Strike playtime on this account."
		default:
			log.Printf("player %s: %v", steam64id, err)
			page.PlaytimeNote = "Couldn't load playtime right now."
		}
	}()
	go func() {
		defer wg.Done()
		profile, err := api.LeetifyProfile(steam64id)
		switch {
		case err == nil:
			page.Leetify = &profile
			page.Teammates = lookUpTeammates(profile.RecentTeammates)
		case errors.Is(err, api.ErrNotFound):
			page.LeetifyNote = "Leetify has no data on this player, so ranks, ratings and matches aren't available."
		default:
			log.Printf("player %s: %v", steam64id, err)
			page.LeetifyNote = "Couldn't load Leetify data right now. Reload the page to try again."
		}
	}()
	wg.Wait()

	var leetifyBans []api.LeetifyBan
	if page.Leetify != nil {
		leetifyBans = page.Leetify.Bans
		page.MapRanks = rankedMaps(page.Leetify.Ranks.Competitive)
		matches := page.Leetify.RecentMatches
		page.Matches = matches[:min(len(matches), recentMatchesShown)]
		page.OlderMatches = matches[len(page.Matches):]
	}
	page.BanState, page.BanLines = describeBans(bans, leetifyBans)
	return page
}

func lookUpTeammates(recent []api.LeetifyTeammate) []teammate {
	sorted := append([]api.LeetifyTeammate(nil), recent...)
	sort.SliceStable(sorted, func(i, j int) bool { return sorted[i].MatchesTogether > sorted[j].MatchesTogether })
	sorted = sorted[:min(len(sorted), teammatesShown)]

	ids := make([]string, len(sorted))
	for i, mate := range sorted {
		ids[i] = mate.Steam64ID
	}
	summaries, err := api.PlayerSummaries(ids)
	if err != nil {
		log.Printf("teammates: %v", err)
	}

	mates := make([]teammate, len(sorted))
	for i, mate := range sorted {
		summary := summaries[mate.Steam64ID]
		mates[i] = teammate{Steam64ID: mate.Steam64ID, Name: summary.Name, Avatar: summary.AvatarMedium, Matches: mate.MatchesTogether}
	}
	return mates
}

func rankedMaps(ranks []api.LeetifyMapRank) []api.LeetifyMapRank {
	var ranked []api.LeetifyMapRank
	for _, rank := range ranks {
		if rank.Rank > 0 {
			ranked = append(ranked, rank)
		}
	}
	sort.SliceStable(ranked, func(i, j int) bool {
		if ranked[i].Rank != ranked[j].Rank {
			return ranked[i].Rank > ranked[j].Rank
		}
		return mapName(ranked[i].Map) < mapName(ranked[j].Map)
	})
	return ranked
}

func describeBans(bans *api.Bans, leetifyBans []api.LeetifyBan) (string, []string) {
	var lines []string
	if bans != nil {
		if bans.VACBans > 0 {
			lines = append(lines, count(bans.VACBans, "VAC ban"))
		}
		if bans.GameBans > 0 {
			lines = append(lines, count(bans.GameBans, "game ban"))
		}
		if bans.VACBans > 0 || bans.GameBans > 0 {
			lines = append(lines, lastBan(bans.DaysSinceLastBan))
		}
		if bans.CommunityBanned {
			lines = append(lines, "Banned from the Steam Community")
		}
		switch bans.TradeBan {
		case "":
		case "probation":
			lines = append(lines, "On trade probation")
		default:
			lines = append(lines, "Trade banned")
		}
	}
	for _, ban := range leetifyBans {
		line := fmt.Sprintf("%s ban since %s", platformName(ban.Platform), monthYear(ban.BannedSince))
		if ban.Nickname != "" {
			line += " as " + ban.Nickname
		}
		lines = append(lines, line)
	}

	switch {
	case len(lines) > 0:
		return "banned", lines
	case bans == nil:
		return "unknown", nil
	}
	return "clean", nil
}

func count(n int, noun string) string {
	if n == 1 {
		return "1 " + noun
	}
	return fmt.Sprintf("%d %ss", n, noun)
}

func lastBan(days int) string {
	switch days {
	case 0:
		return "Last ban today"
	case 1:
		return "Last ban yesterday"
	}
	return fmt.Sprintf("Last ban %s days ago", thousands(days))
}

func platformName(platform string) string {
	switch strings.ToLower(platform) {
	case "faceit":
		return "FACEIT"
	case "esea":
		return "ESEA"
	case "":
		return "Platform"
	}
	return strings.ToUpper(platform[:1]) + platform[1:]
}
