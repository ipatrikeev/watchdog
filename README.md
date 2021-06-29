# Watchdog

Watchdog is a simple service for monitoring other services

## Notifiers

Notifiers send status messages about monitored services' health

Currently, only *Telegram* notifier is supported

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

## Config structure

```yaml
entities:
  - name: "Example"
    health-url: "https://example.com/health"
    check-period: "1m"
    valid-statuses: [ 200, 204 ]
senders:
  - name: "telegram"
    params:
      token: "<BOT_TOKEN>"
      channel-id: "<CHANNEL_ID>"
```
