import socket
import os

def main():
    host = 'localhost'
    port = int(os.environ.get('PORT_SERVER', 8080))

    # 创建 socket 对象
    server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

    # 绑定地址和端口
    server_socket.bind((host, port))

    # 监听连接
    server_socket.listen(1)
    print("服务器正在监听端口 {} ...".format(port))

    # 等待客户端连接
    client_socket, addr = server_socket.accept()
    print("连接来自: {}".format(addr))

    # 接收数据并发送回客户端
    data = client_socket.recv(1024).decode('utf-8')
    print("接收到的数据: {}".format(data))
    client_socket.sendall(data.encode('utf-8'))

    # 关闭连接
    client_socket.close()

if __name__ == "__main__":
    main()
