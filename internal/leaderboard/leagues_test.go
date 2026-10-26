package leaderboard

import (
	"fmt"
	"testing"
	"time"

	"github.com/chuma-beep/mathua/internal/storage"
)

func leagueRepo(t *testing.T) *storage.SQLiteStore {
	t.Helper()
	repo, err := storage.NewSQLiteStore(":memory:")
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	t.Cleanup(func() { repo.Close() })
	return repo
}

// TestMyLeague_SeatsAndRanksTheCallersCohort covers the whole read: the first
// reader of a week bootstraps the leagues, the caller is seated, and the
// standings are ordered by this week's XP with the caller's place returned
// separately rather than injected into the list.
func TestMyLeague_SeatsAndRanksTheCallersCohort(t *testing.T) {
	repo := leagueRepo(t)
	var ids []string
	for i := 0; i < 8; i++ {
		st, err := repo.CreateStudent(fmt.Sprintf("learner%d", i))
		if err != nil {
			t.Fatalf("create %d: %v", i, err)
		}
		if err := repo.SetUsername(st.ID, fmt.Sprintf("learner%d", i)); err != nil {
			t.Fatalf("username %d: %v", i, err)
		}
		if err := repo.AwardXP(st.ID, 10*(i+1), "test", fmt.Sprintf("week:%d", i)); err != nil {
			t.Fatalf("award %d: %v", i, err)
		}
		ids = append(ids, st.ID)
	}

	// AwardXP stamps the live week. Pinning a past Monday made this test pass
	// until the calendar moved and then fail with an empty board, so the week
	// is read from the same clock the awards used.
	now := time.Now().UTC()
	board, err := NewComputer(repo).MyLeague(ids[0], now)
	if err != nil {
		t.Fatalf("my league: %v", err)
	}
	if board.Tier != "bronze" || board.LeagueNo != 1 {
		t.Errorf("seat = %s/%d, want bronze/1", board.Tier, board.LeagueNo)
	}
	if board.Size != 8 || len(board.Members) != 8 {
		t.Fatalf("size = %d members = %d, want 8/8", board.Size, len(board.Members))
	}
	if want := WeekStart(now); !board.WeekStart.Equal(want) {
		t.Errorf("week start = %s, want %s", board.WeekStart, want)
	}
	if board.Me == nil || board.Me.Rank != 8 || board.Me.XP != 10 {
		t.Errorf("me = %+v, want rank 8 with 10 XP", board.Me)
	}
	// Highest XP first, and each row carries the public username instead of an
	// account id.
	if board.Members[0].XP != 80 || board.Members[0].ID != "learner7" {
		t.Errorf("first row = %+v, want learner7 at 80 XP", board.Members[0])
	}
	if !board.Members[7].IsMe {
		t.Error("caller's line is not marked")
	}
	// The rules travel with the payload so the interface cannot describe a
	// different threshold from the one applied.
	if board.Promote != storage.LeagueMovers(8) || board.Demote != board.Promote {
		t.Errorf("rules = %d/%d, want %d", board.Promote, board.Demote, storage.LeagueMovers(8))
	}
}

// A league line follows the same rule as the public board, and the caller
// still recognises their own anonymised seat.
func TestMyLeague_AnonymisesOptedOutMembers(t *testing.T) {
	repo := leagueRepo(t)
	var ids []string
	for i := 0; i < 8; i++ {
		st, err := repo.CreateStudent(fmt.Sprintf("learner%d", i))
		if err != nil {
			t.Fatalf("create %d: %v", i, err)
		}
		if err := repo.SetUsername(st.ID, fmt.Sprintf("learner%d", i)); err != nil {
			t.Fatalf("username %d: %v", i, err)
		}
		if err := repo.AwardXP(st.ID, 10*(i+1), "test", fmt.Sprintf("privacy:%d", i)); err != nil {
			t.Fatalf("award %d: %v", i, err)
		}
		ids = append(ids, st.ID)
	}
	hidden := ids[7] // highest XP, rank 1
	if err := repo.SetAvatarURL(hidden, "https://photos.example/learner7.jpg"); err != nil {
		t.Fatalf("avatar url: %v", err)
	}
	if err := repo.UpdateSettings(hidden, `{"leaderboard_privacy":"hidden","avatar_custom":true,"avatar_dicebear":{"style":"bottts","seed":"learner7"},"avatar_version":4}`); err != nil {
		t.Fatalf("settings: %v", err)
	}

	now := time.Now().UTC()
	computer := NewComputer(repo)
	board, err := computer.MyLeague(ids[0], now)
	if err != nil {
		t.Fatalf("my league: %v", err)
	}
	if len(board.Members) != 8 {
		t.Fatalf("members = %d, want 8", len(board.Members))
	}
	gone := board.Members[0]
	if gone.ID != "anonymous-1" || gone.Name != "Anonymous" {
		t.Errorf("hidden member = %+v, want anonymous-1/Anonymous", gone)
	}
	if gone.Avatar == nil || gone.Avatar.Kind != "generated" || gone.Avatar.Style != "notionists" {
		t.Errorf("hidden member avatar = %+v, want a generic generated face", gone.Avatar)
	}
	if gone.Avatar != nil && (gone.Avatar.URL != "" || gone.Avatar.Version != 0) {
		t.Errorf("hidden member avatar leaks learner data: %+v", gone.Avatar)
	}
	if gone.IsMe {
		t.Error("another learner's seat marked as the caller")
	}
	// The named members are untouched.
	if board.Members[1].ID != "learner6" || board.Members[1].Name != "learner6" {
		t.Errorf("second seat = %+v, want learner6 unchanged", board.Members[1])
	}

	// The caller's own seat is anonymised but still theirs.
	own, err := computer.MyLeague(hidden, now)
	if err != nil {
		t.Fatalf("own league: %v", err)
	}
	if own.Me == nil || own.Me.Rank != 1 {
		t.Fatalf("me = %+v, want rank 1", own.Me)
	}
	mine := own.Members[0]
	if !mine.IsMe || mine.Name != "Anonymous" || mine.ID != "anonymous-1" {
		t.Errorf("own seat = %+v, want anonymous-1/Anonymous and IsMe", mine)
	}
}

func TestMyLeague_GuestAndZeroXPGetNoSeat(t *testing.T) {
	repo := leagueRepo(t)
	now := time.Date(2026, 10, 21, 12, 0, 0, 0, time.UTC)
	computer := NewComputer(repo)

	guest, err := computer.MyLeague("guest_abc", now)
	if err != nil {
		t.Fatalf("guest league: %v", err)
	}
	if guest.Tier != "" || len(guest.Members) != 0 || guest.Me != nil {
		t.Errorf("guest got a league: %+v", guest)
	}

	quiet, err := repo.CreateStudent("quiet")
	if err != nil {
		t.Fatalf("create quiet: %v", err)
	}
	empty, err := computer.MyLeague(quiet.ID, now)
	if err != nil {
		t.Fatalf("zero-xp league: %v", err)
	}
	if empty.Tier != "" || empty.Me != nil {
		t.Errorf("zero-XP learner got a league: %+v", empty)
	}
}

// TestMyLeague_SecondReadDoesNotReseatAnyone pins the retry case: two reads of
// the same week return the same seats and the same standings.
func TestMyLeague_SecondReadDoesNotReseatAnyone(t *testing.T) {
	repo := leagueRepo(t)
	var ids []string
	for i := 0; i < 12; i++ {
		st, err := repo.CreateStudent(fmt.Sprintf("learner%d", i))
		if err != nil {
			t.Fatalf("create %d: %v", i, err)
		}
		if err := repo.AwardXP(st.ID, i+1, "test", fmt.Sprintf("twice:%d", i)); err != nil {
			t.Fatalf("award %d: %v", i, err)
		}
		ids = append(ids, st.ID)
	}
	now := time.Date(2026, 10, 21, 12, 0, 0, 0, time.UTC)
	computer := NewComputer(repo)

	first, err := computer.MyLeague(ids[3], now)
	if err != nil {
		t.Fatalf("first: %v", err)
	}
	second, err := computer.MyLeague(ids[3], now)
	if err != nil {
		t.Fatalf("second: %v", err)
	}
	if first.Tier != second.Tier || first.LeagueNo != second.LeagueNo || first.Size != second.Size {
		t.Errorf("seat changed between reads: %+v then %+v", first, second)
	}
	if first.Me == nil || second.Me == nil || first.Me.Rank != second.Me.Rank {
		t.Errorf("rank changed between reads: %+v then %+v", first.Me, second.Me)
	}
}

func TestMyLeague_EmptyOrGuestStudentID(t *testing.T) {
	repo := leagueRepo(t)
	comp := NewComputer(repo)
	now := time.Date(2026, 10, 21, 12, 0, 0, 0, time.UTC)

	// Empty student ID
	emptyRes, err := comp.MyLeague("", now)
	if err != nil {
		t.Fatalf("MyLeague empty student: %v", err)
	}
	if len(emptyRes.Members) != 0 || emptyRes.Me != nil {
		t.Errorf("expected empty members for empty student, got %+v", emptyRes)
	}

	// Guest student ID
	guestRes, err := comp.MyLeague("guest-abc-123", now)
	if err != nil {
		t.Fatalf("MyLeague guest student: %v", err)
	}
	if len(guestRes.Members) != 0 || guestRes.Me != nil {
		t.Errorf("expected empty members for guest student, got %+v", guestRes)
	}
}
