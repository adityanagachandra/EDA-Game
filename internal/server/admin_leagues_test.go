package server

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/chuma-beep/mathua/internal/storage"
)

// The owner's view of the league pyramid.
//
// The public league endpoint answers "where am I", which leaves an administrator unable to see a
// cohort at all: a report saying "my league is broken" cannot be located, and a tier holding one
// learner looks identical to a tier that does not exist. These tests cover the read that closes
// that gap, and the two ways it can be wrong: exposing accounts to someone who may not read them,
// and trusting a date that is used as a join key.

type adminLeaguesBody struct {
	Week    string `json:"week"`
	Promote int    `json:"promote"`
	Demote  int    `json:"demote"`
	Seats   int    `json:"seats"`
	Leagues []struct {
		Tier     string `json:"tier"`
		LeagueNo int    `json:"league_no"`
		Size     int    `json:"size"`
		Movers   int    `json:"movers"`
		Members  []struct {
			Rank           int    `json:"rank"`
			Name           string `json:"name"`
			Username       string `json:"username"`
			StudentID      string `json:"student_id"`
			XP             int    `json:"xp"`
			AvatarURL      string `json:"avatar_url"`
			AvatarDicebear *struct {
				Style string `json:"style"`
				Seed  string `json:"seed"`
			} `json:"avatar_dicebear"`
			AvatarPreset *int `json:"avatar_preset"`
		} `json:"members"`
	} `json:"leagues"`
}

// seedLeagueXP gives the env's student enough lifetime XP to be seated in the current week, and
// a second account so the pyramid has more than one row.
func seedLeagueXP(t *testing.T, e *adminEnv) string {
	t.Helper()
	if err := e.store.AwardXP(e.studentID, 120, "test", "league:ada"); err != nil {
		t.Fatalf("award xp: %v", err)
	}
	// The second account is created through the same signup route the rest of the suite uses, so
	// the roster carries a name and a username rather than a bare id.
	rec := e.do(t, "POST", "/api/auth/signup", "",
		`{"name":"Grace","username":"grace","password":"Engine!n1","email":"grace@example.com"}`)
	if rec.Code != 200 {
		t.Fatalf("signup grace: %d %s", rec.Code, rec.Body.String())
	}
	var out struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("decode signup: %v", err)
	}
	// Resolve the new account's id through the admin list rather than guessing it.
	rec = e.do(t, "GET", "/api/admin/users?q=grace", e.adminTok, "")
	if rec.Code != 200 {
		t.Fatalf("find grace: %d %s", rec.Code, rec.Body.String())
	}
	var list struct {
		Users []struct {
			ID string `json:"id"`
		} `json:"users"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &list); err != nil {
		t.Fatalf("decode users: %v", err)
	}
	if len(list.Users) != 1 {
		t.Fatalf("found %d accounts named grace, want 1", len(list.Users))
	}
	if err := e.store.AwardXP(list.Users[0].ID, 40, "test", "league:grace"); err != nil {
		t.Fatalf("award xp to grace: %v", err)
	}
	return list.Users[0].ID
}

func TestAdminLeaguesListsThePyramidAndOneRoster(t *testing.T) {
	e := newAdminEnv(t)
	graceID := seedLeagueXP(t, e)
	if err := e.store.SetAvatarImage(e.studentID, "image/png", []byte("custom-photo")); err != nil {
		t.Fatalf("set custom avatar: %v", err)
	}
	if err := e.store.UpdateSettings(e.studentID, `{"avatar_custom":true,"avatar_version":7,"avatar_preset":2}`); err != nil {
		t.Fatalf("set ada avatar settings: %v", err)
	}
	if err := e.store.UpdateSettings(graceID, `{"avatar_dicebear":{"style":"notionists","seed":"grace-seed"}}`); err != nil {
		t.Fatalf("set grace avatar settings: %v", err)
	}

	rec := e.do(t, "GET", "/api/admin/leagues", e.adminTok, "")
	if rec.Code != 200 {
		t.Fatalf("leagues: %d %s", rec.Code, rec.Body.String())
	}
	var body adminLeaguesBody
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Week == "" {
		t.Fatal("week is empty")
	}
	if body.Promote != storage.PromoteCount || body.Demote != storage.DemoteCount {
		t.Errorf("promote/demote = %d/%d, want %d/%d",
			body.Promote, body.Demote, storage.PromoteCount, storage.DemoteCount)
	}
	if len(body.Leagues) != 1 {
		t.Fatalf("leagues = %+v, want one cohort for the first competition week", body.Leagues)
	}
	first := body.Leagues[0]
	if first.Tier != storage.LeagueTiers[0] {
		t.Errorf("tier = %q, want %q", first.Tier, storage.LeagueTiers[0])
	}
	if first.Size != 2 || body.Seats != 2 {
		t.Errorf("size = %d, seats = %d, want 2 and 2", first.Size, body.Seats)
	}
	// A short cohort must not advertise movers. With two learners, moving any of them would
	// shuffle the same people between tiers, which is what minLeagueToMove prevents.
	if first.Movers != 0 {
		t.Errorf("movers = %d for a cohort of two, want 0", first.Movers)
	}
	if len(first.Members) != 0 {
		t.Errorf("members = %d in the pyramid read, want none: rosters are fetched per league",
			len(first.Members))
	}

	// One league, with standings, ordered by XP and carrying the account id so the page can link
	// to it.
	rec = e.do(t, "GET", fmt.Sprintf("/api/admin/leagues?tier=%s&league=%d", first.Tier, first.LeagueNo), e.adminTok, "")
	if rec.Code != 200 {
		t.Fatalf("roster: %d %s", rec.Code, rec.Body.String())
	}
	body = adminLeaguesBody{}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode roster: %v", err)
	}
	if len(body.Leagues) != 1 || len(body.Leagues[0].Members) != 2 {
		t.Fatalf("roster = %+v, want the one league and its two members", body.Leagues)
	}
	members := body.Leagues[0].Members
	if members[0].XP != 120 || members[0].StudentID != e.studentID {
		t.Errorf("first row = %+v, want Ada with 120 XP", members[0])
	}
	if members[0].Rank != 1 || members[1].Rank != 2 {
		t.Errorf("ranks = %d, %d, want 1, 2", members[0].Rank, members[1].Rank)
	}
	if members[1].StudentID != graceID || members[1].XP != 40 {
		t.Errorf("second row = %+v, want Grace with 40 XP", members[1])
	}
	if members[1].Username != "grace" || members[1].Name != "Grace" {
		t.Errorf("second row identity = %q/%q, want Grace/grace", members[1].Name, members[1].Username)
	}
	wantCustom := "/api/avatar/" + e.studentID + "?v=7"
	if members[0].AvatarURL != wantCustom {
		t.Errorf("custom avatar url = %q, want %q", members[0].AvatarURL, wantCustom)
	}
	if members[0].AvatarPreset == nil || *members[0].AvatarPreset != 2 {
		t.Errorf("legacy preset = %v, want 2", members[0].AvatarPreset)
	}
	if members[1].AvatarDicebear == nil ||
		members[1].AvatarDicebear.Style != "notionists" ||
		members[1].AvatarDicebear.Seed != "grace-seed" {
		t.Errorf("dicebear pick = %+v, want notionists/grace-seed", members[1].AvatarDicebear)
	}
}

func TestAdminLeaguesRefusesALearnerAndABadWeek(t *testing.T) {
	e := newAdminEnv(t)
	seedLeagueXP(t, e)

	if rec := e.do(t, "GET", "/api/admin/leagues", e.studentTok, ""); rec.Code != 403 {
		t.Errorf("learner read = %d, want 403", rec.Code)
	}
	if rec := e.do(t, "GET", "/api/admin/leagues", "", ""); rec.Code != 401 {
		t.Errorf("anonymous read = %d, want 401", rec.Code)
	}
	// A week key is a join key into league_members, so anything that is not a Monday is refused
	// before it reaches the store.
	for _, week := range []string{"2026-10-22", "not-a-date", "2026-13-01"} {
		rec := e.do(t, "GET", "/api/admin/leagues?week="+week, e.adminTok, "")
		if rec.Code != 400 {
			t.Errorf("week %q = %d, want 400", week, rec.Code)
		}
	}
	// Half a filter is refused rather than silently ignored: a page that asked for tier=bronze
	// and got the whole pyramid would show the wrong thing.
	if rec := e.do(t, "GET", "/api/admin/leagues?tier=bronze", e.adminTok, ""); rec.Code != 400 {
		t.Errorf("tier without league = %d, want 400", rec.Code)
	}
	// A league that does not exist is a 404 rather than an empty roster.
	rec := e.do(t, "GET", "/api/admin/leagues?tier=diamond&league=9", e.adminTok, "")
	if rec.Code != 404 {
		t.Errorf("missing league = %d, want 404", rec.Code)
	}
}
