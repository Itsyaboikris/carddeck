//go:generate stringer -type=Suit,Rank

package carddeck

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

type Suit uint8 

const (
	Spade Suit = iota 
	Diamond
	Club
	Heart
	Joker
)


var suits = [...]Suit{Spade, Diamond, Club, Heart}

type Rank uint8

const (
	_ Rank = iota
	Ace
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

const (
	minRank = Ace
	maxRank = King
)

type Card struct {
	Suit
	Rank
}

func (c Card) String() string {
	if c.Suit == Joker {
		return c.Suit.String()
	}

	return fmt.Sprintf("%s of %ss", c.Rank.String(), c.Suit.String())
}

func New(opts ...func([]Card) []Card) []Card {
	var cards []Card

	for _, suit := range suits {
		for rank := minRank; rank <= maxRank; rank++ {
			cards = append(cards, Card{Suit: suit, Rank: rank})
		}
	}

	for _, opt := range opts {
		cards = opt(cards)
	}

	return cards
}

func DefaultSort(cards []Card) []Card {
	sort.Slice(cards, Less(cards))
	return cards
}

func Sort(less func(cards []Card) func(i, j int) bool) func([]Card) []Card {
	return func(cards []Card) []Card {
		sort.Slice(cards, less(cards))
		return cards
	}
}

func Less(cards []Card) func(i, j int) bool {
	return func(i, j int) bool {
		return absRank(cards[i]) < absRank(cards[j])
	}
}

func absRank(c Card) int {
	return int(c.Suit)*int(maxRank) + int(c.Rank)
}

var shuffleRand = rand.New(rand.NewSource(time.Now().Unix()))

func Shuffle(cards []Card) []Card {
	retrn := make([]Card, len(cards))

	perm := shuffleRand.Perm(len(cards))

	for i, j := range perm {
		retrn[i] = cards[j]
	}

	return retrn
}

func Jokers(n int) func([]Card) []Card {
	return func(cards []Card) []Card {
		for i := 0; i < n; i++ {
			cards = append(cards, Card{
				Rank: Rank(i),
				Suit: Joker,
			})
		}
		return cards
	}
}

func Filter(f func(card Card) bool) func([]Card) []Card {
	return func(cards []Card) []Card {
		var retrn []Card

		for _, c := range cards {
			if !f(c) {
				retrn = append(retrn, c)
			}
		}

		return retrn
	}
}

func Deck(n int) func([]Card) []Card {
	return func(cards []Card) []Card {
		var retrn []Card

		for i := 0; i < n; i++ {
			retrn = append(retrn, cards...)
		}

		return retrn
	}
}
