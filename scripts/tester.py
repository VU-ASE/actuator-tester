import socket
import time

UDP_IP = "192.168.0.101"
UDP_PORT = 8080
# MESSAGE = "{\"channel\":0,\"value\":0.7}"

sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM) # UDP

for i in range(3):
    sock.sendto(bytes(f"{{\"channel\":{i},\"value\":0.7\}}", "utf-8"), (UDP_IP, UDP_PORT))
    time.sleep(1)

    sock.sendto(bytes(f"{{\"channel\":{i},\"value\":-0.7}}", "utf-8"), (UDP_IP, UDP_PORT))
    time.sleep(1)

    sock.sendto(bytes(f"{{\"channel\":{i},\"value\":0}}", "utf-8"), (UDP_IP, UDP_PORT))
    time.sleep(1)

    print(f"finished channel {i}")
