package main

import (
	"github.com/Simply-Box/go-chess-bot/board"
	"github.com/Simply-Box/go-chess-bot/movegen"

	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("Chess Bot is now online")
	fmt.Println("_________________")

	// Import the starting gamestate of a new match
	gameState := board.FENtoGameState(board.GetStartingFEN())
	board.PrintBoard(gameState.Board)
	fmt.Println("_________________")

	// Set to true if you want to test yourself
	humanPlaying := false

	// Create a reader for input
	reader := bufio.NewReader(os.Stdin)

	for {
		// Generate moves
		allMoves := movegen.GenerateAllMoves(&gameState)

		// Check game end conditions and stop loop
		if board.CheckGameEnd(gameState.Results) {
			board.PrintResult(gameState.Results, gameState.WhiteToMove)
			break
		}

		if humanPlaying {
			// Take a human input
			fmt.Print("Move (e.g. e2e4): ")
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading input:", err)
				// skip to next iteration to get new input
					// maybe move this code up so you don't have to redo above code
				continue
			}
			move := strings.TrimSpace(input)
			if move == "exit" { break }

			// Check and apply the human move
			err = board.CheckHumanMove(&gameState, &allMoves, move)
			
			// prints error msg 
			if err != nil {
				fmt.Println(err)
				continue
			}
		} else {
			// Selects the computer move
			move := allMoves[1] // takes the second move with index 1

			// Applies the computer move
			board.ApplyMove(&gameState, move)
		}

		board.PrintBoard(gameState.Board)
		fmt.Println("_________________")
	}
}
