class Game:
    def __init__(self):
        self.board = Board()
        self.pieces = []
        self._initialize_pieces()

    def _initialize_pieces(self):
        for file in self.board.files:
            self.pieces += [self.board.squares[f"{file}2"].add_piece(Pawn("white"))]
            self.pieces += [self.board.squares[f"{file}7"].add_piece(Pawn("black"))]

    def make_move(self, move):
        print(move)

    def start_game(self):
        while True:
            self.board.print_board()
            move = input("Enter your move: ")
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
                board[f"{file}{rank}"] = Square(rank, file)
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
        return self.squares[algebraic_notation]


class Square:
    def __init__(self, rank, file):
        self.rank = rank
        self.file = file
        self.piece = None

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
        self.printable = "♟" if colour == "black" else "♙"

    def list_moves(self):
        pass


if __name__ == "__main__":
    chessGame = Game()
    chessGame.start_game()
