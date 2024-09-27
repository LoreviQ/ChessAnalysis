from database import DBConnection
from flask import Flask, jsonify, make_response, request
from flask_cors import CORS


class Server:
    def __init__(self, db_connection):
        self.db_connection = db_connection
        self.server = Flask(__name__)
        CORS(self.server)
        self.setup_routes()

    def setup_routes(self):
        @self.server.route("/readiness", methods=["GET"])
        def readiness():
            """
            Health check endpoint
            """
            return make_response("OK", 200)

        @self.server.route("/update_moves", methods=["POST"])
        def update_moves():
            """
            Update moves endpoint
            """
            data = request.get_json()
            self.db_connection.insert_move(data)
            return jsonify({"status": "success"})

    def run(self):
        self.server.run(debug=True, use_reloader=False)


if __name__ == "__main__":
    db = DBConnection()
    server = Server(db)
    server.run()
