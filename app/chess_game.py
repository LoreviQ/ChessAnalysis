import re


class Game:
    def __init__(self):
        self.board = Board()
        self.pieces = []
        self._initialize_pieces()
        self.turn = "white"
        self.previous_move = None

    def _initialize_pieces(self):
        for file in self.board.files:
            self.pieces += [self.board.squares[f"{file}2"].add_piece(Pawn("white"))]
            self.pieces += [self.board.squares[f"{file}7"].add_piece(Pawn("black"))]

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
            if move in piece.list_possible_moves(self.previous_move):
                self._move_piece(piece, move)
                return
        print("Invalid move.")

    def _move_piece(self, piece, move):
        destination = re.search(r"([a-h][1-8])(=[QRBN])?([+#]?)", move)
        destination_square = self.board.get_square(
            destination.group(1)[0], destination.group(1)[1]
        )
        origin = piece.position.string
        take = "x" if "x" in move else ""

        long_notation = f"{piece.string}{origin}{take}{destination.group(0)}"
        piece.position.remove_piece()
        if take:
            if destination_square.piece:
                destination_square.piece.taken()
            else:
                # en passant
                taken_square = self.board.get_square(
                    self.previous_move[-2:][0], self.previous_move[-2:][1]
                )
                taken_square.piece.taken()
        destination_square.add_piece(piece)
        self.turn = "black" if self.turn == "white" else "white"
        self.previous_move = long_notation
        print(long_notation)

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

    def get_square(self, file, rank):
        try:
            return self.squares[f"{file}{rank}"]
        except KeyError:
            return None


class Square:
    def __init__(self, rank, file, board):
        self.rank = rank
        self.file = file
        self.piece = None
        self.board = board
        self.string = f"{file}{rank}"

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
        self.string = ""
        self.direction = 1 if colour == "white" else -1

    def list_possible_moves(self, previous_move):
        if self.active is False:
            return None
        moves = []
        origin_file = self.position.file
        origin_rank = self.position.rank

        # forward 1 square
        new_file = origin_file
        new_rank = int(origin_rank) + self.direction
        forward_square = self.position.board.get_square(new_file, new_rank)
        if forward_square and forward_square.piece is None:
            move = forward_square.string
            # promotion
            if new_rank in [1, 8]:
                moves.append(f"{move}=Q")
                moves.append(f"{move}=R")
                moves.append(f"{move}=B")
                moves.append(f"{move}=N")
            else:
                moves.append(move)
            # forward 2 squares from starting position
            if (origin_rank == "2" and self.direction == 1) or (
                origin_rank == "7" and self.direction == -1
            ):
                new_rank = int(origin_rank) + 2 * self.direction
                double_forward_square = self.position.board.get_square(
                    new_file, new_rank
                )
                if double_forward_square and double_forward_square.piece is None:
                    moves.append(double_forward_square.string)

        # capture diagonally
        new_file = chr(ord(origin_file) - 1)
        new_rank = int(origin_rank) + self.direction
        left_diagonal_square = self.position.board.get_square(new_file, new_rank)
        new_file = chr(ord(origin_file) + 1)
        right_diagonal_square = self.position.board.get_square(new_file, new_rank)
        if left_diagonal_square and left_diagonal_square.piece:
            moves.append(f"{origin_file}x{left_diagonal_square.string}")
        if right_diagonal_square and right_diagonal_square.piece:
            moves.append(f"{origin_file}x{right_diagonal_square.string}")

        # en passant
        if (self.position.rank == "5" and self.direction == 1) or (
            self.position.rank == "4" and self.direction == -1
        ):
            en_passant_left = f"{left_diagonal_square.file}{int(left_diagonal_square.rank)+self.direction}{left_diagonal_square.file}{int(left_diagonal_square.rank)-self.direction}"
            if previous_move == en_passant_left:
                moves.append(f"{origin_file}x{left_diagonal_square.string}")
            en_passant_right = f"{right_diagonal_square.file}{int(right_diagonal_square.rank)+self.direction}{right_diagonal_square.file}{int(right_diagonal_square.rank)-self.direction}"
            if previous_move == en_passant_right:
                moves.append(f"{origin_file}x{right_diagonal_square.string}")
        return moves


if __name__ == "__main__":
    chessGame = Game()
    chessGame.start_game()
