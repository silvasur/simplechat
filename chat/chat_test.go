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

	if _, _, err = Join("test", "Baz"); err != RoomIsFull {
		t.Fatalf("Got error \"%s\", expected \"%s\"", err, RoomIsFull)
	}
}

func TestLeaving(t *testing.T) {
	InitRooms(10)

	b1, _, err := Join("test", "Foo")
	if err != nil {
		t.Fatalf("Could not join room: %s", err)
	}

	b1.Leave()

	// The room is now empty. It should no longer exist.
	if _, ok := rooms["test"]; ok {
		t.Error("Room test still exists, although no user is left")
	}
}
