import sqlite3
import os
import logging

# Настройка логирования
logging.basicConfig(filename='app.log', level=logging.ERROR, format='%(asctime)s - %(levelname)s - %(message)s')

# Установите путь к файлу базы данных
db_path = os.path.join(os.path.dirname(__file__), 'recruitment.db')

# Подключение к базе данных SQLite с использованием контекстного менеджера
def connect_db():
    try:
        conn = sqlite3.connect(db_path)
        return conn
    except sqlite3.DatabaseError as e:
        logging.error(f"Database connection error: {e}")
        raise Exception("Failed to connect to the database.")

# Создание таблицы клиентов (если не существует)
def create_table():
    try:
        with connect_db() as conn:
            cur = conn.cursor()
            cur.execute(''' 
                CREATE TABLE IF NOT EXISTS clients (
                    id INTEGER PRIMARY KEY AUTOINCREMENT,
                    name TEXT UNIQUE NOT NULL,
                    formula TEXT NOT NULL
                );
            ''')
            cur.execute('''
            CREATE TABLE IF NOT EXISTS admins (
                    id INTEGER UNIQUE NOT NULL      
                );
            ''')
            conn.commit()
    except sqlite3.DatabaseError as e:
        logging.error(f"Error creating table: {e}")
        raise

# Добавление клиента
def add_client(name, formula):
    try:
        with connect_db() as conn:
            cur = conn.cursor()
            cur.execute('INSERT INTO clients (name, formula) VALUES (?, ?)', (name, formula))
            conn.commit()
    except sqlite3.IntegrityError as e:
        logging.error(f"Client '{name}' already exists: {e}")
        raise
    except sqlite3.DatabaseError as e:
        logging.error(f"Error adding client '{name}': {e}")
        raise Exception("Error adding client to the database.")
    
# Добавление администратора
def add_admin(id):
    # Проверка, что id является целым числом
    if not isinstance(id, int):
        raise ValueError("The provided id must be an integer.")
    
    try:
        with connect_db() as conn:
            cur = conn.cursor()
            cur.execute('INSERT INTO admins (id) VALUES (?)', (id,))
            conn.commit()
    except sqlite3.IntegrityError as e:
        logging.error(f"Admin '{id}' already exists: {e}")
        raise
    except sqlite3.DatabaseError as e:
        logging.error(f"Error adding admin with id '{id}': {e}")
        raise Exception("Error adding admin to the database.")
    
# Удаление администратора 
def remove_admin(id) -> bool:
    # Проверка, что id является целым числом
    if not isinstance(id, int):
        raise ValueError("The provided id must be an integer.")
    
    try:
        with connect_db() as conn:
            cur = conn.cursor()
            cur.execute('SELECT 1 FROM admins WHERE id = ?', (id,))
            admin_exists = cur.fetchone()

            if admin_exists:
                cur.execute('DELETE FROM admins WHERE id = ?', (id,))
                conn.commit()
                return True
            else:
                return False
    except sqlite3.DatabaseError as e:
        logging.error(f"Error removing admin with id '{id}': {e}")
        raise

# Получение всех администраторов
def get_all_admins():
    try:
        with connect_db() as conn:
            cur = conn.cursor()
            cur.execute("SELECT * FROM admins")
            admins = cur.fetchall()
            return admins
    except sqlite3.DatabaseError as e:
        logging.error(f"Error fetching admins: {e}")
        raise

# Удаление клиента
def delete_client(name: str) -> bool:
    try:
        with connect_db() as conn:
            cur = conn.cursor()
            cur.execute('SELECT 1 FROM clients WHERE name = ?', (name,))
            client_exists = cur.fetchone()

            if client_exists:
                cur.execute('DELETE FROM clients WHERE name = ?', (name,))
                conn.commit()
                return True
            else:
                return False
    except sqlite3.DatabaseError as e:
        logging.error(f"Error deleting client '{name}': {e}")
        raise

# Изменение формулы клиента
def update_client_formula(name: str, new_formula: str) -> bool:
    try:
        with connect_db() as conn:
            cur = conn.cursor()
            cur.execute('SELECT 1 FROM clients WHERE name = ?', (name,))
            client_exists = cur.fetchone()

            if client_exists:
                cur.execute('UPDATE clients SET formula = ? WHERE name = ?', (new_formula, name))
                conn.commit()
                return True
            else:
                return False
    except sqlite3.DatabaseError as e:
        logging.error(f"Error updating formula for client '{name}': {e}")
        raise

def calculate_payment(client_name, amount):
    try:
        with connect_db() as conn:
            cur = conn.cursor()
            cur.execute('SELECT formula FROM clients WHERE name = ?', (client_name,))
            result = cur.fetchone()

            if result:
                formula = result[0]
                # Преобразование строки формулы в арифметическое выражение
                # Использование eval с ограничениями может быть небезопасным,
                # вместо этого лучше использовать более безопасные способы вычислений
                try:
                    # Преобразуем формулу в математическое выражение, заменив "payment" на amount
                    safe_formula = formula.replace("payment", "(" + str(amount) + ")")
                    payment = eval(safe_formula, {"__builtins__": None}, {})
                    return payment, formula
                except Exception as calc_error:
                    logging.error(f"Error evaluating formula for client '{client_name}': {calc_error}")
                    raise
            else:
                logging.error(f"Client '{client_name}' not found in the database.")
                return None, None
    except Exception as e:
        logging.error(f"Error calculating payment for client '{client_name}': {e}")
        raise

# Получение всех клиентов
def get_all_clients():
    try:
        with connect_db() as conn:
            cur = conn.cursor()
            cur.execute("SELECT * FROM clients")
            clients = cur.fetchall()
            return clients
    except sqlite3.DatabaseError as e:
        logging.error(f"Error fetching clients: {e}")
        raise