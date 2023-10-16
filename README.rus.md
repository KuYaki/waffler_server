# Вафлер Сервер
[![en](https://img.shields.io/badge/lang-en-red.svg)](README.md)

Этот репозиторий Github является серверной частью Waffler API.

Waffler — сервис, который использует ИИ для обработки источников информации на предмет расизма и противоречий.

В настоящее время Waffler поддерживает Telegram и YouTube.

Более подробную информацию вы можете найти в [описании API](REST_Rus.md).


## Начало Работы

### Файл Env 

Пример файла .env можно найти [здесь](deployments/.env.example).

Вам нужно будет создать свой собственный файл .env в корне этого репозитория.

#### Telegram API ID

Обратите внимание на переменные TELEGRAM_* в примере файла .env.

Чтобы использовать Waffler, вам понадобится собственный [идентификатор API Telegram](https://core.telegram.org/api/obtaining_api_id).

После обновления файла .env вам нужно будет войти в систему.

Для этого соберите приложение и запустите его.

```
go mod download

CGO_ENABLED=0 go build -a -installsuffix cgo -o waffler ./cmd/api

./waffler
```

Вам будет предложено ввести код, который вы получите через Telegram.

После этого ваш токен будет храниться локально в ./telegram_sessions, и вам не придется повторно входить в систему.

#### Токен Chat GPT

Для обработки источников вам понадобится [Токен ChatGPT](https://platform.openai.com/signup).

### Запуск Docker

Чтобы запустить локальную версию Waffler, вам понадобиться установить [Docker](https://docs.docker.com/get-docker/).

После установки запустите контейнеры командой:

```
docker compose -f ./deployments/Docker-compose_localhost.yaml --env-file ./.env up
```