import sqlite3
import os

# Установите путь к файлу базы данных
db_path = os.path.join(os.path.dirname(__file__), 'recruitment.db')

# Подключение к базе данных SQLite
def connect_db():
    return sqlite3.connect(db_path)

# Создание таблицы клиентов (если не существует)
def create_table():
    conn = connect_db()
    cur = conn.cursor()
    cur.execute('''
        CREATE TABLE IF NOT EXISTS clients (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT UNIQUE NOT NULL,
            formula TEXT NOT NULL
        );
    ''')
    conn.commit()
    cur.close()
    conn.close()

# Добавление клиента
def add_client(name, formula):
    conn = connect_db()
    cur = conn.cursor()
    try:
        cur.execute('INSERT INTO clients (name, formula) VALUES (?, ?)', (name, formula))
        conn.commit()
    except sqlite3.IntegrityError:
        print(f"Client '{name}' already exists.")
    finally:
        cur.close()
        conn.close()

# Удаление клиента
def delete_client(name: str) -> bool:
    conn = connect_db()
    cur = conn.cursor()

    # Проверка наличия клиента
    cur.execute('SELECT 1 FROM clients WHERE name = ?', (name,))
    client_exists = cur.fetchone()

    if client_exists:
        # Удаление клиента
        cur.execute('DELETE FROM clients WHERE name = ?', (name,))
        conn.commit()
        cur.close()
        conn.close()
        return True
    else:
        cur.close()
        conn.close()
        return False

# Изменение формулы клиента
def update_client_formula(name: str, new_formula: str) -> bool:
    conn = connect_db()
    cur = conn.cursor()

    # Проверка наличия клиента
    cur.execute('SELECT 1 FROM clients WHERE name = ?', (name,))
    client_exists = cur.fetchone()

    if client_exists:
        # Обновление формулы клиента
        cur.execute('UPDATE clients SET formula = ? WHERE name = ?', (new_formula, name))
        conn.commit()
        cur.close()
        conn.close()
        return True
    else:
        cur.close()
        conn.close()
        return False

# Расчет выплаты
def calculate_payment(client_name, amount):
    conn = connect_db()
    cur = conn.cursor()
    cur.execute('SELECT formula FROM clients WHERE name = ?', (client_name,))
    result = cur.fetchone()
    if result:
        formula = result[0]
        payment = eval(formula.replace("amount", str(amount)))
        print(f"\nPayment for {client_name} (formula: '{formula}') with amount {amount}: {payment}")
    else:
        print(f"Client '{client_name}' not found.")
    cur.close()
    conn.close()
    return payment, formula

def get_all_clients():
    conn = connect_db()
    cur = conn.cursor()
    cur.execute("SELECT * FROM clients")
    clients = cur.fetchall()
    cur.close()
    conn.close()
    return clients