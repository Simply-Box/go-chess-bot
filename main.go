package main

import (
	"github.com/Simply-Box/go-chess-bot/board"
	//"chessEngine/movegen"
	"fmt"
    "bufio"
    "os"
    "strings"
)

func main() {
    fmt.Println("Starting Chess Engine...")
    gameState := board.ConvertFENtoGameState(board.GetStartingFEN())
    board.PrintBoard(gameState.Board)
    reader := bufio.NewReader(os.Stdin)

    for {
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

        allMoves := board.GenerateAllMoves(&gameState)

        err = board.MakeMove(&gameState, &allMoves, move)
        if err != nil {
            fmt.Println(err)
            continue
        }

        board.PrintBoard(gameState.Board)
	}
}