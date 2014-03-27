package chat

import (
	"fmt"
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

	if _, _, err = Join("test", ""); err != EmptyNick {
		t.Fatalf("Got error \"%s\", expected \"%s\"", err, EmptyNick)
	}

	if _, _, err = Join("test", "abcdefghijklmnopqrstuvwxyz"); err != NickTooLong {
		t.Fatalf("Got error \"%s\", expected \"%s\"", err, NickTooLong)
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

func checkMsg(t *testing.T, m Message, typ MsgType, user string, txt string) {
	fail := fmt.Sprintf("Expected a %s message from %s", typ, user)
	if txt != "" {
		fail += fmt.Sprintf(" with message '%s'.", txt)
	}

	failed := false

	if m.Type != typ {
		failed = true
		fail += fmt.Sprintf(" Type is wrong (%s)", m.Type)
	}
	if m.User != user {
		failed = true
		fail += fmt.Sprintf(" User is wrong (%s)", m.User)
	}
	if txt != "" && m.Text != txt {
		failed = true
		fail += fmt.Sprintf(" Text is wrong (%s)", m.Text)
	}

	if failed {
		t.Error(fail)
	}
}

func TestChatting(t *testing.T) {
	InitRooms(10)

	// In this test, we will ignore the errors of Join(), since we already tested that stuff in the tests above.

	b1, _, _ := Join("test", "Foo")
	b2, _, _ := Join("test", "Bar")

	checkMsg(t, <-b1.Receive, MsgJoin, "Bar", "")
	checkMsg(t, <-b2.Receive, MsgJoin, "Bar", "")

	b2.Say("Hello")
	checkMsg(t, <-b1.Receive, MsgChat, "Bar", "Hello")
	checkMsg(t, <-b2.Receive, MsgChat, "Bar", "Hello")

	b1.Say(":)")
	checkMsg(t, <-b1.Receive, MsgChat, "Foo", ":)")
	checkMsg(t, <-b2.Receive, MsgChat, "Foo", ":)")

	b3, _, _ := Join("test", "Baz")

	checkMsg(t, <-b1.Receive, MsgJoin, "Baz", "")
	checkMsg(t, <-b2.Receive, MsgJoin, "Baz", "")
	checkMsg(t, <-b3.Receive, MsgJoin, "Baz", "")

	b3.Say("!")

	checkMsg(t, <-b1.Receive, MsgChat, "Baz", "!")
	checkMsg(t, <-b2.Receive, MsgChat, "Baz", "!")
	checkMsg(t, <-b3.Receive, MsgChat, "Baz", "!")

	b2.Leave()

	checkMsg(t, <-b1.Receive, MsgLeave, "Bar", "")
	checkMsg(t, <-b3.Receive, MsgLeave, "Bar", "")
}
