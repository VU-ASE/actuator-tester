import socket
import time

rover_id = int(input("Please input the id of your rover: "))

ip_component = 100 + rover_id

UDP_IP = f"192.168.0.{ip_component}"
UDP_PORT = 8080

channel = int(input("Select which channel to test 0, 1, or 2: "))

sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM) # UDP

values = [0, 0.2, 0.5, 0.7, 1.0, 0.7, 0.5, 0.2, 0, -0.2, -0.5, -0.7, -1.0, -0.7, -0.5, -0.2, 0.0]

for val in values:
    print(f"Testing channel {channel} with: {val}")
    sock.sendto(bytes(f"{{\"channel\":{channel},\"value\":{val}}}", "utf-8"), (UDP_IP, UDP_PORT))
    time.sleep(0.5)

print("done")
