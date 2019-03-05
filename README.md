# Smartest Telegram Bot

Telegram bot built in go.

A lot of code pulled from [here](https://github.com/go-telegram-bot-api/telegram-bot-api/blob/13c54dc548f7ca692fe434d4b7cac072b0de0e0b/types.go#L129).

## Development + Deployment

To spin up a container with hot reloading, make your own file *my.env* and place it the root of this project. Put in it these variables:
```
TELE_KEY=<key given from fatherbot>
REDDIT_CLIENT_ID=<id of app generated on reddit>
REDDIT_CLIENT_SECRET=<secret of the app>
REDDIT_USERNAME=<your reddit username>
REDDIT_PASSWORD=<your reddit password>
```

To actually run the app:
```
docker-compose up
```

### Database Schema

```SQL
CREATE TABLE `messages` (
  `MessageID` int(11) NOT NULL,
  `ChatID` bigint(20) NOT NULL,
  `Date` int(11) NOT NULL,
  `FromID` int(11) DEFAULT NULL,
  `FromUserName` varchar(100) DEFAULT NULL,
  `ReplyToMessageID` int(11) DEFAULT NULL,
  `ForwardedFromUserID` int(11) DEFAULT NULL,
  `ForwardedFromChatID` bigint(20) DEFAULT NULL,
  `PhotoFileID` varchar(200) DEFAULT NULL,
  `VideoFileID` varchar(200) DEFAULT NULL,
  `DocumentFileID` varchar(200) DEFAULT NULL,
  `StickerID` varchar(200) DEFAULT NULL,
  `Text` varchar(4096) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

```
