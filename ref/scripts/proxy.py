#!/usr/bin/python

import socket, thread

host = 'localhost'
port = 4123

def receiver(c):
	try:
		while True:
			msg = c.recv(1024)
			print msg
			c.send('HTTP/1.1 200 OK\n\nHello')
	except Exception as e:
		print e

if __name__ == '__main__':
	s = socket.socket()
	s.bind((host, port))
	s.listen(10)
	while True:
		c, _ = s.accept()
		thread.start_new_thread(receiver, tuple([c]))
