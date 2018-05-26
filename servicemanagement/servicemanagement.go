package servicemanagement

import (
	"errors"
	"fmt"
	models "line-boi/models"
	"net"
	"strings"
	"time"
)

var (
	bankCoreIPAddress = "192.168.212.212"
	bankServiceInfo   = map[string]string{
		"hulk":               "9906",
		"genesis":            "9900",
		"eve":                "9904",
		"ems":                "9901",
		"minio":              "9800",
		"maersk":             "9903",
		"simpleredirectbank": "3001",
		"lawson":             "9902",
		"portainer":          "3002",
		"phpmyadmin":         "3000",
	}
)

// PingService provides the function that send the serviceName and services information to match and validate is it online or not.
func PingService(messageServiceName string, servicesInfo *models.ServicesInfo, timeOut time.Duration) string {
	serviceInfo, err := FindServiceName(messageServiceName, servicesInfo)
	if err != nil {
		return "Sorry, the name did not match to any services in our system."
	}
	if len(serviceInfo.ServiceName) > 0 && len(serviceInfo.Port) > 0 {
		serviceStatus := ping(serviceInfo.ServiceName, serviceInfo.IPAddress, serviceInfo.Port, timeOut)
		message, _ := isServiceOnline(serviceInfo.ServiceName, serviceStatus)
		return message
	}
	return ""
}

func FindServiceName(messageText string, servicesInfo *models.ServicesInfo) (*models.ServiceInfo, error) {
	for _, serviceDetail := range *servicesInfo {
		if strings.Contains(strings.ToLower(messageText), strings.ToLower(serviceDetail.ServiceName)) {
			return &serviceDetail, nil
		}
	}
	return nil, errors.New("the name did not match to any services in our system.")
}

// StartPingAllServices provides the function ping to allservice that we send through input.
func StartPingAllServices(servicesInfo *models.ServicesInfo, timeOut time.Duration) []string {
	var lstServiceDowns []string
	for _, serviceDetail := range *servicesInfo {
		serviceStatus := ping(serviceDetail.ServiceName, serviceDetail.IPAddress, serviceDetail.Port, timeOut)
		if message, isOnline := isServiceOnline(serviceDetail.ServiceName, serviceStatus); !(isOnline) {
			lstServiceDowns = append(lstServiceDowns, message)
		}
	}
	return lstServiceDowns
}

func isServiceOnline(serviceName string, status bool) (string, bool) {
	if status {
		return fmt.Sprintf("%s service is working pretty well.", serviceName), true
	} else {
		return fmt.Sprintf("%s service is down, please contact admin.", serviceName), false
	}

}

func ping(serviceName string, ipAddress string, port string, timeOut time.Duration) bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%s", ipAddress, port), timeOut)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// NewBankCoreServiceInfo provides the all service information that using in BankCore project.
func NewBankCoreServiceInfo() *models.ServicesInfo {
	bankCoreServices := models.ServicesInfo{}
	for BankServiceName, BankServicePort := range bankServiceInfo {
		serviceInfo := models.ServiceInfo{
			ServiceName: BankServiceName,
			IPAddress:   bankCoreIPAddress,
			Port:        BankServicePort}
		bankCoreServices = append(bankCoreServices, serviceInfo)
	}
	return &bankCoreServices
}
