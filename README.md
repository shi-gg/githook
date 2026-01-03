[![](https://img.shields.io/discord/828676951023550495?color=5865F2&logo=discord&logoColor=white)](https://lunish.nl/support)
![](https://ghcr-badge.egpl.dev/luna-devv/githook/latest_tag)
![](https://ghcr-badge.egpl.dev/luna-devv/githook/size)

[![ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/I3I6AFVAP)

**⚠️ In development, breaking changes ⚠️**

## About
GitHook is a Discord -> GitHub middleware designed to improve GitHub notifications for discord. Basically just a modern and nicer version of the discord built-in GitHub webhook support.

If you need help developing with this, join **[our Discord Server](https://discord.com/invite/yYd6YKHQZH)**.

![image](https://github.com/user-attachments/assets/4a52bba8-b71f-413a-be10-f0a5e626906d)

## Deploy
To deploy this project, create the following `docker-compose.yml`:
```yml
services:
  app:
    image: ghcr.io/luna-devv/githook:latest
    container_name: githook
    ports:
      - "8080:8080"
    restart: unless-stopped
    volumes:
      - ./.env:/app/.env
```

Create a `.env` file with the following values:
```env
REDIS_USR=""
REDIS_PW=""
REDIS_ADDR="127.0.0.1:6379"
SECRET="replace-me-with-a-random-string"
```

To deploy the project, run:
```sh
docker compose up -d
```

## Develope
Clone this repo and start the server with the following commands:
```bash
git clone https://github.com/shi-gg/githook
go run .
```

## Usage
To create a webhook, you have to follow the following steps:
1. Go to [localhost:8080/create?url=<discord-webhook>](http://localhost:8080/create?url=); replace `<discord-webhook>` with an actual webhook url.
2. Go to your repostitories webhook settings.
3. Enter `http://localhost:8080/incomming/<hash>`; replace `<hash>` with the hash returned from step 1 and 2.
4. Set **Content Type** to `application/json`, **Secret** to none, enable *SSL Verification** and enable all events.

![image](https://github.com/user-attachments/assets/20f5d319-482a-418d-a20a-b81ac79cc572)

*Please don't be stupid, localhost won't work with GitHub*
