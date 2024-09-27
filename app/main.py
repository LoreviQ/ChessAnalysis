from flask import Flask, jsonify, request
from flask_cors import CORS

app = Flask(__name__)
CORS(app)


@app.route("/readiness")
def readiness():
    """
    Health check endpoint
    """
    return jsonify({"status": "success"})


@app.route("/update_moves", methods=["POST"])
def update_moves():
    """
    Update moves endpoint
    """
    data = request.get_json()
    print(data)
    return jsonify({"status": "success"})


if __name__ == "__main__":
    app.run(debug=True)
