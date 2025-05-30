### Go-stream-it

This repository hosts my _extremely minimal_ Telegram bot, for managing my shrimps video streamings to YouTube remotely, using Telegram bot commands.

The project has been developed on/for Linux only, tested on a Raspberry Pi 4, running Raspberry Pi OS.

Commands:
- `/start`: start streaming to YouTube
- `/stop`: stop streaming to YouTube
- `/status`: get the streaming status (running/not running)

## Build the bot

1. Build the bot for the target system (here a 64-bit Raspberry Pi):

    ```bash
    GOOS=linux GOARCH=arm64 go build -o go-stream-it
    ```

1. Make it executable:

    ```bash
    chmod +x ./go-stream-it
    ```

Or simply:

```bash
make build
```

## Create a systemd service

Create the systemd service file to describe this service:

```bash
sudo nano /etc/systemd/system/go-stream-it.service
```

Paste the following (edit the <fields>):

```bash
[Unit]
Description=Go Telegram Bot Service
After=network.target

[Service]
Type=simple
User=pagans
Group=pagans
ExecStart=<path-to-executable>
WorkingDirectory=<path-to-executable-folder>
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
```

Now, register the service and start it right away:

```bash
sudo systemctl daemon-reexec
sudo systemctl daemon-reload
sudo systemctl enable go-stream-it.service
sudo systemctl start go-stream-it.service
```

You can always check the service's status by running:

```bash
sudo systemctl status go-stream-it.service
```

More logs:

```bash
journalctl -u go-stream-it.service -f
```

## Run without service

To just run the bot, type the command:

```bash
make run # must run after the build step
```

## Restart service

After code changes, simply re-build the bot and restart the service:

```bash
make build
make deploy
```

## Env variables

| Name    | Used for |
| -------- | ------- |
| YOUTUBE_STREAMING_URL  | Streaming the video content to the right YouTube channel    |
| TELEGRAM_BOT_TOKEN* | Authenticating the Telegram bot creation     |

The environment variables have to be defined in a `.env` file in the root folder of the repo.

*To run this code you should have set up a Telegram bot, and have received a bot token.