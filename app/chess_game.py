"""
Module to play chess in the terminal.
"""

import re


class Game:
    """
    A class to manage a chess game.
    """

    def __init__(self):
        self.board = Board()
        self.pieces = []
        self._initialize_pieces()
        self.turn = "white"
        self.previous_moves = []

    def _initialize_pieces(self):
        for file in self.board.files:
            self.pieces += [self.board.squares[f"{file}2"].add_piece(Pawn("white"))]
            self.pieces += [self.board.squares[f"{file}7"].add_piece(Pawn("black"))]
        self.pieces += [self.board.squares["a3"].add_piece(Rook("white"))]
        self.pieces += [self.board.squares["h3"].add_piece(Rook("white"))]

    def _make_move(self, move):
        """
        Checks if the move is valid and makes the move if it is.
        """
        # Parse move string
        try:
            piece_str, take, destination, promotion, checkmate = self._regex_match(
                move, log=True
            )
        except ValueError:
            self._declare_invalid_move()
            return
        # Find the piece that can make the move and move it
        possible_moves = self._list_moves(self.turn, piece_str)
        for piece, possible_piece_moves in possible_moves.items():
            for possible_move in possible_piece_moves:
                if possible_move == move:
                    self._move_piece(
                        piece,
                        piece_str,
                        take,
                        destination,
                        promotion,
                        checkmate,
                    )
                    return
        # If no piece can make the move, declare it invalid
        self._declare_invalid_move()

    def _move_piece(self, piece, piece_str, take, destination, promotion, checkmate):
        destination_square = self.board.get_square(destination[0], destination[1])
        origin = piece.position.string
        long_notation = f"{piece_str}{origin}{take}{destination}{promotion}{checkmate}"

        piece.position.remove_piece()  # Remove piece from origin square
        if take:
            if destination_square.piece:
                destination_square.piece.taken()
            else:
                # en passant
                taken_square = self.board.get_square(
                    self.previous_moves[-2:][0], self.previous_moves[-2:][1]
                )
                taken_square.piece.taken()
        destination_square.add_piece(piece)
        # Promotion
        if promotion:
            piece = piece.promote(promotion[1])
            self.pieces += [destination_square.add_piece(piece)]
        self.turn = "black" if self.turn == "white" else "white"
        self.previous_moves += [long_notation]

    def _list_moves(self, colour=None, piece_str=None):
        """
        Lists all possible moves for the current player.
        """
        if colour is None:
            colour = self.turn
        possible_moves = {}
        for piece in self.pieces:
            if (
                piece.active
                and piece.colour == colour
                and (piece_str is None or piece.string == piece_str)
            ):
                possible_moves[piece] = piece.list_possible_moves(self)
        return possible_moves

    def _regex_match(self, move, log=False):
        # regex to match algebraic notation
        pattern = r"([KQRBN]?)(x?)([a-h][1-8])(=[QRBN])?([+#]?)"
        regex_match = re.search(pattern, move)
        if not regex_match:
            raise ValueError("Invalid move.")
        piece_str = regex_match.group(1)
        take = regex_match.group(2)
        destination = regex_match.group(3)
        promotion = regex_match.group(4)
        checkmate = regex_match.group(5)
        if log:
            print(
                f"piece: {piece_str}, take: {take}, destination: {destination}, promotion: {promotion}, checkmate: {checkmate}"
            )
        return (
            piece_str,
            take,
            destination,
            promotion if promotion else "",
            checkmate,
        )

    def _declare_invalid_move(self):
        print("Invalid move.")
        self._list_moves()

    def start_game(self):
        """
        Starts the game and allows players to make moves.
        Expects mvoe input in algebraic notation.
        """
        while True:
            self.board.print_board()
            user_input = input(f"{self.turn.capitalize()}'s move: ")
            if user_input.lower() == "quit":
                print("Game over.")
                break
            if user_input.lower() == "moves":
                print(self.previous_moves)
                continue
            if user_input.lower() == "list_moves":
                possible_moves = self._list_moves()
                print(f"{self.turn.capitalize()}'s possible moves:")
                for piece, moves in possible_moves.items():
                    print(
                        f"{piece.printable}-{piece.position.string}: {', '.join(moves)}"
                    )
                continue
            self._make_move(user_input)


class Board:
    """
    A class to represent a chess board.
    """

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
        """
        Prints the board with the pieces in their current positions.
        """
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
        """
        Returns the square at the given file and rank.
        """
        try:
            return self.squares[f"{file}{rank}"]
        except KeyError:
            return None


class Square:
    """
    A class to represent a square on a chess board.
    """

    def __init__(self, rank, file, board):
        self.rank = rank
        self.file = file
        self.piece = None
        self.board = board
        self.string = f"{file}{rank}"

    def add_piece(self, piece):
        """
        Adds a piece to the square.
        """
        self.piece = piece
        self.piece.position = self
        return piece

    def remove_piece(self):
        """
        Removes a piece from the square.
        """
        self.piece = None


class Piece:
    """
    A class to represent a chess piece.
    """

    def __init__(self, colour):
        self.active = True
        self.colour = colour
        self.position = None

    def taken(self):
        """
        Sets the piece as taken.
        """
        self.active = False
        self.position = None

    def list_possible_moves(self, _):
        """
        Returns a list of possible moves for the piece.
        Not implemented in the base class.
        """
        return []


class Pawn(Piece):
    """
    A class to represent a pawn chess piece.
    """

    def __init__(self, colour):
        super().__init__(colour)
        self.printable = "♟" if colour == "white" else "♙"
        self.string = ""
        self.direction = 1 if colour == "white" else -1

    def list_possible_moves(self, game):
        """
        Returns a list of possible moves for the pawn.
        """
        if self.active is False:
            return []
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
        if (
            left_diagonal_square
            and left_diagonal_square.piece
            and left_diagonal_square.piece.colour != self.colour
        ):
            moves.append(f"{origin_file}x{left_diagonal_square.string}")
        if (
            right_diagonal_square
            and right_diagonal_square.piece
            and right_diagonal_square.piece.colour != self.colour
        ):
            moves.append(f"{origin_file}x{right_diagonal_square.string}")
            # TODO capture promotion

        # en passant
        if (self.position.rank == "5" and self.direction == 1) or (
            self.position.rank == "4" and self.direction == -1
        ):
            en_passant_left = f"{left_diagonal_square.file}{int(left_diagonal_square.rank)+self.direction}{left_diagonal_square.file}{int(left_diagonal_square.rank)-self.direction}"
            if game.previous_moves[-1] == en_passant_left:
                moves.append(f"{origin_file}x{left_diagonal_square.string}")
            en_passant_right = f"{right_diagonal_square.file}{int(right_diagonal_square.rank)+self.direction}{right_diagonal_square.file}{int(right_diagonal_square.rank)-self.direction}"
            if game.previous_moves[-1] == en_passant_right:
                moves.append(f"{origin_file}x{right_diagonal_square.string}")

        return moves

    def promote(self, promotion):
        """
        Promotes the pawn to a queen, rook, bishop, or knight.
        """
        self.taken()
        if promotion == "Q":
            return Queen(self.colour)
        if promotion == "R":
            return Rook(self.colour)
        if promotion == "B":
            return Bishop(self.colour)
        if promotion == "N":
            return Knight(self.colour)
        raise ValueError("Invalid promotion.")


class Rook(Piece):
    """
    A class to represent a rook chess piece.
    """

    def __init__(self, colour):
        super().__init__(colour)
        self.printable = "♜" if colour == "white" else "♖"
        self.string = "R"

    def list_possible_moves(self, _):
        """
        Returns a list of possible moves for the rook.
        """
        if self.active is False:
            return []
        moves = []
        origin_file = self.position.file
        origin_rank = self.position.rank

        # forward
        for rank in range(int(origin_rank) + 1, 9):
            square = self.position.board.get_square(origin_file, str(rank))
            if square.piece and square.piece.colour != self.colour:
                moves.append(f"{self.string}x{square.string}")
                break
            moves.append(f"{self.string}{square.string}")

        # backward
        for rank in range(int(origin_rank) - 1, 0, -1):
            square = self.position.board.get_square(origin_file, str(rank))
            if square.piece and square.piece.colour != self.colour:
                moves.append(f"{self.string}x{square.string}")
                break
            moves.append(f"{self.string}{square.string}")

        # left
        for file in range(ord(origin_file) - 1, ord("a") - 1, -1):
            square = self.position.board.get_square(chr(file), origin_rank)
            if square.piece and square.piece.colour != self.colour:
                moves.append(f"{self.string}x{square.string}")
                break
            moves.append(f"{self.string}{square.string}")

        # right
        for file in range(ord(origin_file) + 1, ord("h") + 1):
            square = self.position.board.get_square(chr(file), origin_rank)
            if square.piece and square.piece.colour != self.colour:
                moves.append(f"{self.string}x{square.string}")
                break
            moves.append(f"{self.string}{square.string}")
        return moves


class Knight(Piece):
    """
    A class to represent a knight chess piece.
    """

    def __init__(self, colour):
        super().__init__(colour)
        self.printable = "♞" if colour == "white" else "♘"
        self.string = "N"


class Bishop(Piece):
    """
    A class to represent a bishop chess piece.
    """

    def __init__(self, colour):
        super().__init__(colour)
        self.printable = "♝" if colour == "white" else "♗"
        self.string = "B"


class Queen(Piece):
    """
    A class to represent a queen chess piece.
    """

    def __init__(self, colour):
        super().__init__(colour)
        self.printable = "♛" if colour == "white" else "♕"
        self.string = "Q"


class King(Piece):
    """
    A class to represent a king chess piece.
    """

    def __init__(self, colour):
        super().__init__(colour)
        self.printable = "♚" if colour == "white" else "♔"
        self.string = "K"


if __name__ == "__main__":
    chessGame = Game()
    chessGame.start_game()
