Run your own Goni server with Docker

## Before you start
[Create 'Slack App'](https://api.slack.com/apps), and change `GONI_SLACK_CLIENT`, `GONI_SLACK_SECRET` with your APIKEY in `dashboard/dockerfile`.

## Usage
``` bash
git clone https://github.com/monitflux/goni
cd ./docker
docker-compose up
```
