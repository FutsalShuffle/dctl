#!/bin/bash
set -e
#first cd current dir
cd "$(dirname "${BASH_SOURCE[0]}")"

export DEFAULT_USER="1000";
export DEFAULT_GROUP="1000";

export USER_ID=`id -u`
export GROUP_ID=`id -g`
export USER=$USER

if [ "$USER_ID" == "0" ];
then
    export USER_ID=$DEFAULT_USER;
fi

if [ "$GROUP_ID" == "0" ];
then
    export GROUP_ID=$DEFAULT_GROUP;
fi

test -e "./.env" || { cp .env.example .env; };
#load .env
export $(egrep -v '^#' .env | xargs)
export PROJECT_PREFIX={{.Name}}

{{$projectName := .Name}}{{$gitlab := .Gitlab}}{{$docker := .Docker}}
if [ $# -eq 0 ]
  then
    echo "HELP:"
    echo "make env - copy .env.example to .env"
    echo "make db - load init database dump to mysql"
    echo "db import FILE - load FILE to mysql"
    echo "db renew - load dump from repo, fresh db and apply"
    echo "build - make docker build"
    echo "up - docker up in console"
    echo "up -d - docker up daemon"
    echo "down - docker down"
    echo "run - run in main container from project root"
    echo "build-docker - build containers with ci or latest tags"
    echo "push-docker - push containers with ci or latest tags to registry"
    echo "build-docker-prod - build containers with prod-latest tag"
    echo "push-docker-prod - push containers with with prod-latest tag"
    echo "db - enter database container"
fi

if [ "$1" == "make" ];
  then
    if [ "$2" == "env" ];
        then
            cp .env.example .env
    fi
    if [ "$2" == "db" ];
        then
         applyDump "../bitrix/database/init.sql";
    fi
    if [ "$2" == "vendor" ];
        then
           docker-compose -p {{$projectName}} run php composer install
           docker-compose -p {{$projectName}} run php npm i
    fi
fi
{{ if eq .Commands.Db.Vendor "mysql" }}
function applyDump {
    cat $1 | docker exec -i {{$projectName}}_{{.Commands.Db.Container}} mysql -u $MYSQL_USER -p"$MYSQL_PASSWORD" $MYSQL_DATABASE;
    return $?
}
{{end}}
{{ if .Commands.Db.Vendor }}
if [ "$1" == "db" ];
  then

    {{ if eq .Commands.Db.Vendor "mysql" }}
    if [ "$2" == "" ];
        then
        docker exec -it {{$projectName}}_{{.Commands.Db.Container}} mysql -u $MYSQL_USER -p"$MYSQL_PASSWORD" $MYSQL_DATABASE;
    fi

    if [ "$2" == "import" ];
        then
        applyDump $3
    fi

    if [ "$2" == "export" ];
        then
        docker exec -it {{$projectName}}_{{.Commands.Db.Container}} su mysql -c "export MYSQL_PWD='$MYSQL_PASSWORD'; mysqldump -u $MYSQL_USER -p"$MYSQL_PASSWORD" $MYSQL_DATABASE"
    fi

    if [ "$2" == "renew" ];
        then
        rm -rf "../docker/data/mysql/dump" || echo "old dump not found"
        git clone $DATABASE_REPO ../docker/data/mysql/dump
        applyDump "../docker/containers/mysql/drop_all_tables.sql"
        applyDump "../docker/data/mysql/dump/database.sql"
    fi {{end}}

    {{- if eq .Commands.Db.Vendor "postgres" }}
    if [ "$2" == "" ];
        then
        docker exec -it {{$projectName}}_{{.Commands.Db.Container}} psql -U $DB_USER $DB_NAME;
    fi
    if [ "$2" == "export" ];
        then
        docker exec -it {{$projectName}}_{{.Commands.Db.Container}} pg_dump -c -v -f ./dump.sql -U $DB_USER $DB_NAME
    fi
    if [ "$2" == "import" ];
        then
        cat $1 | docker exec -i {{$projectName}}_{{.Commands.Db.Container}} psql -U $DB_USER $DB_NAME;
    fi
    {{end}}
fi {{end}}

if [ "$1" == "build" ];
  then
    docker-compose build
fi

if [ "$1" == "up" ];
  then
    if [ "$2" == "-d" ];
        then
            docker-compose -p {{$projectName}} up -d;
        else
            docker-compose -p {{$projectName}} up
    fi
fi

if [ "$1" == "down" ];
  then
    docker-compose -p {{$projectName}} down
fi

if [ "$1" == "fulldown" ];
  then
    docker-compose -p {{$projectName}} down --rmi local
fi

{{ if .Commands.Run.Container }}
function runInContainer {
    local command=$@
    echo $command;
    docker exec -i {{$projectName}}_{{.Commands.Run.Container}} bash -c "cd /var/www/html/;$command"
    return $?
}
if [ "$1" == "run" ];
  then
    if [ "$2" == "" ];
        then
        docker exec -u www-data -it {{$projectName}}_{{.Commands.Run.Container}} bash
        else
        runInContainer "${@:2}"
    fi
fi {{end}}

if [ "$1" == "build-docker" ];
  then
    {{range $index, $container := .Containers}}
    if [ "$2" == "{{$index}}" ];
        then
            docker build {{$container.Build.Context}} \
            --file {{$container.Build.Context}}/{{$container.Build.Dockerfile}} \
            {{range $argName, $argVal := $container.Build.Args}}--build-arg {{$argName}}={{$argVal}} \
            {{end}}--build-arg USER_ID=$USER_ID \
            --build-arg GROUP_ID=$GROUP_ID \
            $(if [ -n "${CI}" ]; then echo "--tag {{if $docker.Registry}}{{$docker.Registry}}/{{end}}{{$projectName}}/{{$index}}:${CI_COMMIT_REF_NAME}" ; fi) \
            -t {{if $docker.Registry}}{{$docker.Registry}}/{{end}}{{$projectName}}/{{$index}}:latest;
    fi
    {{end}}
    if [ "$2" == "" ];
        then
          cd "$(dirname "${BASH_SOURCE[0]}")"{{range $index, $c := .Containers}}
          ./dctl.sh build-docker {{$index}}{{end}}
    fi
fi

if [ "$1" == "push-docker" ];
  then
    {{range $index, $container := .Containers}}
    if [ "$2" == "{{$index}}" ];
        then
          if [ -n "${CI}" ];
          then
            docker push {{if $docker.Registry}}{{$docker.Registry}}/{{end}}{{$projectName}}/{{$index}}:${CI_COMMIT_REF_NAME}
          else
            docker push {{if $docker.Registry}}{{$docker.Registry}}/{{end}}{{$projectName}}/{{$index}}:latest
          fi
    fi
    {{end}}
    if [ "$2" == "" ];
        then
          cd "$(dirname "${BASH_SOURCE[0]}")"{{range $index, $c := .Containers}}
          ./dctl.sh push-docker {{$index}}{{end}}
    fi
fi
{{$containers := .Containers}}
{{range $index, $environment := .K8.Environments}}
if [ "$1" == "push-docker-{{$environment}}" ];
  then
    {{range $index, $container := $containers}}
    if [ "$2" == "{{$index}}" ];
        then
          docker push {{if $docker.Registry}}{{$docker.Registry}}/{{end}}{{$projectName}}/{{$index}}:{{$environment}}-latest
    fi
    {{end}}
    if [ "$2" == "" ];
        then
          cd "$(dirname "${BASH_SOURCE[0]}")"{{range $index, $c := $containers}}
          ./dctl.sh push-docker-{{$environment}} {{$index}}{{end}}
    fi
fi

if [ "$1" == "build-docker-{{$environment}}" ];
  then
    {{range $index, $container := $containers}}
    if [ "$2" == "{{$index}}" ];
        then
            docker build {{$container.Build.Context}} \
            --file {{$container.Build.Context}}/{{$container.Build.Dockerfile}} \
            {{range $argName, $argVal := $container.Build.Args}}--build-arg {{$argName}}={{$argVal}} \
            {{end}}--build-arg USER_ID=$USER_ID \
            --build-arg GROUP_ID=$GROUP_ID \
            -t {{if $docker.Registry}}{{$docker.Registry}}/{{end}}{{$projectName}}/{{$index}}:{{$environment}}-latest;
    fi
    {{end}}
    if [ "$2" == "" ];
        then
          cd "$(dirname "${BASH_SOURCE[0]}")"{{range $index, $c := $containers}}
          ./dctl.sh build-docker-{{$environment}} {{$index}}{{end}}
    fi
fi
{{- end -}}

{{range $index, $command := .Commands.Extra}}
if [ "$1" == "{{$command.Name}}" ];
  then
    {{$command.Command}}
fi
{{end}}
