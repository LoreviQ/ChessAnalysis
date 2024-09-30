import threading
import tkinter as tk

from database import DBConnection
from server import initialize_server


def start_app():
    db_connection = DBConnection()
    my_server = initialize_server(db_connection)
    server_thread = threading.Thread(
        target=my_server.run, kwargs={"debug": True, "use_reloader": False}
    )
    server_thread.start()
    root = tk.Tk()
    root.title("Server Status")

    data_label = tk.Label(root, text="No data received", font=("Helvetica", 16))
    data_label.pack(pady=10)
    root.mainloop()


if __name__ == "__main__":
    start_app()
