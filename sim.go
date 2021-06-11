package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
)

var (
	poolSize  = flag.Int("pool_size", 32, "Number of players in the pool")
	steps     = flag.Int("steps", 100000, "Number of steps to simulate")
	packValue = flag.Float64("pack_value", 200, "Gem value of a pack")
	event     = flag.String("event", "premier", "Event type (premier, quick)")
)

type player struct {
	wins, losses int
}

func main() {
	flag.Parse()
	rand.Seed(time.Now().UnixNano())

	var entryFee float64
	var gems map[player]float64
	var packs map[player]float64

	switch *event {
	case "premier":
		entryFee = 1500
		gems = map[player]float64{
			{wins: 0, losses: 3}: 50,
			{wins: 1, losses: 3}: 100,
			{wins: 2, losses: 3}: 250,
			{wins: 3, losses: 3}: 1000,
			{wins: 4, losses: 3}: 1400,
			{wins: 5, losses: 3}: 1600,
			{wins: 6, losses: 3}: 1800,
			{wins: 7, losses: 0}: 2200,
			{wins: 7, losses: 1}: 2200,
			{wins: 7, losses: 2}: 2200,
		}
		packs = map[player]float64{
			{wins: 0, losses: 3}: 1,
			{wins: 1, losses: 3}: 1,
			{wins: 2, losses: 3}: 2,
			{wins: 3, losses: 3}: 2,
			{wins: 4, losses: 3}: 3,
			{wins: 5, losses: 3}: 4,
			{wins: 6, losses: 3}: 5,
			{wins: 7, losses: 0}: 6,
			{wins: 7, losses: 1}: 6,
			{wins: 7, losses: 2}: 6,
		}

	case "quick":
		entryFee = 750
		gems = map[player]float64{
			{wins: 0, losses: 3}: 50,
			{wins: 1, losses: 3}: 100,
			{wins: 2, losses: 3}: 200,
			{wins: 3, losses: 3}: 300,
			{wins: 4, losses: 3}: 450,
			{wins: 5, losses: 3}: 650,
			{wins: 6, losses: 3}: 850,
			{wins: 7, losses: 0}: 950,
			{wins: 7, losses: 1}: 950,
			{wins: 7, losses: 2}: 950,
		}
		packs = map[player]float64{
			{wins: 0, losses: 3}: 1.2,
			{wins: 1, losses: 3}: 1.22,
			{wins: 2, losses: 3}: 1.24,
			{wins: 3, losses: 3}: 1.26,
			{wins: 4, losses: 3}: 1.3,
			{wins: 5, losses: 3}: 1.35,
			{wins: 6, losses: 3}: 1.4,
			{wins: 7, losses: 0}: 2,
			{wins: 7, losses: 1}: 2,
			{wins: 7, losses: 2}: 2,
		}

	default:
		fmt.Printf("Event type %q not suppported\n", *event)
		os.Exit(1)
	}

	if *poolSize < 2 {
		fmt.Printf("Pool size must be at least 2\n")
		os.Exit(1)
	}

	players := make([]player, *poolSize, *poolSize)
	var results []player
	for s := 0; s < *steps; s++ {
		// Randomly pick a winner and a loser, retry if they are the same
		i := rand.Intn(len(players))
		j := rand.Intn(len(players))
		if i == j {
			s--
			continue
		}

		players[i].wins++
		players[j].losses++

		// Every time a player reaches 7 wins or 3 losses, we replace it with a new player
		if players[i].wins >= 7 {
			results = append(results, players[i])
			players[i] = player{}
		}
		if players[j].losses >= 3 {
			results = append(results, players[j])
			players[j] = player{}
		}
	}

	finishes := float64(len(results))
	totalEntryFee := entryFee * finishes
	wins := 0
	gemPayout := 0.0
	packPayout := 0.0
	stats := map[player]int{}
	for _, p := range results {
		wins += p.wins
		gemPayout += gems[p]
		packPayout += packs[p]
		stats[p]++
	}

	fmt.Printf("Number of players = %d\n", len(results))
	fmt.Printf("Results ({win loss}: count) = %v\n", stats)
	fmt.Printf("Avg wins = %.2f (total %d)\n", float64(wins)/finishes, wins)

	fmt.Println()
	fmt.Printf("Entry fee = %.0f (total %.0f)\n", entryFee, totalEntryFee)
	fmt.Printf("Avg gem payout = %.0f (total %.0f)\n", gemPayout/finishes, gemPayout)
	fmt.Printf("Avg packs payout = %.2f (total %.0f)\n", packPayout/finishes, packPayout)
	ev1 := totalEntryFee - gemPayout
	fmt.Printf("House EV = %.2f (total %.0f)\n", ev1/finishes, ev1)
	ev2 := ev1 - packPayout**packValue
	fmt.Printf("House EV w/ pack valued at %.0f gems = %.2f (total %.0f)\n", *packValue, ev2/finishes, ev2)
}
