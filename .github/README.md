# Echoes

Self-hosted blog engine

![home page](home-page.png)

![blog page](blog-page.png)

## Installation

### Docker

```bash
docker build -t echoes .

docker run --name echoes \
    -p 5000:80 \
    -e ECHOES_POSTGRES_ADDR=localhost:5432 \
    -e ECHOES_POSTGRES_DB=echoesdb \
    -e ECHOES_POSTGRES_SSL_MODE=disable \
    -e ECHOES_POSTGRES_USER=user \
    -e ECHOES_POSTGRES_PASS=1234 \
    -e ECHOES_REDIS_ADDR=localhost:6379 \
    -v $(PWD)/config.yaml:/app/config.yaml \
    -v $(PWD)/colorscheme.css:/app/assets/css/colorscheme.css \
    -v $(PWD)/images:/app/assets/images \
    -d echoes
```

### Docker Compose

Example docker-compose.yaml

```yaml
services:
  web:
    container_name: echoes-web
    build: .
    restart: always
    ports:
      - "5000:80"
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
      migrate:
        condition: service_completed_successfully
    environment:
      - ECHOES_POSTGRES_ADDR=echoes-postgres:5432
      - ECHOES_POSTGRES_DB=echoesdb
      - ECHOES_POSTGRES_SSL_MODE=disable
      - ECHOES_POSTGRES_USER=user
      - ECHOES_POSTGRES_PASS=1234
      - ECHOES_REDIS_ADDR=echoes-redis:6379
    networks:
      - mainnet
    healthcheck:
      test: curl --fail http://localhost:80/health || exit 1
      interval: 1m
      timeout: 30s
      retries: 5
      start_period: 30s
    volumes:
      - $(PWD)/config.yaml:/app/config.yaml
      - $(PWD)/colorscheme.css:/app/assets/css/colorscheme.css
      - $(PWD)/images:/app/assets/images
  postgres:
    container_name: echoes-postgres
    image: postgres
    restart: always
    environment:
      - POSTGRES_USER=user
      - POSTGRES_PASSWORD=1234
      - POSTGRES_DB=echoesdb  
    volumes:
      - pg-data:/var/lib/postgresql/data
    networks:
      - mainnet
    healthcheck:
      test: [ "CMD-SHELL", "pg_isready" ]
      interval: 10s
      timeout: 5s
      retries: 3
  migrate:
    container_name: echoes-migrate
    image: migrate/migrate
    networks:
      - mainnet
    volumes:
      - ./migrations:/migrations
    command:
      [
        "-path",
        "/migrations",
        "-database",
        "postgres://user:1234@echoes-postgres:5432/echoesdb?sslmode=disable",
        "up"
      ]
    depends_on:
      postgres:
        condition: service_healthy
    links:
      - postgres
  redis:
    container_name: echoes-redis
    image: redis
    restart: always
    networks:
      - mainnet
    volumes:
      - cache:/data
    healthcheck:
      test: [ "CMD", "redis-cli", "ping" ]
      interval: 10s
      timeout: 5s
      retries: 3

networks:
  mainnet:
    driver: bridge

volumes:
  pg-data:
  cache:
```

## Configuration

It can be configured using either config.yaml or environment variables

Example config.yaml

```yaml
server:
  addr: "0.0.0.0:80"
  session_key: "super_secret_session_key"
  root_pass: "root"
redis:
  addr: "localhost:6379"
  db: 0
  password: ""
postgres:
  username: "user"
  password: "1234"
  ssl_mode: "disable"
  addr: "localhost:5432/echoesdb"
profile:
  name: "Name Surname"
  bio: >
    Hi there!
    Add your bio in config.yaml file.
  picture: "/assets/images/icon.svg"
feed:
  title: "microblog recent posts"
  desc: "30 latest posts from microblog"
  link: "https://website.com/blog"
  detail_link: "https://website.com/posts/"
  author: "Name Surname"
  email: "address@email.com"
website:
  title: "echoes"
  logo: "/assets/images/icon.svg"
  bg_img: "/assets/images/bg.webp"
```

### Changing colorscheme

You can change colorscheme in assets/css/colorscheme.css

Example colorscheme.css

```css
/* Catppuccin Latte */
:root {
    --bg0: rgba(220, 224, 232, 0);
    --bg1: #e6e9ef;
    --bg2: #ccd0da;
    --bg3: #bcc0cc;

    --text: #4c4f69;
    --text-alt: #eff1f5;
    --primary: #11111b;
    --warning: #fe640b;
    --danger: #e64553;
    --success: #40a02b;

    --badge-fg: #191724;
    --badge-bg: #ea9d34;

    --announce-bg: #ebf0f8;
    --announce-fg: #4c4f69;

    --links-bg: #eff1f5;
    --links-fg: #000000;

    --card-bg: #eff1f5;
    --card-fg: #4c4f69;

    --input-bg: #e6e9ef;
    --input-fg: #4c4f69;

    --code-bg: #dce0e8;
    --code-fg: #4c4f69;

    --pin-fg: #fe640b;

    --bg-img-cover: rgba(220, 224, 232, 0.6);
}
```