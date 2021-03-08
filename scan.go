package main

import (
	"bytes"
	"fmt"
	"net"
	"sort"
	"strconv"
	"time"
)

const timeout = time.Second * 15

type ScanItem struct {
	IP   string
	Port uint16
}

type ScanResult struct {
	IP       string
	Port     uint16
	Services string
}

type ScanError struct {
	Err error
}

type Scanner struct {
	ScanItems   []ScanItem
	ScanResults []ScanResult
	ScanErrors  []error // did not used
}

func (s *Scanner) Start() {
	var startTime time.Time
	var endTime time.Duration
	resChan := make(chan ScanResult, len(s.ScanItems))
	errChan := make(chan error, len(s.ScanItems))
	fmt.Printf("There are %d targets to scan.\n\n", len(s.ScanItems))
	startTime = time.Now()
	for _, item := range s.ScanItems {
		go item.Execute(resChan, errChan)
	}

outer:
	for {
		select {
		case res := <-resChan:
			s.ScanResults = append(s.ScanResults, res)
		case err := <-errChan:
			s.ScanErrors = append(s.ScanErrors, err)
		default:
			if len(s.ScanResults)+len(s.ScanErrors) == len(s.ScanItems) {
				endTime = time.Since(startTime)
				break outer
			}
		}
	}
	s.Result()
	fmt.Printf("Done in %v.\n", endTime)
}

func (s *Scanner) Result() {
	// Sort ScanResult only by IP and Port
	sort.Slice(s.ScanResults, func(i, j int) bool {
		res := bytes.Compare(net.ParseIP(s.ScanResults[i].IP), net.ParseIP(s.ScanResults[j].IP))
		if res == 0 {
			return s.ScanResults[i].Port < s.ScanResults[j].Port
		}
		return res < 0
	})
	IPLength := 0
	portLength := 5
	for _, val := range s.ScanResults {
		if len(val.IP) > IPLength {
			IPLength = len(val.IP)
		}
	}

	if len(s.ScanResults) == 0 {
		fmt.Println("All targets are close.")
		return
	}

	format := fmt.Sprintf("%%%ds  %%%dv  %%v\n", IPLength, portLength)
	fmt.Printf(format, "IPs", "Port", "Services")
	for _, val := range s.ScanResults {
		fmt.Printf(format, val.IP, val.Port, val.Services)
	}
}

func NewScanner(s string, ports []uint16) (*Scanner, error) {
	items := make([]ScanItem, 0)
	for i := len(s) - 1; i >= 0; i-- {
		switch s[i] {
		case '/':
			ip, ipNet, err := net.ParseCIDR(s)
			if err != nil {
				return nil, err
			}
			for ip = ip.Mask(ipNet.Mask); ipNet.Contains(ip); inc(ip) {
				for _, port := range ports {
					items = append(items, ScanItem{
						IP:   ip.String(),
						Port: port,
					})
				}
			}
			items = items[1 : len(items)-1]
			goto ret
		case '.', ':':
			ip := net.ParseIP(s)
			if ip == nil {
				return nil, fmt.Errorf("invalid IP address: %v", s)
			}
			for _, port := range ports {
				items = append(items, ScanItem{
					IP:   ip.String(),
					Port: port,
				})
			}
			goto ret
		}
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("invalid IP address: %v", s)
	}

ret:
	return &Scanner{
		ScanItems: items,
	}, nil
}

func inc(ip net.IP) {
	for i := len(ip) - 1; i >= 0; i-- {
		ip[i]++
		if ip[i] > 0 {
			return
		}
	}
}

func (e *ScanError) Error() string {
	if e == nil {
		return "<nil>"
	}
	return "ScanError: " + e.Err.Error()
}

func (i ScanItem) Execute(resChan chan<- ScanResult, errChan chan<- error) {
	address := "[" + i.IP + "]:" + strconv.FormatUint(uint64(i.Port), 10)
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		errChan <- &ScanError{Err: err}
		return
	}
	conn.Close()
	resChan <- ScanResult{
		IP:       i.IP,
		Port:     i.Port,
		Services: GetServiceName("tcp", i.Port),
	}
}
