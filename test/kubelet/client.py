import socket
import os



def main():
#     import time
#     start_time = time.time()  # 获取当前时间
#     while True:
#         current_time = time.time()  # 获取当前时间
#         elapsed_time = current_time - start_time  # 计算已经过去的时间
#         # 如果已经过了10秒，退出循环
#         if elapsed_time >= 5:
#             break
#         # 这里可以放置你想要循环执行的代码
#         print("循环执行中...")
#         with open("/home/python/client.txt", "a+") as f:
#             f.write("hahaha" + '\n')
#         time.sleep(1)  # 每次循环暂停1秒钟，避免过快执行
    host = 'localhost'
    port = int(os.environ.get('PORT_CLIENT', 8080))

    # 创建 socket 对象
    client_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

    # 连接服务器
    client_socket.connect((host, port))

    # 发送数据
    message = "Hello, server!"
    print("发送消息: {}".format(message))
    client_socket.sendall(message.encode('utf-8'))

    # 接收数据
    received_data = client_socket.recv(1024).decode('utf-8')
    with open("/home/python/client-"+str(port)+".txt", "w+") as f:
        f.write(received_data + '\n')
#     print("接收到的消息: {}".format(received_data))

    # 关闭连接
    client_socket.close()

if __name__ == "__main__":
    main()
