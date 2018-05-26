# line-boi

Line-boi is the example ping service that we receive the message from line bot through webhook.



### How to instaill

#### Required:
* Golang
* Echo library from github.com/labstack/echo
* Line bot library from github.com/line/line-bot-sdk-go/linebot

***Choose which one between (Ngrok, or Heroku) to build your server online.***
* [Ngrok](https://ngrok.com/) to forward your IP Address to public for let linebot can send message through webhook (Recommend for private service, and newbie developer because this is the easy way to run public your IP Address to connect to linebot)
* [Heroku](https://www.heroku.com/) if your service are public now, you can use heroku to build your project to cloud service.


### Steps:
- you have to export ``` CHANNEL_SECRET, and CHANNEL_TOKEN ``` from [Line Bot](https://access.line.me/dialog/oauth/weblogin?client_id=1459630796&redirect_uri=https://business.line.me/auth?redirectUri%3Dhttps://business.line.me/sso/auth?response_type%253Dcode%2526scope%253Dopenid%2526client_id%253D1%2526redirect_uri%253Dhttps%25253A%25252F%25252Fadmin-official.line.me%25252Fs%2526state%253DXsBoBb7xedfN&response_type=code&state=3OtFfK#/bulk) as environment
- you can setting your configuration service at servicemanagement file
- type command ```go run main.go```
- last step make this server to online through heroku, or Ngrok, but remember one thing if you using heroku to make this server online, please make sure your server can ping to your private service.

### Example for using ngrok to forward IP Address
- Install [Ngrok](https://ngrok.com/) And sign up Ngrok account
- Unzil ngrok folder and connect program to your account by ``` $ ./ngrok authtoken {{your_ngrok_token}}``` token can get after sign in to ngrok website
- Last step, type command ```$ ./ngrok http {{Port}} ``` Port must be the same as your server (default in this project is 6000)
- Then you will got the information about your server, please find the Forwarding field and copy ```http link example {{http://5d03e282.ngrok.io}}``` and paste it into line bot account in linewebsite.
