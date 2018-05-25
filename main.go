package main

import (
	"fmt"
	"line-boi/cache"
	"line-boi/models"
	"line-boi/service"
	"line-boi/service/delivery/http"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/bluele/gcache"
	"github.com/labstack/echo"
	"github.com/line/line-bot-sdk-go/linebot"
)

var Cache = cache.NewCache(50)

func main() {

	go startService()

	for {
		messages := make(chan string, 1)
		go taskCheckAllBankService(messages)
		manageTask(messages)
		time.Sleep(time.Second * 45)
	}

}

func restartService(text string) {

	allData := Cache.Gcache.GetALL()

	for key, value := range allData {
		if strings.Contains(key.(string), "restart") {
			if value.(string) == text {
				cmdStr := "test.sh"
				os.Setenv("SERVICE_NAME_RESTART", strings.Split(key.(string), ":")[0])
				cmd := exec.Command("/bin/sh", cmdStr)

				data, err := cmd.Output()

				if err != nil {
					fmt.Println(err)
					log.Fatal(err)
				}

				fmt.Println(strings.NewReader(string(data)))
			}
		}
	}
}

func manageTask(messages chan string) {
	for {
		select {
		case msg := <-messages:
			{
				allCache := Cache.Gcache.GetALL()
				fmt.Println(allCache)
				var allUserID []string
				for key, value := range allCache {
					if strings.Contains(key.(string), ("line_id")) {
						allUserID = append(allUserID, value.(string))
					}
				}
				if len(allUserID) > 0 && len(msg) > 0 {
					_, err := connectLineBot().Multicast(allUserID, linebot.NewTextMessage(msg)).Do()
					if err != nil {
						log.Fatal(err)
					}

				}
				fmt.Println("service errors :", msg)
				return
			}
		}

	}
}

func taskCheckAllBankService(ch chan string) {
	ch <- strings.Join(service.StartPingAllServices(service.BankServicesInfo), ",")
}

func connectLineBot() *linebot.Client {
	bot, err := linebot.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}
	return bot

}
func startService() {
	e := echo.New()

	http.NewServiceHTTPHandler(e, connectLineBot(), Cache)
	e.Logger.Fatal(e.Start(":6000"))

}

func init() {
	service.BankServicesInfo = service.GetBankCoreServiceInfo()
}

// export CHANNEL_SECRET="34a3b63f5ade75d28a9b20f226914e6a"
// export CHANNEL_TOKEN="93Yob89UyPpSYFaWCOsTxEvhVuQHFCxr+rhO/9iOYtk2F0h+Z7mQDVA+EaVSdIG+wm75JXS7b0EeBRVEHFkytQ0YS/D1z3kJuzcfkKFQ/glknI8biE0WBOVv9j1v2QDN/LCO7cJ7COQZ343exnIbIAdB04t89/1O/w1cDnyilFU="

func NewCache(size int) *models.CacheService {
	return &models.CacheService{
		Gcache: gcache.New(size).
			AddedFunc(func(key, value interface{}) {
				fmt.Println("added key:", key)
			}).Build(),
	}

}
