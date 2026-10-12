package routes

import (
	"fmt"
	"html/template"
	"math"
	"strconv"
	"strings"
	"time"
	"unicode"
)

var now = time.Now

var funcs = template.FuncMap{
	"thousands":      thousands,
	"premier":        premier,
	"signed":         signed,
	"rating100":      rating100,
	"points100":      points100,
	"percent":        percent,
	"fraction":       fraction,
	"decimal":        decimal,
	"hours":          hours,
	"date":           shortDate,
	"dayMonth":       dayMonth,
	"monthYear":      monthYear,
	"ago":            ago,
	"mapName":        mapName,
	"skillGroup":     skillGroupShort,
	"skillGroupName": skillGroupName,
	"faceitTier":     faceitTier,
	"presence":       presence,
	"flag":           flag,
	"source":         source,
	"outcome":        outcome,
	"bar":            bar,
	"deref":          deref,
	"derefFloat":     derefFloat,
}

func thousands(n int) string {
	digits := strconv.Itoa(n)
	sign := ""
	if n < 0 {
		sign, digits = "-", digits[1:]
	}
	var b strings.Builder
	for i, d := range digits {
		if i > 0 && (len(digits)-i)%3 == 0 {
			b.WriteByte(',')
		}
		b.WriteRune(d)
	}
	return sign + b.String()
}

type premierRating struct {
	High, Low string
	Tier      int
}

func premier(rating int) *premierRating {
	if rating <= 0 {
		return nil
	}
	p := &premierRating{Tier: min(rating/5000+1, 7)}
	if rating < 1000 {
		p.Low = strconv.Itoa(rating)
		return p
	}
	p.High = thousands(rating / 1000)
	p.Low = fmt.Sprintf(",%03d", rating%1000)
	return p
}

func deref(v *int) int {
	if v == nil {
		return 0
	}
	return *v
}

func derefFloat(v *float64) float64 {
	if v == nil {
		return 0
	}
	return *v
}

func signed(v float64) string {
	return strings.Replace(fmt.Sprintf("%+.2f", v), "-", "\u2212", 1)
}

func rating100(v float64) string {
	return signed(v * 100)
}

func points100(v float64) string {
	return fmt.Sprintf("%.1f", v*100)
}

func percent(v float64) string {
	return fmt.Sprintf("%.0f%%", v)
}

func fraction(v float64) string {
	return percent(v * 100)
}

func decimal(v float64, places int) string {
	return strconv.FormatFloat(v, 'f', places, 64)
}

func hours(minutes int) string {
	switch {
	case minutes <= 0:
		return "0 h"
	case minutes < 60:
		return "under 1 h"
	}
	return thousands(minutes/60) + " h"
}

func shortDate(t time.Time) string {
	return t.Format("2 Jan 2006")
}

func dayMonth(t time.Time) string {
	if t.Year() == now().Year() {
		return t.Format("2 Jan")
	}
	return t.Format("2 Jan 2006")
}

func monthYear(t time.Time) string {
	return t.Format("Jan 2006")
}

func ago(t time.Time) string {
	days := int(now().Sub(t).Hours() / 24)
	plural := func(n int, unit string) string {
		if n == 1 {
			return "1 " + unit + " ago"
		}
		return fmt.Sprintf("%d %ss ago", n, unit)
	}
	switch {
	case days < 1:
		return "today"
	case days < 2:
		return "yesterday"
	case days < 14:
		return plural(days, "day")
	case days < 60:
		return plural(days/7, "week")
	case days < 730:
		return plural(days/30, "month")
	}
	return plural(days/365, "year")
}

func mapName(name string) string {
	for _, prefix := range []string{"de_", "cs_", "ar_", "dz_"} {
		name = strings.TrimPrefix(name, prefix)
	}
	words := strings.Split(name, "_")
	for i, w := range words {
		if w != "" {
			r := []rune(w)
			r[0] = unicode.ToUpper(r[0])
			words[i] = string(r)
		}
	}
	return strings.Join(words, " ")
}

var skillGroups = []struct{ short, name string }{
	{"S1", "Silver I"}, {"S2", "Silver II"}, {"S3", "Silver III"}, {"S4", "Silver IV"},
	{"SE", "Silver Elite"}, {"SEM", "Silver Elite Master"},
	{"GN1", "Gold Nova I"}, {"GN2", "Gold Nova II"}, {"GN3", "Gold Nova III"}, {"GNM", "Gold Nova Master"},
	{"MG1", "Master Guardian I"}, {"MG2", "Master Guardian II"}, {"MGE", "Master Guardian Elite"},
	{"DMG", "Distinguished Master Guardian"}, {"LE", "Legendary Eagle"}, {"LEM", "Legendary Eagle Master"},
	{"SMFC", "Supreme Master First Class"}, {"GE", "The Global Elite"},
}

func skillGroupShort(rank int) string {
	if rank < 1 || rank > len(skillGroups) {
		return ""
	}
	return skillGroups[rank-1].short
}

func skillGroupName(rank int) string {
	if rank < 1 || rank > len(skillGroups) {
		return ""
	}
	return skillGroups[rank-1].name
}

func faceitTier(level int) int {
	switch {
	case level >= 10:
		return 5
	case level >= 8:
		return 4
	case level >= 4:
		return 3
	case level >= 2:
		return 2
	}
	return 1
}

type presenceState struct {
	Label, Class string
}

func presence(state int, game string) presenceState {
	if game != "" {
		return presenceState{"Playing " + game, "ingame"}
	}
	labels := []string{"Offline", "Online", "Busy", "Away", "Snooze", "Looking to trade", "Looking to play"}
	if state < 1 || state >= len(labels) {
		return presenceState{"Offline", "offline"}
	}
	return presenceState{labels[state], "online"}
}

func flag(countryCode string) string {
	if len(countryCode) != 2 {
		return ""
	}
	var b strings.Builder
	for _, c := range strings.ToUpper(countryCode) {
		if c < 'A' || c > 'Z' {
			return ""
		}
		b.WriteRune(0x1F1E6 + c - 'A')
	}
	return b.String()
}

func source(dataSource string, rankType *int) string {
	switch {
	case dataSource == "matchmaking" && rankType != nil && *rankType == 11:
		return "Premier"
	case dataSource == "matchmaking":
		return "Matchmaking"
	case dataSource == "faceit":
		return "FACEIT"
	case dataSource == "":
		return ""
	}
	label := strings.ReplaceAll(dataSource, "_", " ")
	return strings.ToUpper(label[:1]) + label[1:]
}

func outcome(result string) string {
	switch result {
	case "win":
		return "W"
	case "loss":
		return "L"
	case "tie":
		return "T"
	}
	return "?"
}

func bar(v float64) string {
	return strconv.FormatFloat(math.Max(0, math.Min(100, v)), 'f', 1, 64)
}
