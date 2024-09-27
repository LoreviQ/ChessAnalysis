from flask import Flask, jsonify, make_response, request
from flask_cors import CORS

server = Flask(__name__)
CORS(server)

received_data = []


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
    received_data.append(data)
    print(data)
    return jsonify({"status": "success"})


if __name__ == "__main__":
    server.run(debug=True)
