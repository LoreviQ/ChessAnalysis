"""
Module to define the server class and its routes
"""

from database import DBConnection
from flask import Flask, jsonify, make_response, request
from flask_cors import CORS


def initialize_server(db_connection):
    """
    Start the Flask server
    """
    server = Flask(__name__)
    CORS(server)
    setup_routes(server, db_connection)
    return server


def setup_routes(server, db_connection):
    """
    Define the server routes
    """

    @server.route("/readiness", methods=["GET"])
    def readiness():
        """
        Health check endpoint
        """
        return make_response("OK", 200)

    @server.route("/update_moves", methods=["POST"])
    def update_moves():
        """
        Update moves endpoint
        """
        data = request.get_json()
        moves = data.get("moves")
        if not moves:
            return jsonify({"status": "error", "message": "No moves provided"})
        game_id = data.get("game_id")
        print(f"Received moves: {moves}")
        print(f"Game ID: {game_id}")
        db_connection.insert_move(moves, game_id)
        return jsonify({"status": "success"})


if __name__ == "__main__":
    db = DBConnection()
    my_server = initialize_server(db)
    my_server.run(debug=True, use_reloader=False)
