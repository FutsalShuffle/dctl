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
export PROJECT_PREFIX=example


if [ $# -eq 0 ]
  then
    echo "HELP:"
    echo "make env - copy .env.example to .env"
    echo "make db - load init bitrix database dump to mysql"
    echo "db import FILE - load FILE to mysql"
    echo "db renew - load dump from repo, fresh db and apply"
    echo "build - make docker build"
    echo "up - docker up in console"
    echo "up silent - docker up daemon"
    echo "down - docker down"
    echo "run - run in main container from project root"
fi

function runInRabbitMq {
    local command=$@
    echo $command;
    docker exec -i example_rabbitmq bash -c "$command"
    return $?
}

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
           docker-compose -p example run php composer install
           docker-compose -p example run php npm i
    fi
fi


function applyDump {
    cat $1 | docker exec -i example_postgres mysql -u $MYSQL_USER -p"$MYSQL_PASSWORD" $MYSQL_DATABASE;
    return $?
}



if [ "$1" == "db" ];
  then

    
    if [ "$2" == "" ];
        then
        docker exec -it example_postgres mysql -u $MYSQL_USER -p"$MYSQL_PASSWORD" $MYSQL_DATABASE;
    fi

    if [ "$2" == "import" ];
        then
        applyDump $3
    fi

    if [ "$2" == "export" ];
        then
        docker exec -it example_postgres su mysql -c "mysqldump -u $MYSQL_USER -p"$MYSQL_PASSWORD" $MYSQL_DATABASE"
    fi

    if [ "$2" == "renew" ];
        then
        rm -rf "../docker/data/mysql/dump" || echo "old dump not found"
        git clone $DATABASE_REPO ../docker/data/mysql/dump
        applyDump "../docker/containers/mysql/drop_all_tables.sql"
        applyDump "../docker/data/mysql/dump/database.sql"
    fi 
fi 

if [ "$1" == "build" ];
  then
    docker-compose build
fi

if [ "$1" == "up" ];
  then
    if [ "$2" == "silent" ];
        then
            docker-compose -p example up -d;
        else
            docker-compose -p example up
    fi
fi

if [ "$1" == "down" ];
  then
    docker-compose -p example down
fi

if [ "$1" == "fulldown" ];
  then
    docker-compose -p example down --rmi local
fi


function runInContainer {
    local command=$@
    echo $command;
    docker exec -i example_php bash -c "cd /var/www/html/;$command"
    return $?
}
if [ "$1" == "run" ];
  then
    if [ "$2" == "" ];
        then
        docker exec -u www-data -it example_php bash
        else
        runInContainer "${@:2}"
    fi
fi 

if [ "$1" == "rabbitmq" ];
  then

    if [ "$2" == "up" ];
        then
            runInRabbitMq "rabbitmqctl delete_user guest"
            runInRabbitMq "rabbitmqctl add_vhost /"
            runInRabbitMq "rabbitmqctl add_user $RABBITMQ_LOGIN $RABBITMQ_PASSWORD"
            runInRabbitMq "rabbitmqctl set_user_tags $RABBITMQ_LOGIN administrator"
            runInRabbitMq "rabbitmqctl set_permissions -p / $RABBITMQ_LOGIN '.*' '.*' '.*'"
    fi

fi

if [ "$1" == "build" ];
  then
    
    if [ "$2" == "nginx" ];
        then
            docker build -t example/nginx:latest -f ./Dockerfile ./containers/nginx;
    fi
    
    if [ "$2" == "php" ];
        then
            docker build -t example/php:latest -f ./Dockerfile ./containers/php;
    fi
    
    if [ "$2" == "postgres" ];
        then
            docker build -t example/postgres:latest -f ./Dockerfile ./containers/postgres;
    fi
    
    if [ "$2" == "redis" ];
        then
            docker build -t example/redis:latest -f ./containers/redis/Dockerfile ./../;
    fi
    
    if [ "$2" == "smtp" ];
        then
            docker build -t example/smtp:latest -f ./containers/smtp/Dockerfile ./../;
    fi
    
    if [ "$2" == "supervisor" ];
        then
            docker build -t example/supervisor:latest -f ./WorkerDockerfile ./containers/php;
    fi
    
fi
