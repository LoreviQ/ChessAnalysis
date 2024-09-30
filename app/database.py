"""
Mode to interact with the SQLite database
"""

import json
import sqlite3

from chess_game import convert_notation

QUERIES = {
    "CREATE_MOVES": """
        INSERT INTO moves (game_id, move_data) VALUES (?, ?)
    """,
    "CREATE_GAME": """
        INSERT INTO games (chessdotcom_id) VALUES (?) RETURNING id
    """,
    "GET_LATEST_GAME_ID": """
        SELECT id FROM games WHERE chessdotcom_id = ? ORDER BY created_at DESC LIMIT 1
    """,
}


class DBConnection:
    """
    Class to interact with the SQLite database
    Creates the database and tables if they do not exist
    """

    def __init__(self):
        self.db_path = "database.db"
        self._init_db()

    def _get_connection(self):
        conn = sqlite3.connect(self.db_path)
        return conn, conn.cursor()

    def _init_db(self):
        conn, cursor = self._get_connection()
        with open("schema.sql", "r", encoding="utf-8") as schema_file:
            schema = schema_file.read()
        cursor.executescript(schema)
        conn.commit()
        conn.close()

    def insert_move(self, moves, chessdotcom_id=None):
        """
        Insert moves into the database
        :param moves: list of moves in short algebraic notation
        :param chessdotcom_id: chess.com game ID
        """
        conn, cursor = self._get_connection()
        moves_str = self._standardize_moves(moves)
        cursor.execute(QUERIES["GET_LATEST_GAME_ID"], (chessdotcom_id,))
        result = cursor.fetchone()
        if result:
            game_id = result[0]
        else:
            cursor.execute(QUERIES["CREATE_GAME"], (chessdotcom_id,))
            game_id = cursor.fetchone()[0]
        cursor.execute(QUERIES["CREATE_MOVES"], (game_id, moves_str))
        conn.commit()
        conn.close()

    def fetch_latest_move(self):
        conn, cursor = self._get_connection()
        cursor.execute("SELECT move_data FROM moves ORDER BY id DESC LIMIT 1")
        result = cursor.fetchone()
        conn.close()
        return result[0] if result else "No data received"

    def _standardize_moves(self, moves):
        """
        :param moves: list of moves in short algebraic notation
        :return: string of space-seperated moves in long algebraic notation
        """
        moves = [move for i, move in enumerate(moves) if i % 3 != 0]
        moves = convert_notation(moves)
        return json.dumps(moves)
