package chat

import (
	"testing"
)

func TestJoining(t *testing.T) {
	InitRooms(2)

	b1, r1, err := Join("test", "Foo")
	if err != nil {
		t.Fatalf("Could not join room: %s", err)
	}
	defer b1.Leave()

	if _, _, err = Join("test", "Foo"); err != NickAlreadyInUse {
		t.Fatalf("Got error \"%s\", expected \"%s\"", err, NickAlreadyInUse)
	}

	b2, r2, err := Join("test", "Bar")
	if err != nil {
		t.Fatalf("Could not join room a second time: %s", err)
	}
	defer b2.Leave()

	// since we joined to the same room, r1 and r2 should be equal
	if r1 != r2 {
		t.Error("r1 and r2 are not equal")
	}

	buddies := r1.ListBuddies()
	seen := make(map[string]bool)
	for _, b := range buddies {
		seen[b] = true
	}
	if !seen["Foo"] || !seen["Bar"] {
		t.Error("Foo or Bar missing in buddy list")
	}

	if _, _, err = Join("test", "Baz"); err != RoomIsFull {
		t.Fatalf("Got error \"%s\", expected \"%s\"", err, RoomIsFull)
	}
}

func TestLeaving(t *testing.T) {
	InitRooms(10)

	b1, r, err := Join("test", "Foo")
	if err != nil {
		t.Fatalf("Could not join room: %s", err)
	}

	b2, _, err := Join("test", "Bar")

	b1.Leave()

	buddies := r.ListBuddies()
	for _, b := range buddies {
		if b == "Foo" {
			t.Error("Foo is still in buddy list, even after leaving")
		}
	}

	b2.Leave()

	// The room is now empty. It should no longer exist.
	if _, ok := rooms["test"]; ok {
		t.Error("Room test still exists, although no user is left")
	}
}
