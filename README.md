## Установка
```curl -s https://raw.githubusercontent.com/FutsalShuffle/dctl/master/installer.sh | bash -s``` \
Windows: \
Загрузить ```https://raw.githubusercontent.com/FutsalShuffle/dctl/master/installer.bat``` и исполнить.

Поддерживаемые системы:
1) Linux amd64
2) MacOS amd64
3) MacOS arm64
4) Windows amd64

### Изначальная структура проекта
Проект может быть любой структуры. Dctl исполняется в корне проекта (не src/app). \
Итоговая структура получается: \
- src/app
- ...
- dctl.sh
- .dctl/ - тут все файлы dctl
- docker-compose.yaml
- docker-compose.prod.yaml
- .env
- dctl.yaml

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
  enabled: false #По умолчанию false. Генерация k8 файлов. Пока не реализовано до конца.
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
        image: image-php #Из какого образа запускать тесты, вместе с registry
        build:
          args:
            USER_ID: "$USER_ID"
            GROUP_ID: "$GROUP_ID"
      scripts:
        - composer install
        - cp -u .env.example .env && ln -nf .env ./../.env
        - php artisan test
      allow_failure: true
  deploy: #Джобы для стадии deploy
    - name: dev
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
Раннер должен быть настроен на схему docker in docker. (Дефолтный образ - docker:latest или конкретной версии).
