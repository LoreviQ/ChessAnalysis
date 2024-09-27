import threading
import time
import tkinter as tk

import requests
from server import received_data, server


# Function to run the Flask server
def run_server():
    server.run(debug=True, use_reloader=False)


# Function to check server status
def check_server_status():
    try:
        response = requests.get("http://127.0.0.1:5000/readiness")
        if response.status_code == 200:
            return "Server is active"
        else:
            return "Server is inactive"
    except requests.exceptions.RequestException:
        return "Server is inactive"


# Function to update the GUI
def update_gui():
    status = check_server_status()
    status_label.config(text=status)

    # Update received data
    if received_data:
        data_label.config(text=str(received_data[-1]))
    else:
        data_label.config(text="No data received")

    # Schedule the function to run again after 5 seconds
    root.after(5000, update_gui)


# Start the Flask server in a separate thread
server_thread = threading.Thread(target=run_server)
server_thread.daemon = True
server_thread.start()

# Create the Tkinter window
root = tk.Tk()
root.title("Server Status")

status_label = tk.Label(root, text="Checking server status...", font=("Helvetica", 16))
status_label.pack(pady=10)

data_label = tk.Label(root, text="No data received", font=("Helvetica", 16))
data_label.pack(pady=10)

# Start the periodic update
update_gui()

# Run the Tkinter event loop
root.mainloop()
