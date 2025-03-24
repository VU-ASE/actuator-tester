import socket
import time

UDP_IP = "192.168.0.101"
UDP_PORT = 8080

sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM) # UDP

sock.sendto(bytes(f"{{\"channel\":0,\"value\":0.7}}", "utf-8"), (UDP_IP, UDP_PORT))
time.sleep(1)

sock.sendto(bytes(f"{{\"channel\":0,\"value\":-0.7}}", "utf-8"), (UDP_IP, UDP_PORT))
time.sleep(1)

sock.sendto(bytes(f"{{\"channel\":0,\"value\":0}}", "utf-8"), (UDP_IP, UDP_PORT))
time.sleep(1)

print("finished channel 0")


speed_values = [0.2, 0.5, 0.7, 0, -0.2, -0.5, -0.7, 0]

for val in speed_values:
    print(f"current speed: {val}")

    sock.sendto(bytes(f"{{\"channel\":1,\"value\":{val}}}", "utf-8"), (UDP_IP, UDP_PORT))
    sock.sendto(bytes(f"{{\"channel\":2,\"value\":{val}}}", "utf-8"), (UDP_IP, UDP_PORT))
    time.sleep(3)

print("done")
