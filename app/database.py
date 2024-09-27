import os
import sqlite3


class DBConnection:
    def __init__(self):
        self.db_path = "database.db"
        self._init_db()

    def _get_connection(self):
        conn = sqlite3.connect(self.db_path)
        return conn, conn.cursor()

    def _init_db(self):
        conn, cursor = self._get_connection()
        cursor.execute(
            """
            CREATE TABLE IF NOT EXISTS moves (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                move_data TEXT NOT NULL
            )
        """
        )
        conn.commit()
        conn.close()

    def insert_move(self, data):
        # conn, cursor = self._get_connection()
        # cursor.execute("INSERT INTO moves (move_data) VALUES (?)", (data,))
        # conn.commit()
        # conn.close()
        standardize_moves(data)

    def fetch_latest_move(self):
        conn, cursor = self._get_connection()
        cursor.execute("SELECT move_data FROM moves ORDER BY id DESC LIMIT 1")
        result = cursor.fetchone()
        conn.close()
        return result[0] if result else "No data received"


def standardize_moves(moves):
    standardized_moves = ""
    for i in range(0, len(moves), 3):
        standardized_moves += f"{moves[i + 1]},{moves[i + 2]},"
    print(standardized_moves[:-1])
