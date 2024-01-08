## Установка
```curl -s https://raw.githubusercontent.com/FutsalShuffle/dctl/master/installer.sh | bash -s``` \
Windows: \
Загрузить ```https://raw.githubusercontent.com/FutsalShuffle/dctl/master/installer.bat``` и исполнить.

Поддерживаемые системы:
1) Linux amd64
2) MacOS amd64
3) MacOS arm64
4) Windows amd64 (? Возможно не работае self-update)

### Изначальная структура проекта
Сам проект может быть любой структуры. Но конфигурационные файлы должны быть в .dctl (докерфайлы, и прочее). Dctl исполняется в корне проекта (не src/app). \
Итоговая структура получается: \
- src/app
- ...
- dctl.sh
- .dctl/ - тут все файлы dctl
- docker-compose.yaml
- docker-compose.prod.yaml
- .env
- dctl.yaml \

В итоге docker-compose должен смотреть на ./.dctl/containers/... Если указан только image, то dctl автоматически сгенерирует Dockerfile в нужное место.
### Перенос существующего проекта с docker-compose dctl сборкой на новый dctl
В идеале - сделать ```dctl --init {тип}``` в корне проекта. \
Далее нужно перенести Dockerfile, конфиги для них из ./docker в ./dctl/ (если они отличаются от стандартных, к примеру используются доп пакеты). **Обратите внимание** на context. В старых сборках он был указан обычно на ./containers/*/, теперь он на корень проекта. Соответственно ADD директивы и подобные должны быть в соответствии с этим контекстом. \
После этого нужно настроить сам dctl.yaml файл. Прописать нехватающие контейнеры (можно скопировать из ./docker/docker-compose.yaml - но опять же обращайте внимание на build.context и путь до Dockerfile), комманды. \
После настройки нужно прописать ```dctl``` в корне проекта и вы получите готовую сборку.

### Команды
1) ```dctl``` - запуск на текущем каталоге (корень проекта)
2) ```dctl --update``` Обновление
3) ```dctl --init {projectType}``` - инициализация проекта (докерфайлы, базовый dctl.yaml). После исполнения нужно еще сделать ```dctl```. Допустимые значения: laravel, bitrix, symfony, django, next
4) ```dctl --version``` Узнать текущую версию dctl

Команды самого dctl.sh:
1) ```./dctl.sh up (-d)``` - запуск проекта в docker-compose
2) ```./dctl.sh run``` - зайти в run основной контейнер (чаще всего это php, node, python)
3) ```./dctl.sh db``` - зайти в бд контейнер (postgres, mysql)
4) ```./dctl.sh down``` - остановить проект
5) ```./dctl.sh make vendor``` - установить composer завимости
6) ```./dctl.sh make env``` - скопировать .env.example в .env \

Это основные команды, за остальными в сам .dctl

### Конфигурация dctl.yaml
```yaml
name: project # Название проекта, к примеру todolist
docker:
  enabled: true #Включить генерацию docker-compose и остального для него. По умолчанию true
  registry: "" #Docker registry. Для ci/cd. По умолчанию пусто
k8: 
  enabled: false #По умолчанию false. Генерация k8 файлов.
deployments: #Если используется k8.enabled: true
  db: #Название деплоймента (в данном случае бд)
    service: true #Нужен ли сервис (обычно да, чтобы деплоймент был доступен по порту)
    replicas: 1 #Сколько реплик (подов)
    pvc: #Постоянное хранилище (для данных бд)
      - name: db-data
        storage: 512Mi #Размер хранилища
    containers:
      postgres: #Название контейнера
        env: #Такое лучше выносить в секреты, либо в енвы в отдельном файле в .dctl/kube/secrets/{env}.yaml.
          POSTGRES_DB: dev
          POSTGRES_USER: dev
          POSTGRES_PASSWORD: dev
          PGDATA: "/var/lib/postgresql/data/pgdata"
        pvc: #Указываем, что используем pvc из деплоймента в этом контейнере
          - name: db-data #Название из pvc выше
            mountPath: '/var/lib/postgresql/data/pgdata' #Куда монтируется volume (для данных бд)
containers:
  ... #Список контейнеров. По большей мере сюда можно скопировать контент контейнеров из docker-compose.yaml. Повторяет его синтаксис.
commands:
  db:
    vendor: mysql #или postgres
    container: postgres #название контейнера
  run:
    container: php #в каком контейнере исполняется dctl.sh run
  extra: #Доп команды в dctl.sh
    - name: composer-install #название команды
      command: ./dctl.sh run composer install #команда
gitlab: #Настройки для gitlab ci/cd
  only_when: merge_request #Доступные значения: merge_request, always, never, merge_request_master
  cache: 
    paths: #Кеширование путей (к примеру vendor, node_modules)
      - vendor/
  tests: #Джобы на стадии tests
    - name: test
      docker:
        image: project/image-php #Из какого образа запускать тесты, можно без registry урл, можно с тегом и без (без это ветка из ci/cd)
        build:
          args:
            USER_ID: "$USER_ID"
            GROUP_ID: "$GROUP_ID"
      scripts:
        - composer install
        - cp -u .env.example .env && ln -nf .env ./../.env
        - php artisan test
      allow_failure: true
      only:
         - ...
  deploy: #Джобы для стадии deploy
    - name: dev
      only:
         - ...
      scripts: #Примерная команда для деплоя на dev через ssh (зайти, переключить ветку на текущий MR, сделать пулл и перезапустить приложение)
        - mkdir ~/.ssh
        - echo -n "$CI_GITLAB_PRIVATE_KEY" | base64 -d > ~/.ssh/id_rsa
        - chmod 600 ~/.ssh/id_rsa
        - echo -e "Host *\n\tStrictHostKeyChecking no\n\n" > ~/.ssh/config
        - eval $(ssh-agent -s)
        - ssh-add ~/.ssh/id_rsa
        - ssh $CI_SSH_CONNECT "cd /home/user/project; git fetch && git checkout $CI_COMMIT_REF_NAME; git pull origin $CI_COMMIT_REF_NAME; sh ./dctl.sh down; sh ./dctl.sh up -d"
```

### CI/CD
Сборкой образов в registry занимается сам dctl. Если указан registry и gitlab.tests или gitlab.deploy, то автоматически создастся .gitlab-ci.yaml с build стейджем и джобами на все контейнеры. \
Сборка происходит в ```{registry}/{projectName}/{containerName}:$CI_COMMIT_REF_NAME```. \
Запуская build-docker и push-docker в CI мы получим собранные и запушенные образы с тегом текущей ветки. Если же не из CI - latest. \
build-docker-prod и push-docker-prod соберет и запушит контейнеры с тегом prod-latest. Их можно использовать для прода в docker-compose.prod.yaml, или же кубере (?) \
Раннер должен быть настроен на схему docker in docker. (Дефолтный образ - docker:latest или конкретной версии).

### Кубы
Контейнеры по дефолту берутся с тегом `{{env}}-latest`. Для билда и пуша в registry можно использовать `./dctl.sh build-docker-{{env}} && ./dctl.sh push-docker-{{env}}` \
Файлы конфигураций генерируются с числовым префиксом, чтобы при `apply -f` на директорию они создавались в корректном порядке. \
Как правило, штуки наподобие бд редиса, воркеры - это отдельный деплоймент, который не входит в основной с приложением (куда обычно будет входить приложение, nginx). \
При этом для воркеров или подобных скриптов не нужно указывать service, т.к им не нужен выходной порт (как правило). \
Связь внутри одного деплоймента - по названию контейнера, как в docker-compose. При этом связь между деплойментами - по названию сервиса из *-service.yaml. При генерации это **{НазваниеПроекта}-{НазваниеКонтейнера}**. \
Если нам нужно шарить статику условно из django/php в nginx (какие-нибудь ассеты для фронта), то можно сбилдить их в контейнер, либо сделать emptyDir. Пример:
```
...
deployments:
    app:
        emptyDir:
            enabled: true
            sizeLimit: 512Mi
        containers:
            php:
                emptyDir:
                    enabled: true
                    mountPath: /var/www/html/public
            nginx:
                emptyDir:
                    enabled: true
                    mountPath: /var/www/html/public
```
И билдить их при старте контейнера. Но это не оптимальный вариант, если речь про загружаемые юзерами файлы, которые нужно сохранять. Тогда нужно задуматься о s3, или же о PVC.

Ingress должен быть уникальным. Если на машине несколько приложений, то нужно добавлять к пути какой-нибудь уникальный префикс. \
#### Секреты в кубах
Можно прописать k8.useSealedSecrets: true, тогда вместо Opaque Secrets (base64) будет темплейт для Sealed Secrets. Такие секреты можно безопасно хранить в репозитории, т.к создаются и расшифровываются только по сертификату, который находится на сервере.
#### Площадки в кубах
По дефолту создается dev и prod окружения. Если у вас их больше, можно прописать в k8.environments (array) список площадок. К примеру:
```yaml
k8:
  environments:
    - dev
    - prod
    - stage
    - test
```
