package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

type PortResult struct {
	Host string `json:"host"`
	Ports []PortStatus `json:"ports"`
}

type PortStatus struct {
	Port string `json:"port"`
	Status string `json:"status"`
	Open bool `json:"open"`
}

func raw_connect(host string, port string) PortStatus {
	ps := PortStatus{}
	ps.Port = port

	timeout := time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
	if err != nil {
			ps.Status = fmt.Sprint("CLOSED : ", err)
			ps.Open = false
	}
	if conn != nil {
			defer conn.Close()
			ps.Status = "OPEN"
			ps.Open = true
	}
	return ps
}

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// Greet returns a greeting for the given name
func (a *App) PortCheck(host string, start int, num int) string {

	pr := PortResult{}
	pr.Host = host

	for i := start; i < start + num; i++ {
		port := fmt.Sprintf("%d", i)
		ps := raw_connect(host, port)
		pr.Ports = append(pr.Ports, ps)
	}

	outputJson, _ := json.Marshal(&pr)
	return fmt.Sprint(string(outputJson))
}