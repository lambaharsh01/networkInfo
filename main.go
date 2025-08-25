package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/user"
	"runtime"
	"time"

	"github.com/lambaharsh01/networkInfo/helper"
	"github.com/lambaharsh01/networkInfo/models"
)

func main() {

	done := make(chan bool)
	go helper.ShowProcessing(done, time.Now())

	localHost, err:= localHostInfo() 
	if err!=nil {
		log.Fatalf("failed to get localhost info %v", err)
	}

	localNetwork, err := getNetworkDevices()
	if err!=nil {
		log.Fatalf("failed to get local network info %v", err)
	}

	var report models.JSONReport = models.JSONReport{
		LocalHost: localHost,
		NetworkDevices:localNetwork,
	}

	var out *string = flag.String("out", "", "Output JSON file") 
	flag.Parse()

	var outFileName string = *out

	if outFileName == "" {
		outFileName = fmt.Sprintf("report_%s.json", time.Now().Format("2006-01-02_15-04-05"))
	}

	outFileName = fmt.Sprintf("logs/%s", outFileName)

	if err:= helper.SaveIndentedJSONReport(outFileName, report); err!=nil {
		log.Fatalf("failed to save report %v", err)
	}

	fmt.Println("\rNetwork scan completed")
	done <- true


}

func localHostInfo() (models.LocalHostInfo, error) {

	var info models.LocalHostInfo = models.LocalHostInfo{
		OS: runtime.GOOS,
	}

	// GET DEVICE NAME
	hostName, err := os.Hostname()
	if err!=nil {
		return info, err
	}
	info.HostName = hostName
	
	// GET USER's OS ACCOUNT INFO
	osUser, err := user.Current()
	if err!=nil {
		return info, err
	}
	info.UserName = osUser.Username

	// GET USER's NIC(Network Interface Card) | Network Adaptor FORM OS
	nics, err := net.Interfaces()
	if err!=nil {
		return info, err
	}
	// fmt.Println(nics)

	var networkDevices []string

	for _, nic := range nics {

		// GET ALL IP ADDs IN THE NETWORK INTERFACE
		ipAddrs, err := nic.Addrs()
		if err!=nil {
			return info, err
		}

		for _, addr := range ipAddrs {

			// CAST IPNET AND FILTER LOOPBACK ADDRS
			if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
				if ipNet.IP.To4() != nil {
					networkDevices = append(networkDevices, ipNet.IP.String())
				}

			}

		}

	}

	info.IPs = networkDevices

	return info, nil
}

func getNetworkDevices() ([]models.NetworkDevicesInfo, error) {

	var allDevices []models.NetworkDevicesInfo


	nics, err := net.Interfaces()
	if err!=nil {
		return allDevices, err
	}
	
	for _, nic := range nics {

		// GET ALL IP ADDs IN THE NETWORK INTERFACE
		ipAddrs, err := nic.Addrs()
		if err!=nil {
			return nil, err
		}

		for _, addr := range ipAddrs {

			// CAST IPNET AND FILTER LOOPBACK ADDRS
			if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
				if ipNet.IP.To4() != nil {
					devices := helper.ScanSubnet(ipNet.IP.String(), ipNet.Mask) 
					allDevices = append(allDevices, devices...)
				}

			}

		}

	}

	return allDevices, nil

}