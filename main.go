package main

import (
	"fmt"
	"strings"
)

// Piece represents a chess piece with a type and color
type Piece struct {
	Type  string
	Color string
}

// Board represents the chess board
type Board [8][8]*Piece

// GameState stores the state of the game, including castling rights and en passant targets
type GameState struct {
	Board                   Board
	WhiteCanCastleKingSide  bool
	WhiteCanCastleQueenSide bool
	BlackCanCastleKingSide  bool
	BlackCanCastleQueenSide bool
	EnPassantTarget         [2]int
}

// Initialize the board with pieces in their starting positions
func (gs *GameState) Initialize() {
	// Initialize white pieces
	gs.Board[0][0] = &Piece{Type: "R", Color: "W"}
	gs.Board[0][1] = &Piece{Type: "N", Color: "W"}
	gs.Board[0][2] = &Piece{Type: "B", Color: "W"}
	gs.Board[0][3] = &Piece{Type: "Q", Color: "W"}
	gs.Board[0][4] = &Piece{Type: "K", Color: "W"}
	gs.Board[0][5] = &Piece{Type: "B", Color: "W"}
	gs.Board[0][6] = &Piece{Type: "N", Color: "W"}
	gs.Board[0][7] = &Piece{Type: "R", Color: "W"}
	for i := 0; i < 8; i++ {
		gs.Board[1][i] = &Piece{Type: "P", Color: "W"}
	}

	// Initialize black pieces
	gs.Board[7][0] = &Piece{Type: "R", Color: "B"}
	gs.Board[7][1] = &Piece{Type: "N", Color: "B"}
	gs.Board[7][2] = &Piece{Type: "B", Color: "B"}
	gs.Board[7][3] = &Piece{Type: "Q", Color: "B"}
	gs.Board[7][4] = &Piece{Type: "K", Color: "B"}
	gs.Board[7][5] = &Piece{Type: "B", Color: "B"}
	gs.Board[7][6] = &Piece{Type: "N", Color: "B"}
	gs.Board[7][7] = &Piece{Type: "R", Color: "B"}
	for i := 0; i < 8; i++ {
		gs.Board[6][i] = &Piece{Type: "P", Color: "B"}
	}

	gs.WhiteCanCastleKingSide = true
	gs.WhiteCanCastleQueenSide = true
	gs.BlackCanCastleKingSide = true
	gs.BlackCanCastleQueenSide = true
	gs.EnPassantTarget = [2]int{-1, -1}
}

// Display the board
func (gs *GameState) Display() {
	for i := 7; i >= 0; i-- {
		for j := 0; j < 8; j++ {
			if gs.Board[i][j] == nil {
				fmt.Print(". ")
			} else {
				fmt.Print(gs.Board[i][j].Type + gs.Board[i][j].Color + " ")
			}
		}
		fmt.Println()
	}
}

// MovePiece moves a piece from one position to another if the move is valid
func (gs *GameState) MovePiece(fromX, fromY, toX, toY int) bool {
	piece := gs.Board[fromX][fromY]
	if piece == nil {
		return false
	}

	target := gs.Board[toX][toY]
	if target != nil && target.Color == piece.Color {
		return false
	}

	switch piece.Type {
	case "P":
		if !validPawnMove(gs, fromX, fromY, toX, toY, piece.Color) {
			return false
		}
	case "R":
		if !validRookMove(gs.Board, fromX, fromY, toX, toY) {
			return false
		}
	case "N":
		if !validKnightMove(fromX, fromY, toX, toY) {
			return false
		}
	case "B":
		if !validBishopMove(gs.Board, fromX, fromY, toX, toY) {
			return false
		}
	case "Q":
		if !validQueenMove(gs.Board, fromX, fromY, toX, toY) {
			return false
		}
	case "K":
		if !validKingMove(gs, fromX, fromY, toX, toY) {
			return false
		}
	}

	if piece.Type == "P" && (toX == 0 || toX == 7) {
		// Pawn promotion
		gs.Board[toX][toY] = &Piece{Type: "Q", Color: piece.Color}
	} else {
		gs.Board[toX][toY] = piece
	}

	gs.Board[fromX][fromY] = nil

	// Handle castling rights and en passant
	if piece.Type == "K" {
		if piece.Color == "W" {
			gs.WhiteCanCastleKingSide = false
			gs.WhiteCanCastleQueenSide = false
		} else {
			gs.BlackCanCastleKingSide = false
			gs.BlackCanCastleQueenSide = false
		}
	}
	if piece.Type == "R" {
		if fromX == 0 && fromY == 0 {
			gs.WhiteCanCastleQueenSide = false
		} else if fromX == 0 && fromY == 7 {
			gs.WhiteCanCastleKingSide = false
		} else if fromX == 7 && fromY == 0 {
			gs.BlackCanCastleQueenSide = false
		} else if fromX == 7 && fromY == 7 {
			gs.BlackCanCastleKingSide = false
		}
	}

	if piece.Type == "P" && abs(fromX-toX) == 2 {
		gs.EnPassantTarget = [2]int{(fromX + toX) / 2, fromY}
	} else {
		gs.EnPassantTarget = [2]int{-1, -1}
	}

	return true
}

// Check if a pawn move is valid
func validPawnMove(gs *GameState, fromX, fromY, toX, toY int, color string) bool {
	board := gs.Board
	if color == "W" {
		if fromX == 1 && toX == fromX+2 && fromY == toY && board[toX][toY] == nil && board[toX-1][toY] == nil {
			return true
		}
		if toX == fromX+1 && fromY == toY && board[toX][toY] == nil {
			return true
		}
		if toX == fromX+1 && (toY == fromY+1 || toY == fromY-1) && board[toX][toY] != nil && board[toX][toY].Color != color {
			return true
		}
		if toX == fromX+1 && toY == fromY+1 && gs.EnPassantTarget[0] == toX && gs.EnPassantTarget[1] == toY {
			board[toX-1][toY] = nil
			return true
		}
	} else {
		if fromX == 6 && toX == fromX-2 && fromY == toY && board[toX][toY] == nil && board[toX+1][toY] == nil {
			return true
		}
		if toX == fromX-1 && fromY == toY && board[toX][toY] == nil {
			return true
		}
		if toX == fromX-1 && (toY == fromY+1 || toY == fromY-1) && board[toX][toY] != nil && board[toX][toY].Color != color {
			return true
		}
		if toX == fromX-1 && toY == fromY-1 && gs.EnPassantTarget[0] == toX && gs.EnPassantTarget[1] == toY {
			board[toX+1][toY] = nil
			return true
		}
	}
	return false
}

// Check if a rook move is valid
func validRookMove(board Board, fromX, fromY, toX, toY int) bool {
	if fromX != toX && fromY != toY {
		return false
	}
	if fromX == toX {
		for y := min(fromY, toY) + 1; y < max(fromY, toY); y++ {
			if board[fromX][y] != nil {
				return false
			}
		}
	} else {
		for x := min(fromX, toX) + 1; x < max(fromX, toX); x++ {
			if board[x][fromY] != nil {
				return false
			}
		}
	}
	return true
}

// Check if a knight move is valid
func validKnightMove(fromX, fromY, toX, toY int) bool {
	return (abs(fromX-toX) == 2 && abs(fromY-toY) == 1) || (abs(fromX-toX) == 1 && abs(fromY-toY) == 2)
}

// Check if a bishop move is valid
func validBishopMove(board Board, fromX, fromY, toX, toY int) bool {
	if abs(fromX-toX) != abs(fromY-toY) {
		return false
	}
	xStep := 1
	if fromX > toX {
		xStep = -1
	}
	yStep := 1
	if fromY > toY {
		yStep = -1
	}
	x, y := fromX+xStep, fromY+yStep
	for x != toX && y != toY {
		if board[x][y] != nil {
			return false
		}
		x += xStep
		y += yStep
	}
	return true
}

// Check if a queen move is valid
func validQueenMove(board Board, fromX, fromY, toX, toY int) bool {
	return validRookMove(board, fromX, fromY, toX, toY) || validBishopMove(board, fromX, fromY, toX, toY)
}

// Check if a king move is valid, including castling
func validKingMove(gs *GameState, fromX, fromY, toX, toY int) bool {
	board := gs.Board
	if abs(fromX-toX) <= 1 && abs(fromY-toY) <= 1 {
		return true
	}

	if fromX == toX && abs(fromY-toY) == 2 {
		if board[fromX][fromY].Color == "W" {
			if fromY < toY && gs.WhiteCanCastleKingSide && board[fromX][fromY+1] == nil && board[fromX][fromY+2] == nil {
				board[fromX][fromY+1] = board[fromX][7]
				board[fromX][7] = nil
				return true
			}
			if fromY > toY && gs.WhiteCanCastleQueenSide && board[fromX][fromY-1] == nil && board[fromX][fromY-2] == nil && board[fromX][fromY-3] == nil {
				board[fromX][fromY-1] = board[fromX][0]
				board[fromX][0] = nil
				return true
			}
		} else {
			if fromY < toY && gs.BlackCanCastleKingSide && board[fromX][fromY+1] == nil && board[fromX][fromY+2] == nil {
				board[fromX][fromY+1] = board[fromX][7]
				board[fromX][7] = nil
				return true
			}
			if fromY > toY && gs.BlackCanCastleQueenSide && board[fromX][fromY-1] == nil && board[fromX][fromY-2] == nil && board[fromX][fromY-3] == nil {
				board[fromX][fromY-1] = board[fromX][0]
				board[fromX][0] = nil
				return true
			}
		}
	}
	return false
}

// Helper functions
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Check if the king of the given color is in check
func (gs *GameState) IsInCheck(color string) bool {
	var kingX, kingY int
	// Find the king
	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			if gs.Board[x][y] != nil && gs.Board[x][y].Type == "K" && gs.Board[x][y].Color == color {
				kingX, kingY = x, y
				break
			}
		}
	}

	// Check if any opposing piece can attack the king
	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			if gs.Board[x][y] != nil && gs.Board[x][y].Color != color {
				if gs.IsValidMove(x, y, kingX, kingY) {
					return true
				}
			}
		}
	}

	return false
}

// Check if the given move is valid according to the piece rules
func (gs *GameState) IsValidMove(fromX, fromY, toX, toY int) bool {
	piece := gs.Board[fromX][fromY]
	if piece == nil {
		return false
	}

	switch piece.Type {
	case "P":
		return validPawnMove(gs, fromX, fromY, toX, toY, piece.Color)
	case "R":
		return validRookMove(gs.Board, fromX, fromY, toX, toY)
	case "N":
		return validKnightMove(fromX, fromY, toX, toY)
	case "B":
		return validBishopMove(gs.Board, fromX, fromY, toX, toY)
	case "Q":
		return validQueenMove(gs.Board, fromX, fromY, toX, toY)
	case "K":
		return validKingMove(gs, fromX, fromY, toX, toY)
	}

	return false
}

// Check if the player is in checkmate
func (gs *GameState) IsCheckmate(color string) bool {
	// If not in check, cannot be checkmate
	if !gs.IsInCheck(color) {
		return false
	}

	// Check if there are any valid moves to escape check
	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			if gs.Board[x][y] != nil && gs.Board[x][y].Color == color {
				for toX := 0; toX < 8; toX++ {
					for toY := 0; toY < 8; toY++ {
						originalPiece := gs.Board[toX][toY]
						if gs.MovePiece(x, y, toX, toY) {
							if !gs.IsInCheck(color) {
								gs.Board[toX][toY] = originalPiece
								gs.Board[x][y] = gs.Board[toX][toY]
								return false
							}
							gs.Board[toX][toY] = originalPiece
							gs.Board[x][y] = gs.Board[toX][toY]
						}
					}
				}
			}
		}
	}

	return true
}

// ParseMove converts a move string (e.g., "e2,e4") to board coordinates
func ParseMove(move string) (int, int, int, int, bool) {
	parts := strings.Split(move, ",")
	if len(parts) != 2 || len(parts[0]) != 2 || len(parts[1]) != 2 {
		return 0, 0, 0, 0, false
	}

	fromX, fromY := int(parts[0][1]-'1'), int(parts[0][0]-'a')
	toX, toY := int(parts[1][1]-'1'), int(parts[1][0]-'a')

	if fromX < 0 || fromX > 7 || fromY < 0 || fromY > 7 || toX < 0 || toX > 7 || toY < 0 || toY > 7 {
		return 0, 0, 0, 0, false
	}

	return fromX, fromY, toX, toY, true
}

func main() {
	var gameState GameState
	gameState.Initialize()
	var currentPlayer = "W"

	for {
		gameState.Display()

		if gameState.IsCheckmate(currentPlayer) {
			fmt.Printf("%s is in checkmate. %s wins!\n", currentPlayer, oppositeColor(currentPlayer))
			break
		}

		if gameState.IsInCheck(currentPlayer) {
			fmt.Printf("%s is in check.\n", currentPlayer)
		}

		var move string
		fmt.Printf("%s's turn. Enter your move (e.g., e2,e4): ", currentPlayer)
		fmt.Scanln(&move)

		fromX, fromY, toX, toY, valid := ParseMove(move)
		if !valid {
			fmt.Println("Invalid move format. Try again.")
			continue
		}

		if !gameState.MovePiece(fromX, fromY, toX, toY) {
			fmt.Println("Invalid move. Try again.")
			continue
		}

		currentPlayer = oppositeColor(currentPlayer)
	}
}

func oppositeColor(color string) string {
	if color == "W" {
		return "B"
	}
	return "W"
}
