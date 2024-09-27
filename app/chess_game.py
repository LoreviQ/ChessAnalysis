class Game:
    def __init__(self):
        self.board = Board()
        self.pieces = []
        self._initialize_pieces()
        self.turn = "white"

    def _initialize_pieces(self):
        for file in self.board.files:
            self.pieces += [self.board.squares[f"{file}2"].add_piece(Pawn("white"))]
            self.pieces += [self.board.squares[f"{file}7"].add_piece(Pawn("black"))]

    def _get_direction(self, colour):
        return 1 if colour == "white" else -1

    def make_move(self, move):
        # Ignore pieces that are not active or not of the current turn
        possible_pieces = [
            piece for piece in self.pieces if piece.colour == self.turn and piece.active
        ]
        match move[0]:
            case _ if move[0] in "abcdefgh":  # pawn move
                possible_pieces = [
                    piece for piece in possible_pieces if isinstance(piece, Pawn)
                ]
            case "R":  # rook move
                pass
            case "N":  # knight move
                pass
            case "B":  # bishop move
                pass
            case "Q":  # queen move
                pass
            case "K":  # king move
                pass
            case "O":  # castling
                pass
            case _:
                print("Invalid move.")
        # Find the piece that can make the move
        for piece in possible_pieces:
            if move in piece.list_possible_moves():
                self._move_piece(piece, move)
                break
        print("Invalid move.")

    def _move_piece(self, piece, move):
        piece.position.remove_piece()
        move_square = self.board.get_square(move[-2:])
        if move_square.piece:
            move_square.piece.taken()
        move_square.add_piece(piece)
        self.turn = "black" if self.turn == "white" else "white"

    def start_game(self):
        while True:
            self.board.print_board()
            move = input(f"{self.turn.capitalize()}'s move: ")
            if move.lower() == "quit":
                print("Game over.")
                break
            self.make_move(move)


class Board:
    def __init__(self):
        self.ranks = "12345678"
        self.files = "abcdefgh"
        self.squares = self._initialize_board()

    def _initialize_board(self):
        board = {}
        for rank in self.ranks:
            for file in self.files:
                board[f"{file}{rank}"] = Square(rank, file, self)
        return board

    def print_board(self):
        for rank in self.ranks[::-1]:
            row = ""
            for file in self.files:
                row += (
                    self.squares[f"{file}{rank}"].piece.printable
                    if self.squares[f"{file}{rank}"].piece
                    and self.squares[f"{file}{rank}"].piece.active
                    else " "
                )
            print(row)

    def get_square(self, algebraic_notation):
        try:
            return self.squares[algebraic_notation]
        except KeyError:
            return None


class Square:
    def __init__(self, rank, file, board):
        self.rank = rank
        self.file = file
        self.piece = None
        self.board = board

    def add_piece(self, piece):
        self.piece = piece
        self.piece.position = self
        return piece

    def remove_piece(self):
        self.piece = None


class Piece:
    def __init__(self, colour):
        self.active = True
        self.colour = colour
        self.position = None

    def taken(self):
        self.active = False
        self.position = None


class Pawn(Piece):
    def __init__(self, colour):
        super().__init__(colour)
        self.printable = "♟" if colour == "white" else "♙"
        self.direction = 1 if colour == "white" else -1

    def list_possible_moves(self):
        if self.active is False:
            return None
        moves = []
        current_square = self.position
        # forward 1 square
        forward_square = self.position.board.get_square(
            f"{current_square.file}{int(current_square.rank) + self.direction}"
        )
        if forward_square and forward_square.piece is None:
            moves.append(f"{forward_square.file}{forward_square.rank}")
            # forward 2 squares from starting position
            if (self.position.rank == "2" and self.direction == 1) or (
                self.position.rank == "7" and self.direction == -1
            ):
                double_forward_square = self.position.board.get_square(
                    f"{current_square.file}{int(current_square.rank) + 2*self.direction}"
                )
                if double_forward_square and double_forward_square.piece is None:
                    moves.append(
                        f"{double_forward_square.file}{double_forward_square.rank}"
                    )
        # capture diagonally
        left_diagonal_square = self.position.board.get_square(
            f"{chr(ord(current_square.file) - 1)}{int(current_square.rank) + self.direction}"
        )
        right_diagonal_square = self.position.board.get_square(
            f"{chr(ord(current_square.file) + 1)}{int(current_square.rank) + self.direction}"
        )
        if left_diagonal_square and left_diagonal_square.piece:
            moves.append(
                f"{current_square.file}x{left_diagonal_square.file}{left_diagonal_square.rank}"
            )
        if right_diagonal_square and right_diagonal_square.piece:
            moves.append(
                f"{current_square.file}x{right_diagonal_square.file}{right_diagonal_square.rank}"
            )
        return moves


if __name__ == "__main__":
    chessGame = Game()
    chessGame.start_game()
