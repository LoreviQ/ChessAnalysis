"""
Module to play chess in the terminal.

Used to convert a list of moves to long algebraic notation.

It's necessary to play the moves in order to convert them 
since short algebraic notation relies on board state.
"""

import re


def convert_notation(moves, debug=False):
    """
    Converts a list of moves to long algebraic notation
    by playing them and returning the game log
    """
    chess_game = Game()
    return chess_game.convert_notation(moves, debug)


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
        self.can_castle = {
            "white": {"short": True, "long": True},
            "black": {"short": True, "long": True},
        }

    def _initialize_pieces(self):
        for file in self.board.files:
            self.pieces += [self.board.squares[f"{file}2"].add_piece(Pawn("white"))]
            self.pieces += [self.board.squares[f"{file}7"].add_piece(Pawn("black"))]
        self.pieces += [self.board.squares["a1"].add_piece(Rook("white"))]
        self.pieces += [self.board.squares["h1"].add_piece(Rook("white"))]
        self.pieces += [self.board.squares["a8"].add_piece(Rook("black"))]
        self.pieces += [self.board.squares["h8"].add_piece(Rook("black"))]
        self.pieces += [self.board.squares["b1"].add_piece(Knight("white"))]
        self.pieces += [self.board.squares["g1"].add_piece(Knight("white"))]
        self.pieces += [self.board.squares["b8"].add_piece(Knight("black"))]
        self.pieces += [self.board.squares["g8"].add_piece(Knight("black"))]
        self.pieces += [self.board.squares["c1"].add_piece(Bishop("white"))]
        self.pieces += [self.board.squares["f1"].add_piece(Bishop("white"))]
        self.pieces += [self.board.squares["c8"].add_piece(Bishop("black"))]
        self.pieces += [self.board.squares["f8"].add_piece(Bishop("black"))]
        self.pieces += [self.board.squares["d1"].add_piece(Queen("white"))]
        self.pieces += [self.board.squares["d8"].add_piece(Queen("black"))]
        self.pieces += [self.board.squares["e1"].add_piece(King("white"))]
        self.pieces += [self.board.squares["e8"].add_piece(King("black"))]

    def _make_move(self, move):
        """
        Checks if the move is valid and makes the move if it is.
        :param move: A string representing the move in short algebraic notation.
        (e.g. "e4", "Nf3", "Bxb7", "0-0", "0-0-0")
        return: True if the move was made, False otherwise.
        """
        # Parse move string
        try:
            piece_str, take, destination, promotion, checkmate = self._regex_match(move)
        except ValueError:
            self._declare_invalid_move()
            return False
        if piece_str in ("0-0", "O-O"):
            return self._short_castle()
        if piece_str in ("0-0-0", "O-O-O"):
            return self._long_castle()
        # check/checkmate ignored for move finding purposes
        if checkmate:
            move = f"{move[:-1]}"
        # Find the piece that can make the move and move it
        possible_moves = self._list_moves(self.turn, piece_str)
        for possible_move, piece in possible_moves.items():
            if possible_move == move:
                return self._move_piece(
                    piece,
                    piece_str,
                    take,
                    destination,
                    promotion,
                    checkmate,
                )
        # If no piece can make the move, declare it invalid
        self._declare_invalid_move()
        return False

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

        # disable castling if king or rook moves
        if piece_str == "K":
            self.can_castle[self.turn]["short"] = False
            self.can_castle[self.turn]["long"] = False
        if piece_str == "R":
            if origin[0] == "a":
                self.can_castle[self.turn]["long"] = False
            if origin[0] == "h":
                self.can_castle[self.turn]["short"] = False
        self._change_turn()

        # log move
        self.previous_moves += [long_notation]
        return True

    def _validate_castle_conditions(self, colour, side):
        # Check if castling has been disabled for this player
        if self.can_castle[colour][side] is False:
            return False
        # Check if king and rook are in their home positions
        home_rank = "1" if colour == "white" else "8"
        king = self.board.get_square("e", home_rank).piece
        rook_file = "h" if side == "short" else "a"
        rook = self.board.get_square(rook_file, home_rank).piece
        if (
            isinstance(king, King) is False
            or isinstance(rook, Rook) is False
            or king.active is False
            or rook.active is False
        ):
            return False
        # Check if squares between king and rook are empty
        inbetween_files = ["f", "g"] if side == "short" else ["b", "c", "d"]
        for file in inbetween_files:
            if self.board.get_square(file, home_rank).piece:
                return False
        # Check if king passes through check
        # TODO
        return True

    def _short_castle(self):
        if not self._validate_castle_conditions(self.turn, "short"):
            self._declare_invalid_move()
            return False
        # Make the move
        home_rank = "1" if self.turn == "white" else "8"
        king = self.board.get_square("e", home_rank).piece
        rook = self.board.get_square("h", home_rank).piece
        f_square = self.board.get_square("f", home_rank)
        g_square = self.board.get_square("g", home_rank)
        king.position.remove_piece()
        rook.position.remove_piece()
        f_square.add_piece(rook)
        g_square.add_piece(king)
        # Disable castling for the rest of the game
        self.can_castle[self.turn]["short"] = False
        self.can_castle[self.turn]["long"] = False
        self._change_turn()
        # Log move
        self.previous_moves += ["O-O"]
        return True

    def _long_castle(self):
        if not self._validate_castle_conditions(self.turn, "long"):
            self._declare_invalid_move()
            return False
        # Make the move
        home_rank = "1" if self.turn == "white" else "8"
        king = self.board.get_square("e", home_rank).piece
        rook = self.board.get_square("a", home_rank).piece
        d_square = self.board.get_square("d", home_rank)
        c_square = self.board.get_square("c", home_rank)
        king.position.remove_piece()
        rook.position.remove_piece()
        d_square.add_piece(rook)
        c_square.add_piece(king)
        # Disable castling for the rest of the game
        self.can_castle[self.turn]["short"] = False
        self.can_castle[self.turn]["long"] = False
        self._change_turn()
        # Log move
        self.previous_moves += ["O-O-O"]
        return True

    def _list_moves(self, colour=None, filter_piece_str=None):
        """
        Lists all possible moves for a player and piece type.
        :param colour: The colour of the player.
        :param filter_piece_str: The type of piece to filter by.
        (e.g. "N", "B", "R", "Q", "K", "" for pawn)
        """
        if colour is None:
            colour = self.turn
        possible_moves = {}
        for piece in self.pieces:
            # Filter out pieces that can't move
            if not (
                piece.active
                and piece.colour == colour
                and ((filter_piece_str is None) or (piece.string == filter_piece_str))
            ):
                continue
            piece_moves = piece.list_possible_moves(self)
            for move in piece_moves:
                if move not in possible_moves:
                    possible_moves[move] = piece
                    continue

                # Multiple pieces that can move to the same square
                piece_str, take, destination, promotion, checkmate = self._regex_match(
                    move
                )
                other = possible_moves[move]
                if other.position.file == piece.position.file:
                    # If the pieces are on the same file, disambiguate by rank
                    piece_move = (
                        f"{piece_str}{piece.position.rank}{take}{destination}"
                        f"{promotion}{checkmate}"
                    )
                    other_move = (
                        f"{piece_str}{other.position.rank}{take}{destination}"
                        f"{promotion}{checkmate}"
                    )
                else:
                    # Otherwise, disambiguate by file
                    piece_move = (
                        f"{piece_str}{piece.position.file}{take}{destination}"
                        f"{promotion}{checkmate}"
                    )
                    other_move = (
                        f"{piece_str}{other.position.file}{take}{destination}"
                        f"{promotion}{checkmate}"
                    )
                # TODO disambiguate by both file and rank (incredibly rare)
                possible_moves[piece_move] = piece
                possible_moves[other_move] = other
                del possible_moves[move]
        # Castling
        if self._validate_castle_conditions(self.turn, "short"):
            possible_moves["O-O"] = None
        if self._validate_castle_conditions(self.turn, "long"):
            possible_moves["O-O-O"] = None
        return possible_moves

    def _log_possible_moves(self, possible_moves):
        print(f"{self.turn.capitalize()}'s possible moves:")
        log_moves = {}
        for move, piece in possible_moves.items():
            if piece is None:
                print(move)
                continue
            piece_unique_print = f"{piece.printable}-{piece.position.string}"
            if piece_unique_print in log_moves:
                log_moves[piece_unique_print].append(move)
            else:
                log_moves[piece_unique_print] = [move]
        for piece_unique_print, moves in log_moves.items():
            print(f"{piece_unique_print}: {', '.join(moves)}")

    def _regex_match(self, move, log=False):
        # regex to match algebraic notation
        pattern = (
            r"^([NBRQK])?([a-h])?([1-8])?(x)?([a-h][1-8])(=[NBRQK])?(\+|#)?$|^O-O(-O)?$"
        )
        regex_match = re.search(pattern, move)
        if not regex_match:
            raise ValueError("Invalid move.")
        if move in ("0-0", "O-O", "0-0-0", "O-O-O"):
            return move, "", "", "", ""
        piece_str = regex_match.group(1)
        origin_file = regex_match.group(2)
        origin_rank = regex_match.group(3)
        take = regex_match.group(4)
        destination = regex_match.group(5)
        promotion = regex_match.group(6)
        checkmate = regex_match.group(7)
        if log:
            print(
                f"piece: {piece_str}, origin: {origin_file}{origin_rank}, "
                f"take: {take}, destination: {destination}, promotion: {promotion}, "
                f"checkmate: {checkmate}"
            )
        return (
            piece_str if piece_str else "",
            take if take else "",
            destination,
            promotion if promotion else "",
            checkmate if checkmate else "",
        )

    def _declare_invalid_move(self):
        print("Invalid move.")
        possible_moves = self._list_moves()
        self._log_possible_moves(possible_moves)

    def _new_game(self):
        self.board = Board()
        self.pieces = []
        self._initialize_pieces()
        self.turn = "white"
        self.previous_moves = []
        self.can_castle = {
            "white": {"short": True, "long": True},
            "black": {"short": True, "long": True},
        }

    def _change_turn(self):
        self.turn = "black" if self.turn == "white" else "white"

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
                self._log_possible_moves(possible_moves)
                continue
            self._make_move(user_input)

    def convert_notation(self, moves, debug=False):
        """
        Converts a list of moves to long algebraic notation
        by playing them and returning the game log
        """
        self._new_game()
        if debug:
            for move in moves:
                while True:
                    # wait for enter key to be pressed
                    user_input = input(f"{self.turn.capitalize()}'s move: {move} ")
                    if user_input.lower() == "print":
                        self.board.print_board()
                        continue
                    if self._make_move(move):
                        break
                    self.board.print_board()
        else:
            for move in moves:
                self._make_move(move)
        return self.previous_moves


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
        self.string = None

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

    def _move_orthogonally(self, moves, origin_file, origin_rank, direction):
        directions = {
            "right": (1, 0),
            "left": (-1, 0),
            "forward": (0, 1),
            "backward": (0, -1),
        }
        file_step, rank_step = directions[direction]
        file = ord(origin_file)
        rank = int(origin_rank)

        while True:
            file += file_step
            rank += rank_step

            if file < ord("a") or file > ord("h") or rank < 1 or rank > 8:
                break

            square = self.position.board.get_square(chr(file), str(rank))
            if square.piece:
                if square.piece.colour != self.colour:
                    moves.append(f"{self.string}x{square.string}")
                break
            moves.append(f"{self.string}{square.string}")

    def _move_diagonally(self, moves, origin_file, origin_rank, direction):
        directions = {
            "top_right": (1, 1),
            "top_left": (-1, 1),
            "bottom_right": (1, -1),
            "bottom_left": (-1, -1),
        }
        file_step, rank_step = directions[direction]
        file = ord(origin_file)
        rank = int(origin_rank)

        while True:
            file += file_step
            rank += rank_step

            if file < ord("a") or file > ord("h") or rank < 1 or rank > 8:
                break

            square = self.position.board.get_square(chr(file), str(rank))
            if square.piece:
                if square.piece.colour != self.colour:
                    moves.append(f"{self.string}x{square.string}")
                break
            moves.append(f"{self.string}{square.string}")


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

        self._forward(moves, origin_file, origin_rank)
        self._diagonal_capture(moves, origin_file, origin_rank, game)
        return moves

    def _forward(self, moves, origin_file, origin_rank):
        """
        Returns the moves obtainable by moving the pawn forward.
        """
        # forward 1 square
        new_rank = int(origin_rank) + self.direction
        forward_square = self.position.board.get_square(origin_file, new_rank)
        if forward_square and forward_square.piece is None:
            move = forward_square.string
            self._check_promotion(moves, move, new_rank)

            # forward 2 squares from starting position
            if (origin_rank == "2" and self.direction == 1) or (
                origin_rank == "7" and self.direction == -1
            ):
                new_rank = int(origin_rank) + 2 * self.direction
                double_forward_square = self.position.board.get_square(
                    origin_file, new_rank
                )
                if double_forward_square and double_forward_square.piece is None:
                    moves.append(double_forward_square.string)

    def _diagonal_capture(self, moves, origin_file, origin_rank, game):
        """
        Returns the moves obtainable by capturing diagonally.
        """
        # capture diagonally
        new_rank = int(origin_rank) + self.direction
        for side in [-1, 1]:
            new_file = chr(ord(origin_file) + side)
            diagonal_square = self.position.board.get_square(new_file, new_rank)
            if (
                diagonal_square
                and diagonal_square.piece
                and diagonal_square.piece.colour != self.colour
            ):
                move = f"{origin_file}x{diagonal_square.string}"
                self._check_promotion(moves, move, new_rank)

            # en passant
            if (self.position.rank == "5" and self.direction == 1) or (
                self.position.rank == "4" and self.direction == -1
            ):
                if diagonal_square:
                    en_passant = (
                        f"{new_file}{int(diagonal_square.rank)+self.direction}"
                        f"{new_file}{int(diagonal_square.rank)-self.direction}"
                    )
                    if game.previous_moves[-1] == en_passant:
                        moves.append(f"{origin_file}x{diagonal_square.string}")

    def _check_promotion(self, moves, move, rank):
        """
        Returns a list of possible promotions for the pawn or
        a list of the move if no promotion is possible.
        """
        if rank in [1, 8]:
            moves.append(f"{move}=Q")
            moves.append(f"{move}=R")
            moves.append(f"{move}=B")
            moves.append(f"{move}=N")
        else:
            moves.append(move)

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

        self._move_orthogonally(moves, origin_file, origin_rank, "forward")
        self._move_orthogonally(moves, origin_file, origin_rank, "backward")
        self._move_orthogonally(moves, origin_file, origin_rank, "left")
        self._move_orthogonally(moves, origin_file, origin_rank, "right")
        return moves


class Knight(Piece):
    """
    A class to represent a knight chess piece.
    """

    def __init__(self, colour):
        super().__init__(colour)
        self.printable = "♞" if colour == "white" else "♘"
        self.string = "N"

    def list_possible_moves(self, _):
        """
        Returns a list of possible moves for the knight.
        """
        if self.active is False:
            return []
        moves = []
        origin_file = self.position.file
        origin_rank = self.position.rank

        possible_moves = [
            (ord(origin_file) - 1, int(origin_rank) + 2),
            (ord(origin_file) + 1, int(origin_rank) + 2),
            (ord(origin_file) - 2, int(origin_rank) + 1),
            (ord(origin_file) + 2, int(origin_rank) + 1),
            (ord(origin_file) - 2, int(origin_rank) - 1),
            (ord(origin_file) + 2, int(origin_rank) - 1),
            (ord(origin_file) - 1, int(origin_rank) - 2),
            (ord(origin_file) + 1, int(origin_rank) - 2),
        ]
        for file, rank in possible_moves:
            if file < ord("a") or file > ord("h") or rank < 1 or rank > 8:
                continue
            square = self.position.board.get_square(chr(file), str(rank))
            if square.piece:
                if square.piece.colour != self.colour:
                    moves.append(f"{self.string}x{square.string}")
                continue
            moves.append(f"{self.string}{square.string}")
        return moves


class Bishop(Piece):
    """
    A class to represent a bishop chess piece.
    """

    def __init__(self, colour):
        super().__init__(colour)
        self.printable = "♝" if colour == "white" else "♗"
        self.string = "B"

    def list_possible_moves(self, _):
        """
        Returns a list of possible moves for the bishop.
        """
        if self.active is False:
            return []
        moves = []
        origin_file = self.position.file
        origin_rank = self.position.rank

        self._move_diagonally(moves, origin_file, origin_rank, "top_right")
        self._move_diagonally(moves, origin_file, origin_rank, "top_left")
        self._move_diagonally(moves, origin_file, origin_rank, "bottom_right")
        self._move_diagonally(moves, origin_file, origin_rank, "bottom_left")
        return moves


class Queen(Piece):
    """
    A class to represent a queen chess piece.
    """

    def __init__(self, colour):
        super().__init__(colour)
        self.printable = "♛" if colour == "white" else "♕"
        self.string = "Q"

    def list_possible_moves(self, _):
        """
        Returns a list of possible moves for the queen.
        """
        if self.active is False:
            return []
        moves = []
        origin_file = self.position.file
        origin_rank = self.position.rank

        self._move_orthogonally(moves, origin_file, origin_rank, "forward")
        self._move_orthogonally(moves, origin_file, origin_rank, "backward")
        self._move_orthogonally(moves, origin_file, origin_rank, "left")
        self._move_orthogonally(moves, origin_file, origin_rank, "right")
        self._move_diagonally(moves, origin_file, origin_rank, "top_right")
        self._move_diagonally(moves, origin_file, origin_rank, "top_left")
        self._move_diagonally(moves, origin_file, origin_rank, "bottom_right")
        self._move_diagonally(moves, origin_file, origin_rank, "bottom_left")
        return moves


class King(Piece):
    """
    A class to represent a king chess piece.
    """

    def __init__(self, colour):
        super().__init__(colour)
        self.printable = "♚" if colour == "white" else "♔"
        self.string = "K"

    def list_possible_moves(self, _):
        """
        Returns a list of possible moves for the king.
        """
        if self.active is False:
            return []
        moves = []
        origin_file = self.position.file
        origin_rank = self.position.rank
        possible_moves = [
            (ord(origin_file) - 1, int(origin_rank) + 1),
            (ord(origin_file), int(origin_rank) + 1),
            (ord(origin_file) + 1, int(origin_rank) + 1),
            (ord(origin_file) - 1, int(origin_rank)),
            (ord(origin_file) + 1, int(origin_rank)),
            (ord(origin_file) - 1, int(origin_rank) - 1),
            (ord(origin_file), int(origin_rank) - 1),
            (ord(origin_file) + 1, int(origin_rank) - 1),
        ]
        for file, rank in possible_moves:
            if file < ord("a") or file > ord("h") or rank < 1 or rank > 8:
                continue
            square = self.position.board.get_square(chr(file), str(rank))
            if square.piece:
                if square.piece.colour != self.colour:
                    moves.append(f"{self.string}x{square.string}")
                continue
            moves.append(f"{self.string}{square.string}")
        return moves


if __name__ == "__main__":
    # True to play the game, False to convert moves
    if False:
        chessGame = Game()
        chessGame.start_game()
    else:
        moves_to_convert = "e4 e5 Nc3 Nc6 f4 exf4 Nf3 Bb4 d4 Bxc3+ bxc3 d5 e5 f6 Bxf4 fxe5 Bxe5 Nxe5 Nxe5 Qe7 Bd3 c5 O-O Nf6 Qf3 Bg4 Qf2 O-O Qg3 Ne4 Bxe4 dxe4 Qxg4 cxd4 cxd4 Rad8 c3 b5 Qxe4 Qa3 Rf3 a5 Raf1 Qxa2 Rxf8+ Rxf8 Rxf8+ Kxf8 Qa8+ Ke7 Nc6+ Ke6 Nxa5 Qe2 Qc6+ Kf5 Qf3+ Qxf3 gxf3 Kf4 Kf2 g5 d5 Ke5 Nc6+ Kxd5 Nd4 b4 cxb4 Kxd4 Kg3 Kc4 Kg4 h6 Kh5 Kxb4 Kxh6 Kc5 Kxg5 Kd6 Kg6 Ke5 h4 Kf4 h5 Kxf3 h6 Ke2 h7 Kd3 h8=Q Kc4 Qe5 Kd3 Kf5 Kc4 Kf4 Kd3 Qe4+ Kc3 Kf3 Kd2 Qe3+ Kc2 Kf2 Kd1 Qe2+ Kc1 Ke3 Kb1 Kd3 Kc1 Qf2 Kd1 Qf1#".split()
        convert_notation(moves_to_convert)
