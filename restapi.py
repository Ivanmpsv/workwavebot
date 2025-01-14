from flask import Flask, request, jsonify
from db_functions import create_table, add_client, delete_client, update_client_formula, calculate_payment, get_all_clients


app = Flask(__name__)

@app.route('/add_client', methods=['POST'])
def api_add_client():
    data = request.json
    name = data.get('name')
    formula = data.get('formula')
    add_client(name, formula)
    return jsonify({"status": 200, "message": f"Client '{name}' added."})

@app.route('/delete_client/<name>', methods=['DELETE'])
def api_delete_client(name):
    if delete_client(name):
        return jsonify({"status": 200, "message": f"Client '{name}' deleted."})
    else:
        return jsonify({"status": "error", "message": f"Client '{name}' not found."})

@app.route('/update_client_formula', methods=['PUT'])
def api_update_client_formula():
    data = request.json
    name = data.get('name')
    new_formula = data.get('new_formula')
    if update_client_formula(name, new_formula):
        return jsonify({"status": 200, "message": f"Client '{name}' formula updated."})
    else:
        return jsonify({"status": "error", "message": f"Client '{name}' not found."})

@app.route('/calculate_payment', methods=['POST'])
def api_calculate_payment():
    data = request.json
    client_name = data.get('client_name')
    amount = data.get('amount')
    payment, formula = calculate_payment(client_name, amount)
    return jsonify({"status": 200, "payment": payment, "formula": formula})

@app.route('/get_all_clients', methods=['GET'])
def api_get_all_clients():
    clients = get_all_clients()
    return jsonify({"status": 200, "clients": clients})

if __name__ == '__main__':
    create_table()
    app.run(host='0.0.0.0', port=5000)
