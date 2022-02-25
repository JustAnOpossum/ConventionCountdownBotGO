# Convention Countdown Bot
This bot will count down the days until a specified convention with a image and caption.

# Where to find the bot?
The bots that I am currently running with the codebase are below

https://telegram.me/FurFestBot

# How to run the bot?

## Requirements
* GO >= 1.11
* MongoDB > 4


```shell
git clone https://github.com/NerdyRedPanda/ConventionCountdownBotGO
cd ConventionCountdownBotGO
go build
cp config.json.example config.json
mkdir ConventionName
MODE=longPoll DATADIR="$PWD/ConventionName" ./conCountdownBot
```

# DATADIR
This environment variable is used to set the data directory for the bot. You must set this before the bot is started on any of the modes. The only file that needs to be placed inside the data dir is config.json. All other sub directories will be created by the bot.

# Font for image
You must set a custom font for the image generation. Place a font in TTF format in the data directory and then set the name of the font in config.json

# Bot Modes
There are a couple of modes that the bot can run to accomplish various tasks. To change the mode set the environment variable **MODE**, ex. ```MODE=longPoll ./conCountdownBot ```

## longPoll
This mode starts the bot in long polling mode in telegram. This is the option if you don't want webhook support.

## webhook
This starts the bot in webhook mode. To see more about what you need for a webhook setup see the telegram API reference https://core.telegram.org/bots/webhooks

Once you have a working webhook set the **webhookURL** and **webhookPort** values in config.json.

## send
This puts the bot in image sending mode. This mode looks at what users are subscribed on telegram and sends the image to them. If the Twitter keys are filled out in config.json the bot will post a tweet to that account.

## upload
This mode puts the bot into upload mode. This mode guides you through the process of uploading images to the bot. If you want to upload images first zip them and run the bot in upload mode. This mode will ask for the path to the zip file among other options to get your photos uploaded to the bot.

# Twitter Setup
Twitter setup is not currently supported with v2 of the Twitter API still being so new and not allowing twitter updates. If you have a V1.1 key then you can fill in the Twitter options in config.json to send a daily Twitter image.

# Example Picture

![img](https://image.ibb.co/gUan7R/photo_2018_01_09_15_58_11.jpg)