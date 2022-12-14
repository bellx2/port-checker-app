package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/showwin/speedtest-go/speedtest"
	"golang.org/x/exp/slices"
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

type NetworkInterface struct {
	IP string `json:"ip"`
	InterfaceName string `json:"name"`
}

type NetworkInfo struct{
	Interface []NetworkInterface `json:"interface"`
	GlobalIP string `json:"global_ip"`
	ISP string `json:"isp"`
}

type SpeedTestResult struct {
	IP string `json:"ip"`
	ISP string `json:"isp"`
	Download float64 `json:"download"`
	Upload float64 `json:"upload"`
	Latency string `json:"latency"`
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

func (a *App) NetInfo() string {
	var currentIP [] string
	addrs, _ := net.InterfaceAddrs()
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				// fmt.Println("Current IP address : ", ipnet.IP.String())
				currentIP = append(currentIP, ipnet.String())
			}
		}
	}
	var iflist [] NetworkInterface
	interfaces, _ := net.Interfaces()
	for _, interf := range interfaces {
		if addrs, err := interf.Addrs(); err == nil {
			for _, addr := range addrs {
				// fmt.Println("[", index, "]", interf.Name, ">", addr)
				if (slices.Contains(currentIP, addr.String())) {
					// fmt.Println("Use name : ", interf.Name, addr.String())
					iflist = append(iflist, NetworkInterface{addr.String(),interf.Name})
				}
			}
		}
	}

	user, _ := speedtest.FetchUserInfo()
	NetworkInfo := NetworkInfo{iflist, user.IP, user.Isp}

	outputJson, _ := json.Marshal(&NetworkInfo)
	return fmt.Sprint(string(outputJson))
}

func (a *App) SpeedTest() string {

	user, _ := speedtest.FetchUserInfo()

	serverList, _ := speedtest.FetchServers(user)
	targets, _ := serverList.FindServer([]int{})

	var TestResults [] SpeedTestResult

	for _, s := range targets {
		s.PingTest()
		s.DownloadTest(false)
		s.UploadTest(false)
		TestResults = append(TestResults, SpeedTestResult{user.IP, user.Isp, s.DLSpeed, s.ULSpeed, s.Latency.String()})
		// fmt.Printf("Latency: %s, Download: %f, Upload: %f\n", s.Latency, s.DLSpeed, s.ULSpeed)
	}
	
	outputJson, _ := json.Marshal(&TestResults)
	return fmt.Sprint(string(outputJson))
}