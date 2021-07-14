# Watchdog

Watchdog is a simple service for monitoring other services

## How to run
Compile and run *watchdog.go* file or use Dockerfile to run the app in docker.

The app expects a config with service list to monitor. By default, *config.yml* is used from the base dir.

You could specify a custom config file in another location by providing it through *-config-path* option:
```yaml
 go run watchdog.go -config-path /custom/path/config.yml
```

## Config structure

```yaml
entities:
  - name: "Example"
    health-url: "https://example.com/health"
    check-period: "1m"
    valid-statuses: [ 200, 204 ]
    fails-allowed: 3 # how many subsequent fails are allowed
notifiers:
  - name: "telegram"
    params:
      token: "<BOT_TOKEN>"
      channel-id: "<CHANNEL_ID>"
  - name: "console"
```

## Notifiers

Notifiers send status messages about monitored services' health

### Console
Print statuses to the console

### Telegram

To add a telegram notifier you will need a channel with a bot added to it as an administrator. You should obtain the bot
API token and the chat id

To do this, follow these steps:

1. Create a bot by writing to @botfather
2. Create a channel
3. Add the created bot to the channel as administrator
4. Obtain the channel id, for example by writing some message to the channel first and then requesting the bot's updates. The
   id can be found in *chat* data

```
https://api.telegram.org/bot<BOT_TOKEN>/getUpdates
```
```json
{
   "ok":true,
   "result":[
      {
         "my_chat_member":{
            "chat":{
               "id": "<CHANNEL_ID>",
               "title":"Service Monitoring",
               "type":"channel"
            }
         }
      }
   ]
}

```
