package service

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
	BankServicesInfo  *models.ServicesInfo
)
var bankServiceInfo = map[string]string{
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

// PingService provides the function that send the serviceName and services information to match and validate is it online or not.
func PingService(messageServiceName string, servicesInfo *models.ServicesInfo) string {
	var serviceName string
	var servicePort string
	var serviceIPAddress string

	serviceInfo, err := FindServiceName(messageServiceName, servicesInfo)

	if err != nil {
		return err.Error()
	}

	serviceName = serviceInfo.ServiceName
	servicePort = serviceInfo.Port
	serviceIPAddress = serviceInfo.IPAddress

	if len(serviceName) > 0 && len(servicePort) > 0 {
		message, _ := validateIsServiceDown(pingServer(serviceName, serviceIPAddress, servicePort))
		return message
	} else {
		return "Sorry, the name did not match to any services in our system."
	}
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
func StartPingAllServices(servicesInfo *models.ServicesInfo) []string {
	var lstServiceDowns []string
	for _, serviceDetail := range *servicesInfo {
		if message, isServiceDown := validateIsServiceDown(pingServer(serviceDetail.ServiceName, serviceDetail.IPAddress, serviceDetail.Port)); isServiceDown {
			lstServiceDowns = append(lstServiceDowns, message)
		}
	}

	return lstServiceDowns

}

func validateIsServiceDown(serviceName string, status bool) (string, bool) {
	if status {
		return fmt.Sprintf("%s service is working pretty well.", serviceName), false
	} else {
		return fmt.Sprintf("%s service is down, please contact admin.", serviceName), true
	}

}

func pingServer(serviceName string, ipAddress string, port string) (string, bool) {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%s", ipAddress, port), time.Duration(5)*time.Second)
	if err != nil {
		return serviceName, false
	}
	conn.Close()
	return serviceName, true
}

// GetBankCoreServiceInfo provides the all service information that using in BankCore project.
func GetBankCoreServiceInfo() *models.ServicesInfo {
	var services models.ServicesInfo
	for BankServiceName, BankServicePort := range bankServiceInfo {
		serviceInfo := models.ServiceInfo{
			ServiceName: BankServiceName,
			IPAddress:   bankCoreIPAddress,
			Port:        BankServicePort}

		services = append(services, serviceInfo)
	}
	return &services
}
