import threading
import tkinter as tk

from database import DBConnection
from server import start_server


def start_app():
    db_connection = DBConnection()
    server_thread = threading.Thread(target=start_server, args=(db_connection,))
    server_thread.start()
    root = tk.Tk()
    root.title("Server Status")

    data_label = tk.Label(root, text="No data received", font=("Helvetica", 16))
    data_label.pack(pady=10)
    root.mainloop()


if __name__ == "__main__":
    start_app()
