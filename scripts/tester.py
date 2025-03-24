import socket
import time

rover_id = int(input("Please input the id of your rover: "))

ip_component = 100 + rover_id

UDP_IP = f"192.168.0.{ip_component}" # change if necessary
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
