package game

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

const initiator = "initiator"
const otherPlayer = "otherPlayer"
const timeLimit = 60
const maxLength = 100
const entriesCount = 0
const entry = "some entry"

func TestStartGame(t *testing.T) {
	game := StartGame(initiator, []string{initiator, otherPlayer}, timeLimit, maxLength, entriesCount)

	if game.Turn != initiator {
		t.Error("it's not the initiator's turn")
	}

	if game.TimeLeft != timeLimit {
		t.Error("time left is not set right")
	}

	if game.MaxLength != maxLength {
		t.Error("max length is not set right")
	}

	if game.EntriesLeft != entriesCount {
		t.Error("entries left is not set right")
	}

	if game.MaxEntries != entriesCount {
		t.Error("max entries is not set right")
	}

	if game.Story == nil || game.Finished || game.VoteKick != nil {
		t.Error("fields are not initialized")
	}
}

func TestAddEntry(t *testing.T) {
	game := StartGame(initiator, []string{initiator, otherPlayer}, timeLimit, maxLength, entriesCount)

	err := game.AddEntry(entry, initiator)

	if err != nil {
		t.Error("add entry should pass with no error")
	}

	if game.Turn != otherPlayer {
		t.Error("the next turn is not set right")
	}

	if len(game.Story) != 1 {
		t.Error("story length is not 1 after adding a single entry")
	}

	if game.Story[0].Text != entry || game.Story[0].Player != initiator {
		t.Error("entry is not set right")
	}
}

func TestAddEntryIncorrectTurn(t *testing.T) {
	game := StartGame(initiator, []string{initiator, otherPlayer}, timeLimit, maxLength, entriesCount)

	err := game.AddEntry(entry, otherPlayer)

	if err == nil {
		t.Error("add entry should not pass with no error")
	}

	if err.Error() != "invalid entry - not this player's turn" {
		t.Error("wrong type of error is returned")
	}

	if game.Turn != initiator {
		t.Error("turn has changed with an invalid entry")
	}

	if len(game.Story) != 0 {
		t.Error("story length is not 0 after attempting to add an entry illegaly")
	}
}

func TestAddEntryMaxLengtViolation(t *testing.T) {
	maxLength := 5
	game := StartGame(initiator, []string{initiator, otherPlayer}, timeLimit, maxLength, entriesCount)

	err := game.AddEntry(entry, initiator)

	if err == nil {
		t.Error("add entry should not pass with no error")
	}

	if err.Error() != fmt.Sprintf("invalid entry - entry is above max length (%v)", maxLength) {
		t.Error("wrong type of error is returned")
	}

	if game.Turn != initiator {
		t.Error("turn has changed with an invalid entry")
	}

	if len(game.Story) != 0 {
		t.Error("story length is not 0 after attempting to add an entry illegaly")
	}
}

func TestAddEntryEndGameOnOneEntryLeft(t *testing.T) {
	game := StartGame(initiator, []string{initiator, otherPlayer}, timeLimit, maxLength, 1)

	err := game.AddEntry(entry, initiator)

	if err != nil {
		t.Error("add entry should pass with no error")
	}

	if !game.Finished {
		t.Error("game should be finished")
	}
}

func TestEndGame(t *testing.T) {
	game := StartGame(initiator, []string{initiator, otherPlayer}, timeLimit, maxLength, entriesCount)

	entriesToEndAfter := 5
	game.EndGame(entriesToEndAfter)

	if game.EntriesLeft != entriesToEndAfter {
		t.Error("entries left is not set right")
	}

	if game.MaxEntries != entriesToEndAfter {
		t.Error("max entries is not set right")
	}
}

func TestTriggerVoteKick(t *testing.T) {
	game := StartGame(initiator, []string{initiator, otherPlayer, "third player"}, timeLimit, maxLength, entriesCount)

	err := game.TriggerVoteKick(initiator, otherPlayer, 0.65, 60)

	if err != nil {
		t.Error("trigger votekick should pass with no error")
	}

	if game.VoteKick.Issuer != initiator {
		t.Error("issuer is not set right")
	}

	if game.VoteKick.Count != 0 {
		t.Error("count is not started from 0")
	}

	if game.VoteKick.Treshold != 2 {
		t.Error("threshold is not calculated correctly")
	}

	if game.VoteKick.Player != otherPlayer {
		t.Error("player to kick is not set right")
	}

	if game.VoteKick.TimeLeft != 60 {
		t.Error("time left is not set right")
	}
}

func TestTriggerVoteKickOnAFinishedGame(t *testing.T) {
	game := StartGame(initiator, []string{initiator, otherPlayer, "third player"}, timeLimit, maxLength, entriesCount)
	game.Finished = true

	err := game.TriggerVoteKick(initiator, otherPlayer, 0.65, 60)

	if err == nil {
		t.Error("trigger votekick on a finished game should return error")
	}

	if err.Error() != "there is no running game" {
		t.Error("wrong type of error is returned")
	}

	if game.VoteKick != nil {
		t.Error("vote kick was triggered with an illegal request")
	}
}

func TestTriggerVoteKickOnANonExistingPlayer(t *testing.T) {
	game := StartGame(initiator, []string{initiator, otherPlayer, "third player"}, timeLimit, maxLength, entriesCount)

	missingPlayer := "no-such-player"
	err := game.TriggerVoteKick(initiator, missingPlayer, 0.65, 60)

	if err == nil {
		t.Error("trigger votekick on a player that is not in the game should return error")
	}

	if err.Error() != fmt.Sprintf("player \"%s\" is not in the game", missingPlayer) {
		t.Error("wrong type of error is returned")
	}

	if game.VoteKick != nil {
		t.Error("vote kick was triggered with an illegal request")
	}
}

func TestTriggerVoteKickWithAnOngoingVote(t *testing.T) {
	game := StartGame(initiator, []string{initiator, otherPlayer, "third player"}, timeLimit, maxLength, entriesCount)

	game.TriggerVoteKick(initiator, otherPlayer, 0.65, 60)
	err := game.TriggerVoteKick(otherPlayer, initiator, 0.65, 60)

	if err == nil {
		t.Error("trigger votekick while another vote is running should return error")
	}

	if err.Error() != fmt.Sprintf("there is an ongoing vote to kick player \"%s\"", game.VoteKick.Player) {
		t.Error("wrong type of error is returned")
	}

	if game.VoteKick == nil {
		t.Error("vote kick should still be running after second trigger failed")
	}
}

func TestTriggerVoteKickEndsAfterTimeLimit(t *testing.T) {
	game := StartGame(initiator, []string{initiator, otherPlayer, "third player"}, timeLimit, maxLength, entriesCount)

	game.TriggerVoteKick(otherPlayer, initiator, 0.65, 1)

	time.Sleep(time.Millisecond * 1100)
	if game.VoteKick != nil {
		t.Error("vote kick should have ended after time limit")
	}
}

func TestVote(t *testing.T) {
	game := StartGame(initiator, []string{initiator, otherPlayer, "third player"}, timeLimit, maxLength, entriesCount)
	game.TriggerVoteKick(initiator, otherPlayer, 0.65, 60)

	err := game.Vote(initiator)

	if err != nil {
		t.Error("vote should not should return error")
	}

	if game.VoteKick.Count != 1 {
		t.Error("vote should have been counter")
	}

	if game.VoteKick.voted[0] != initiator {
		t.Error("voter should have been set in the voted list")
	}
}

func TestVoteOnAFinishedGame(t *testing.T) {
	game := StartGame(initiator, []string{initiator, otherPlayer, "third player"}, timeLimit, maxLength, entriesCount)

	game.TriggerVoteKick(initiator, otherPlayer, 0.65, 60)
	game.Finished = true

	err := game.Vote(initiator)

	if err == nil {
		t.Error("voting on a finished game should return error")
	}

	if err.Error() != "there is no running game" {
		t.Error("wrong type of error is returned")
	}
}

func TestVoteWithNoTriggeredVoteKick(t *testing.T) {
	game := StartGame(initiator, []string{initiator, otherPlayer, "third player"}, timeLimit, maxLength, entriesCount)

	err := game.Vote(initiator)

	if err == nil {
		t.Error("voting without an ongoing votekick should return error")
	}

	if err.Error() != "there is no ongoing vote" {
		t.Error("wrong type of error is returned")
	}
}

func TestVoteFromAPlayerOutsideTheGame(t *testing.T) {
	game := StartGame(initiator, []string{initiator, otherPlayer, "third player"}, timeLimit, maxLength, entriesCount)

	game.TriggerVoteKick(initiator, otherPlayer, 0.65, 60)
	missingPlayer := "no-such-player"

	err := game.Vote(missingPlayer)

	if err == nil {
		t.Error("voting with a player not in the game should return error")
	}

	if err.Error() != fmt.Sprintf("player \"%s\" cannot vote as he's not part of the game", missingPlayer) {
		t.Error("wrong type of error is returned")
	}
}

func TestVoteMoreThanOnce(t *testing.T) {
	game := StartGame(initiator, []string{initiator, otherPlayer, "third player"}, timeLimit, maxLength, entriesCount)

	game.TriggerVoteKick(initiator, otherPlayer, 0.65, 60)
	game.Vote(initiator)
	err := game.Vote(initiator)

	if err == nil {
		t.Error("voting with a player not in the game should return error")
	}

	if err.Error() != fmt.Sprintf("player \"%s\" has already voted for this vote", initiator) {
		t.Error("wrong type of error is returned")
	}
}

func TestVoteKickOnMeetingThreshold(t *testing.T) {
	thirdPlayer := "third player"
	game := StartGame(initiator, []string{initiator, otherPlayer, thirdPlayer}, timeLimit, maxLength, entriesCount)

	game.TriggerVoteKick(otherPlayer, initiator, 0.65, 60)

	game.Vote(otherPlayer)
	game.Vote(thirdPlayer)

	time.Sleep(1100 * time.Millisecond) // have to wait for thread to kick player

	if len(game.Players) != 2 {
		t.Error("the vote kicked player was not removed from the game")
	}

	if game.Turn == initiator {
		t.Error("the vote kicked player was not removed from current turn")
	}

	err := game.AddEntry(entry, initiator)
	if err == nil {
		t.Error("the vote kicked player can still submit entried")
	}
}

func TestTurnsAreSwappedCorrectlyOnTimeRunningOut(t *testing.T) {
	game := StartGame(initiator, []string{initiator, otherPlayer}, 1, maxLength, entriesCount)

	game.setNextTurn()
	if game.Turn != otherPlayer {
		t.Error("the turn did not set correctly after time limit was reached")
	}

	game.setNextTurn()
	if game.Turn != initiator {
		t.Error("the turn did not repeat the first player after players ran out")
	}
}

func TestGameEndsWithNoLeftPlayers(t *testing.T) {
	game := StartGame(initiator, []string{initiator}, 1, maxLength, entriesCount)
	game.Kick(initiator)
	game.setNextTurn()
	if !game.Finished {
		t.Error("the game has no players left but is still not finished")
	}

}

func TestTurnRepeatedAfterAllPlayersPlay(t *testing.T) {
	game := StartGame(initiator, []string{initiator, otherPlayer}, 2, maxLength, entriesCount)

	time.Sleep(3 * time.Second)
	if game.Turn != otherPlayer {
		t.Error("the turn did not set correctly after time limit was reached")
	}
}

func TestKickPlayerNotInTheGame(t *testing.T) {
	game := StartGame(initiator, []string{initiator}, 2, maxLength, entriesCount)

	err := game.Kick(otherPlayer)

	if err == nil {
		t.Error("kick on a missing player should return error")
	}

	if err.Error() != fmt.Sprintf("player \"%s\" is not part of the game", otherPlayer) {
		t.Error("wrong type of error is returned")

	}
}

func TestGameStringMethod(t *testing.T) {
	game := StartGame(initiator, []string{initiator, otherPlayer}, timeLimit, maxLength, entriesCount)

	gameStr := game.String()
	if !strings.Contains(gameStr, fmt.Sprintf("Players in the game: %s, %s", initiator, otherPlayer)) ||
		!strings.Contains(gameStr, fmt.Sprintf(`Next turn: Player "%s"`, initiator)) ||
		!strings.Contains(gameStr, fmt.Sprintf("Max length: %d symbols", maxLength)) ||
		!strings.Contains(gameStr, "Time left:") {
		t.Error("string method missed some output")
	}
}

func TestGameStringMethodOnAGameWithEntriesLeft(t *testing.T) {
	entriesCount := 5
	game := StartGame(initiator, []string{initiator, otherPlayer}, timeLimit, maxLength, entriesCount)

	gameStr := game.String()
	if !strings.Contains(gameStr, fmt.Sprintf("Players in the game: %s, %s", initiator, otherPlayer)) ||
		!strings.Contains(gameStr, fmt.Sprintf(`Next turn: Player "%s"`, initiator)) ||
		!strings.Contains(gameStr, fmt.Sprintf("Max length: %d symbols", maxLength)) ||
		!strings.Contains(gameStr, "Time left:") ||
		!strings.Contains(gameStr, fmt.Sprintf("Entires left: %d", entriesCount)) {
		t.Error("string method missed some output")
	}
}

func TestGameStringMethodOnAGameWithOnlyOneEntryLeft(t *testing.T) {
	entriesCount := 1
	game := StartGame(initiator, []string{initiator, otherPlayer}, timeLimit, maxLength, entriesCount)

	if !strings.Contains(game.String(), "Next entry will be the story ending. Make it a good one!") {
		t.Error("string method missed some output")
	}
}

func TestGameStringMethodOnAGameWithVote(t *testing.T) {
	game := StartGame(initiator, []string{initiator, otherPlayer}, timeLimit, maxLength, entriesCount)
	game.TriggerVoteKick(initiator, otherPlayer, 0.65, 60)

	gameStr := game.String()
	if !strings.Contains(gameStr, fmt.Sprintf("Players in the game: %s, %s", initiator, otherPlayer)) ||
		!strings.Contains(gameStr, fmt.Sprintf(`Next turn: Player "%s"`, initiator)) ||
		!strings.Contains(gameStr, fmt.Sprintf("Max length: %d symbols", maxLength)) ||
		!strings.Contains(gameStr, "Time left:") ||
		!strings.Contains(gameStr, fmt.Sprintf(game.VoteKick.String())) {
		t.Error("string method missed some output")
	}
}

func TestGameStringOnAFinishedGame(t *testing.T) {
	game := StartGame(initiator, []string{initiator, otherPlayer}, timeLimit, maxLength, entriesCount)
	game.Finished = true

	if !strings.Contains(game.String(), "The game has finished. You can now start the next one!") {
		t.Error("string method missed some output")
	}
}

func TestGameStringWithAStory(t *testing.T) {
	game := StartGame(initiator, []string{initiator, otherPlayer}, timeLimit, maxLength, entriesCount)
	game.AddEntry(entry, initiator)

	if !strings.Contains(game.String(), game.Story[0].String()) {
		t.Error("string method missed some output")
	}
}
