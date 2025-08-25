package helper

import (
	"encoding/json"
	"fmt"
	"net"
	"networkInfo/models"
	"os"
	"sync"
	"time"
)

func GenerateIPs(subNet net.IP, mask net.IPMask) []string {
	// GENERATING ALL POSSIBLE IPs IN SUBNET
	var ips []string

	// SHRINK THE SUBNET FROM ~65,000 TO ONLY SCAN AT MOST 254 -- VALID IP RANGE
	ones, bits := mask.Size()
	if ones < 24 {
		mask = net.CIDRMask(24, bits)
		ones, _ = mask.Size()
	}

	var network net.IP = subNet.Mask(mask)
	var hostBits int = bits - ones
	var possibleIPinSubset int = 1<<hostBits // shift bits of a left by b positions 

	var ip net.IP = make(net.IP, len(network))
	copy(ip, network)

	for i:=1; i < possibleIPinSubset - 1; i++ {
		ip[3] = byte(int(network[3]) + i)
		ips = append(ips, ip.String())
	}

	return ips
}

func IPIsReachable(ip string, timeout time.Duration) bool {
	connection, err := net.DialTimeout("tcp", net.JoinHostPort(ip, "80"), timeout) 
	if err!=nil {
		return false
	}

	connection.Close()
	return true
}

func ScanSubnet(ip string, mask net.IPMask) []models.NetworkDevicesInfo {

	var devices []models.NetworkDevicesInfo
	var candidates []string = GenerateIPs(net.ParseIP(ip), mask)

	// CONCURRENT REDUCE TIME // I/O BOUND OPERATIONS
	jobs := make(chan string)
	var wg sync.WaitGroup
	var mu sync.Mutex

	workerFunc := func(workerID int, jobs <- chan string){
		defer wg.Done()
		for job := range jobs {

			if IPIsReachable(job, 400 * time.Millisecond) {

				hostNames, err := net.LookupAddr(job)

				var hostName string 

				if err == nil && len(hostNames) > 0 {
					hostName = hostNames[0]
				}

				mu.Lock()
					devices = append(devices, models.NetworkDevicesInfo{
						IP: job,
						HostName: hostName,
					})
				mu.Unlock()

			}
		}
	}

	var workerCount int = 300
	wg.Add(workerCount)

	for worker := 1; worker <= workerCount; worker++ {
		go workerFunc(worker, jobs)
	}
	
	for _, candidate := range candidates {
		jobs <- candidate
	}

	close(jobs)
	wg.Wait()

	return devices
}


func SaveIndentedJSONReport(fileName string, report models.JSONReport) error {
	jsonData, err := json.MarshalIndent(report, "", "	")
	if err!=nil {
		return err
	}

	if err := os.WriteFile(fileName, jsonData, 0644); err!=nil { // owner (read 4 + write 2), groups (read 4), guest (read 4) => 0644
		return err
	}
	
	return nil
}

func ShowProcessing(done chan bool, start time.Time){
	var spinner []string = []string{"|", "/", "-", "\\"}

	var i int 

	for {
		select {
		case <- done: 
			var duration time.Duration= time.Since(start)
			fmt.Printf("\rScan Took %.2f seconds to complete\n", duration.Seconds())
		return
		default:
			fmt.Printf("\rScanning Network [%s]", spinner[i%len(spinner)])
			i++
			time.Sleep(100 * time.Millisecond)
		}
	}
}