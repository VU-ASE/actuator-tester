import socket
import time

rover_ip = str(input("IP Address of Rover (192.168.0.101): "))
 
UDP_PORT = 8080
print(f"Sending to {rover_ip}:{UDP_PORT}")

channel = int(input("Select which channel to test 0, 1, or 2: "))
sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)

while True:
    val = float(input(f"Enter value for channel {channel}: "))
    print(f"    Setting channel {channel} to {val}")
    sock.sendto(bytes(f"{{\"channel\":{channel},\"value\":{val}}}", "utf-8"), (rover_ip, UDP_PORT))
    time.sleep(0.1)
