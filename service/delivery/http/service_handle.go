package http

import (
	"context"
	"fmt"
	"line-boi/models"
	"line-boi/service"
	"log"
	"strings"

	"github.com/labstack/echo"
	"github.com/line/line-bot-sdk-go/linebot"
)

type HTTPCallBackHanlder struct {
	Bot   *linebot.Client
	Cache *models.CacheService
}

func NewServiceHTTPHandler(e *echo.Echo, linebot *linebot.Client, cache *models.CacheService) {

	hanlders := &HTTPCallBackHanlder{Bot: linebot, Cache: cache}
	e.GET("ping", func(c echo.Context) error {

		return c.String(200, "Line boi Service : We are good thank you for asking us.")
	})

	e.POST("/callback", hanlders.Callback)
}

func (handler *HTTPCallBackHanlder) Callback(c echo.Context) error {

	ctx := c.Request().Context()

	if ctx == nil {
		ctx = context.Background()
	}

	events, err := handler.Bot.ParseRequest(c.Request())

	if err != nil {
		if err == linebot.ErrInvalidSignature {
			c.String(400, "badrequest")
		} else {
			c.String(500, "internal")
		}
	}
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				var userID = event.Source.UserID
				handler.Cache.Gcache.Set(fmt.Sprintf("line_id:%s", userID), userID)
				if restartDocker(message.Text) {
					serviceInfo, err := service.FindServiceName(message.Text, service.GetBankCoreServiceInfo())

					if err != nil {

						log.Print(err)
					}
					handler.Cache.Gcache.Set(fmt.Sprintf("%s:restart", serviceInfo.ServiceName), "adminonly")
				}
				if _, err = handler.Bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(service.PingService(message.Text, service.GetBankCoreServiceInfo()))).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	}

	return c.JSON(200, "")
}

func restartDocker(message string) bool {
	if strings.Contains(strings.ToLower(message), "restart") {
		return true
	} else {
		return false
	}
}
