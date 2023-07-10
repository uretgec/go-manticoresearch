# GO ManticoreSearch

Golang Manticoresearch Http Client

> NOT ready to use in production

## Docker Use

Only run this command. http api localhost:9308 will be ready to use.

```sh
    docker-compose up -d
```

## Recomended Packages

Fiber Client

### NOTE

- This manticoresearch http api is not correct, everthings are wrong wrong wrong and api does not returns standart responses.
- manticore docker file has some bugs. env -> MCL=1 when set only install columnar library install, always returns error and not run container. So default installation env EXTRA=1 must be set, runs http api request through php api and after manticore binary. Returns none standart json or text message. i dont understand why why why.
