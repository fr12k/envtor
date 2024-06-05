# envtor (Environment Injector)

Tool to inject environments variables into a docker-compose.yaml file.

## Use Case

With docker its easy to pass environment variables to a container using the the following syntax.

```bash
docker run --env-file <(printf "ENVIRONMENT1=hello\nENVIRONMENT2=world") alpine env
```

or if you familiar with [teller](https://github.com/tellerops/teller) you can use the following syntax to inject secrets with teller to a container
without saving the secrets to disk.

```bash
docker run --env-file <(teller env) alpine env
```

However, when using docker-compose its not possible to use the same syntax to inject environment variables into a container.

```bash
docker-compose --env-file <(printf "ENVIRONMENT1=hello\nENVIRONMENT2=world") up
```

will inject the environment variables into the docker-compose command but not into the containers defined in the docker-compose file.

## Installation

```bash
go install github.com/fr12k/envtor@latest
```

## Usage

Just by defining the following environment variables `ALL_ENV_VARS` in the docker-compose file the tool will inject all the environment variables piped to it into the docker-compose file and replace the ALL_ENV_VARS with those environment variables.

```yaml
version: "3"
services:
  "1":
    image: alpine
    command: sh -c "env | sort"
    environment:
      - ALL_ENV_VARS
      - GO_ENV=production
```

The above docker-compose file will be modified to the following docker-compose file printed to stdout.

```bash
printf "ENVIRONMENT1=hello\nENVIRONMENT2=world" | build/envtor
services:
    "1":
        command: sh -c "env | sort"
        environment:
            - ENVIRONMENT1=hello
            - ENVIRONMENT2=world
            - GO_ENV=production
        image: alpine
version: "3"
```

It could then be stored to disk or even directly piped to `docker-compose -f - up` to start the containers based on the
modified docker-compose file.

```bash
printf "ENVIRONMENT1=hello\nENVIRONMENT2=world" | build/envtor | docker-compose -f - up

1-1  | ENVIRONMENT1=hello
1-1  | ENVIRONMENT2=world
1-1  | GO_ENV=production
1-1  | HOME=/root
1-1  | HOSTNAME=9cb840a5c12e
1-1  | PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
1-1  | PWD=/
1-1  | SHLVL=1
1-1  | _=/bin/printenv
1-1  | _BASH_BASELINE=5.2.21
1-1  | _BASH_BASELINE_PATCH=21
1-1  | _BASH_LATEST_PATCH=26
1-1  | _BASH_VERSION=5.2.26
2-1  | ENV=test
2-1  | ENVIRONMENT1=hello
2-1  | ENVIRONMENT2=world
2-1  | GO_ENV=production
2-1  | HOME=/root
2-1  | HOSTNAME=0d68411afc9a
2-1  | PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
2-1  | PWD=/
2-1  | SHLVL=1
2-1  | _=/bin/printenv
2-1  | _BASH_BASELINE=5.2.21
2-1  | _BASH_BASELINE_PATCH=21
2-1  | _BASH_LATEST_PATCH=26
2-1  | _BASH_VERSION=5.2.26
1-1 exited with code 0
2-1 exited with code 0
```

As you can see the environment variables `ENVIRONMENT1=hello` and `ENVIRONMENT2=world` are injected into the docker-compose file
without modifying the original file. Also the modified docker-compose file is piped into `docker-compose -f - up` to start the containers
and not saved to disk.

The original docker-compose file can also be used without the injected environment variables.

```bash
docker-compose up

[+] Running 2/2
 ✔ Container envtor-2-1  Recreated                                                                                               0.1s
 ✔ Container envtor-1-1  Recreated                                                                                               0.1s
Attaching to 1-1, 2-1
2-1  | ENV=test
2-1  | GO_ENV=production
2-1  | HOME=/root
2-1  | HOSTNAME=afcbad5fce0d
2-1  | PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
2-1  | PWD=/
2-1  | SHLVL=1
2-1  | _=/bin/printenv
2-1  | _BASH_BASELINE=5.2.21
2-1  | _BASH_BASELINE_PATCH=21
2-1  | _BASH_LATEST_PATCH=26
2-1  | _BASH_VERSION=5.2.26
1-1  | GO_ENV=production
1-1  | HOME=/root
1-1  | HOSTNAME=2964f2f4e932
1-1  | PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
1-1  | PWD=/
1-1  | SHLVL=1
1-1  | _=/bin/printenv
1-1  | _BASH_BASELINE=5.2.21
1-1  | _BASH_BASELINE_PATCH=21
1-1  | _BASH_LATEST_PATCH=26
1-1  | _BASH_VERSION=5.2.26
2-1 exited with code 0
1-1 exited with code 0
```
