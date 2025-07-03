package main

import (
	"github.com/Simply-Box/go-chess-bot/board"
	//"chessEngine/movegen"
	"fmt"
	"bufio"
	"os"
	"strings"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func main() {
	// Start pprof server in the background, http://localhost:6060/debug/pprof/
	// In terminal after started main.go in another terminal:
	//  go tool pprof http://localhost:6060/debug/pprof/profile?seconds=6
	// 	after x seconds you can run commands like: top, web(needs https://graphviz.org/download/)
	go func() {
		log.Println("Starting pprof on http://localhost:6060/debug/pprof/")
		err := http.ListenAndServe("localhost:6060", nil)
		if err != nil {
			log.Fatalf("pprof failed: %v", err)
		}
	}()

	time.Sleep(3 * time.Second) // waits for things like performance checks
  
	log.Println(http.ListenAndServe("localhost:6060", nil))
	fmt.Println("Starting Chess Engine...")
	gameState := board.FENtoGameState(board.GetStartingFEN())
	board.PrintBoard(gameState.Board)
	reader := bufio.NewReader(os.Stdin)

	// Set to true if you want to play yourself
	humanPlaying := false
	if humanPlaying {
		for {
			allMoves := board.GenerateAllMoves(&gameState)

			if board.CheckGameEnd(gameState.Results){
				board.PrintResult(gameState.Results, gameState.WhiteToMove)
				break
			}

			fmt.Print("Move (e.g. e2e4): ")
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading input:", err)
				continue
			}

			move := strings.TrimSpace(input)

			if move == "exit" {
				break
			}

			fmt.Println("You entered:", move)

			err = board.CheckHumanMove(&gameState, &allMoves, move)
			if err != nil {
				fmt.Println(err)
				continue
			}

			board.PrintBoard(gameState.Board)
		}
	} else {
		for {
			allMoves := board.GenerateAllMoves(&gameState)

			if board.CheckGameEnd(gameState.Results){
				board.PrintResult(gameState.Results, gameState.WhiteToMove)
				break
			}

			move := allMoves[1] // takes the second move

			board.ApplyMove(&gameState, move)

			board.PrintBoard(gameState.Board)
			fmt.Println("_________________")
		}
	}
	//time.Sleep(20 * time.Second) // waits for things like performance checks
}