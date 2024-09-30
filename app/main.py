"""
Main application to start the server and create the GUI
"""

import threading
import tkinter as tk

from database import DBConnection
from server import initialize_server


class App:
    """
    Class to create the GUI for the main application
    """

    def __init__(self, db_connection):
        self.db_connection = db_connection
        self.root = tk.Tk()
        self.server = None
        self.server_thread = None

    def start_app(self):
        """
        Start the application
        """
        # Start the server in a separate thread
        self.server = initialize_server(self.db_connection)
        self.server_thread = threading.Thread(target=self._run_server)
        self.server_thread.start()

        # Create the GUI
        self._initial_gui()
        self.root.protocol("WM_DELETE_WINDOW", self._on_close)
        self.root.mainloop()

    def _run_server(self):
        self.server.run(debug=True, use_reloader=False, threaded=True)

    def _initial_gui(self):
        self.root.title("Server Status")
        data_label = tk.Label(
            self.root, text="No data received", font=("Helvetica", 16)
        )
        data_label.pack(pady=10)

    def _on_close(self):
        if self.server_thread.is_alive():
            self.server_thread.join()
        self.root.destroy()


if __name__ == "__main__":
    App(DBConnection()).start_app()
