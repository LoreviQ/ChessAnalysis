import threading
import tkinter as tk

from database import DBConnection
from server import Server


class App:
    def __init__(self):
        self.db_connection = DBConnection()
        self.server = Server(self.db_connection)

        self.root = tk.Tk()
        self.root.title("Server Status")

        self.data_label = tk.Label(
            self.root, text="No data received", font=("Helvetica", 16)
        )
        self.data_label.pack(pady=10)

    def start(self):
        self.update_gui()
        self.root.mainloop()

    def update_gui(self):
        data = self.db_connection.fetch_latest_move()
        if data:
            self.data_label.config(text=str(data[-1]))
        else:
            self.data_label.config(text="No data received")

        # Schedule the function to run again after 5 seconds
        self.root.after(5000, self.update_gui)


if __name__ == "__main__":
    app = App()
    server_thread = threading.Thread(target=app.server.run())
    server_thread.start()
    app.start()
