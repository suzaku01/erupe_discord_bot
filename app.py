from flask import Flask, request, jsonify

app = Flask(__name__)

players = 0  # サンプルとしてのプレイヤー数

@app.route('/add_user', methods=['POST'])
def add_user():
    global players
    players += 1
    return jsonify({'message': 'User added!', 'total_players': players})

@app.route('/delete_user', methods=['POST'])
def delete_user():
    global players
    players -= 1
    return jsonify({'message': 'User removed!', 'total_players': players})

@app.route('/get_players', methods=['GET'])
def get_players():
    return jsonify({"players": players})

if __name__ == "__main__":
    app.run(port=9998)