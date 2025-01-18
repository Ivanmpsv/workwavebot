from flask import Flask, request, jsonify, make_response
import sqlite3
import logging
from db_functions import create_table, add_client, delete_client, update_client_formula, calculate_payment, get_all_clients, add_admin, remove_admin, get_all_admins

app = Flask(__name__)

# Инициализация логирования
logging.basicConfig(level=logging.ERROR, format='%(asctime)s - %(levelname)s - %(message)s')

@app.route('/add_client', methods=['POST'])
def api_add_client():
    data = request.json
    name = data.get('name')
    formula = data.get('formula')

    # Проверка на отсутствие обязательных параметров
    if not name or not formula:
        return make_response(jsonify({"status": "error", "message": "Name and formula are required."}), 400)
    
    try:
        # Попытка добавления клиента
        add_client(name, formula)
        return make_response(jsonify({"status": "success", "message": f"Client '{name}' added."}), 201)
    
    except sqlite3.IntegrityError:
        # Ошибка при добавлении клиента, если имя уже существует
        return make_response(jsonify({"status": "error", "message": f"Client '{name}' already exists."}), 409)
    
    except Exception as e:
        # Логируем ошибку и возвращаем внутреннюю ошибку сервера
        logging.error(f"Error adding client: {e}")
        return make_response(jsonify({"status": "error", "message": "Internal server error."}), 500)


@app.route('/delete_client/<name>', methods=['DELETE'])
def api_delete_client(name):
    if not name:
        return make_response(jsonify({"status": "error", "message": "Client name is required."}), 400)
    
    try:
        if delete_client(name):
            return make_response(jsonify({"status": "success", "message": f"Client '{name}' deleted."}), 200)
        else:
            return make_response(jsonify({"status": "error", "message": f"Client '{name}' not found."}), 404)
    except Exception as e:
        logging.error(f"Error deleting client '{name}': {e}")
        return make_response(jsonify({"status": "error", "message": "Internal server error."}), 500)


@app.route('/add_admin', methods=['POST'])
def api_add_admin():
    data = request.json
    id = data.get('id')

    # Проверка на отсутствие обязательных параметров
    if not id:
        return make_response(jsonify({"status": "error", "message": "Id is required."}), 400)
    
    try:
        # Попытка преобразовать строку в целое число
        int_id = int(id)
    except ValueError:
        return make_response(jsonify({"status": "error", "message": "Admin id is required as number."}), 400)

    try:
        # Попытка добавления админа
        add_admin(int_id)
        return make_response(jsonify({"status": "success", "message": f"Admin '{int_id}' added."}), 201)
    
    except sqlite3.IntegrityError:
        # Ошибка при добавлении админа, если имя уже существует
        return make_response(jsonify({"status": "error", "message": f"Admin '{int_id}' already exists."}), 409)
    
    except Exception as e:
        # Логируем ошибку и возвращаем внутреннюю ошибку сервера
        logging.error(f"Error adding admin: {e}")
        return make_response(jsonify({"status": "error", "message": "Internal server error."}), 500)


@app.route('/remove_admin/<id>', methods=['DELETE'])
def api_remove_admin(id):
    if not id:
        return make_response(jsonify({"status": "error", "message": "Admin id is required."}), 400)
    
    try:
        # Попытка преобразовать строку в целое число
        int_id = int(id)
    except ValueError:
        return make_response(jsonify({"status": "error", "message": "Admin id is required as number."}), 400)
    
    try:
        if remove_admin(int_id):
            return make_response(jsonify({"status": "success", "message": f"Admin '{int_id}' removed."}), 200)
        else:
            return make_response(jsonify({"status": "error", "message": f"Admin '{int_id}' not found."}), 404)
    except Exception as e:
        logging.error(f"Error removing admin '{int_id}': {e}")
        return make_response(jsonify({"status": "error", "message": "Internal server error."}), 500)


@app.route('/get_all_admins', methods=['GET'])
def api_get_all_admins():
    try:
        admins = get_all_admins()
        return make_response(jsonify({"status": "success", "admins": admins}), 200)
    except Exception as e:
        logging.error(f"Error retrieving all Admins: {e}")
        return make_response(jsonify({"status": "error", "message": "Internal server error."}), 500)


@app.route('/update_client_formula', methods=['PUT'])
def api_update_client_formula():
    data = request.json
    name = data.get('name')
    new_formula = data.get('new_formula')

    # Проверка на отсутствие обязательных параметров
    if not name or not new_formula:
        return make_response(jsonify({"status": "error", "message": "Name and new formula are required."}), 400)

    try:
        # Проверка формулы на ошибки
        safe_formula = new_formula.replace("payment", "(" + str(100) + ")")
        payment = eval(safe_formula, {"__builtins__": None}, {})
    except Exception as calc_error:
        logging.error(f"Error updating client formula: {calc_error}")
        return make_response(jsonify({"status": "error", "message": "Formula is not correct"}), 400)

    try:
        # Попытка обновления формулы
        if update_client_formula(name, new_formula):
            return make_response(jsonify({"status": "success", "message": f"Client '{name}' formula updated."}), 200)
        else:
            return make_response(jsonify({"status": "error", "message": f"Client '{name}' not found."}), 404)
    except Exception as e:
        logging.error(f"Error updating client formula: {e}")
        return make_response(jsonify({"status": "error", "message": "Internal server error."}), 500)


@app.route('/calculate_payment', methods=['POST'])
def api_calculate_payment():
    data = request.json
    client_name = data.get('client_name')
    payment = data.get('payment')

    # Проверка на отсутствие обязательных параметров
    if not client_name or payment is None:
        return make_response(jsonify({"status": "error", "message": "Client name and amount are required."}), 400)

    try:
        payment, formula = calculate_payment(client_name, payment)
        if payment is None:
            return make_response(jsonify({"status": "error", "message": f"Client '{client_name}' not found."}), 404)
        return make_response(jsonify({"status": "success", "payment": payment, "formula": formula}), 200)
    except Exception as e:
        logging.error(f"Error calculating payment for client '{client_name}': {e}")
        return make_response(jsonify({"status": "error", "message": "Internal server error."}), 500)


@app.route('/get_all_clients', methods=['GET'])
def api_get_all_clients():
    try:
        clients = get_all_clients()
        return make_response(jsonify({"status": "success", "clients": clients}), 200)
    except Exception as e:
        logging.error(f"Error retrieving all clients: {e}")
        return make_response(jsonify({"status": "error", "message": "Internal server error."}), 500)


if __name__ == '__main__':
    create_table()  # Убедитесь, что таблица существует при запуске
    app.run(debug=False, host='0.0.0.0', port=5000)
