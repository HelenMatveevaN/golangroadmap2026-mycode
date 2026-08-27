package main

//Написать конструктор с functional options для условного Server.
//Functional Options Pattern для Server

import (
	"fmt"
	"time"
)

type Server struct {
	host string
	port int
	timeout time.Duration
	maxConn int
}

type Option func(*Server)

func NewServer(host string, port int, opts ...Option) *Server {
	//init defaults:
	srv := &Server{
		host:		host,
		port:		port,
		timeout:	30 * time.Second, 
		maxConn:	10,
	}

	//use opts with funcs
	for _, opt := range opts {
		opt(srv)
	}

	return srv
}

func WithTimeout(d time.Duration) Option {
	return func(s *Server) {
		s.timeout = d
	}
}

func WithMaxConn(n int) Option {
	return func(s *Server) {
		s.maxConn = n
	}
}

func main() {
	defaultServer := NewServer("localhost", 8080)
	fmt.Printf("Дефолтный сервер:\n%+v\n\n", defaultServer)

	customServer := NewServer(
		"127.0.0.1",
		9000,
		WithTimeout(5*time.Second),
		WithMaxConn(100),
	)
	fmt.Printf("Кастомный сервер:\n%+v\n", customServer)

}