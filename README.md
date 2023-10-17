# Waffler Server
[![en](https://img.shields.io/badge/lang-en-red.svg)](README.md)
[![en](https://img.shields.io/badge/lang-ru-blue.svg)](README.rus.md)

This Github repository is the backend of the Waffler API.

Waffler is a service that uses AI to process information sources for racism and contradictions.

Waffler currently supports telegram and youtube sources.

You can find more information in the [API description](REST_Rus.md).


## Gettiing Started

### Env File

A sample .env file can be found [here](deployments/.env.example).

You will need to make your own .env file in the root of this repository.

#### Telegram API ID

Notice the TELEGRAM_* variables in the example .env file.

To use Waffler you will need your own [Telegram API id](https://core.telegram.org/api/obtaining_api_id).

Once you have updated your .env file you will need to log in once.

To do that build the app and run it once.

```
go mod download

CGO_ENABLED=0 go build -a -installsuffix cgo -o waffler ./cmd/api

./waffler
```

You will be prompted to enter a code that you received through Telegram.

After that, your token will be stored locally in ./telegram_sessions and you won't have to login again.

#### Chat GPT Token

To process sources you will need a [ChatGPT Token](https://platform.openai.com/signup).

### Starting Docker

To setup a local version of Waffler you will need to install [Docker](https://docs.docker.com/get-docker/).

After installation start the containers with the command:

```
docker compose -f ./deployments/Docker-compose_localhost.yaml --env-file ./.env up
```