FROM php:8.2.13-fpm-alpine3.18

RUN mkdir -p /var/run/mysqld
RUN mkdir -p /run/php

ARG APP_ENV='prod'
ARG APP_ENV=${APP_ENV}
ENV APP_ENV=${APP_ENV}

ARG USER_ID='1000'
ARG USER_ID=${USER_ID}
ENV USER_ID=${USER_ID}

ARG GROUP_ID='1000'
ARG GROUP_ID=${GROUP_ID}
ENV GROUP_ID=${GROUP_ID}

RUN apk add ssmtp bash git

RUN wget -O /usr/local/bin/install-php-extensions https://github.com/mlocati/docker-php-extension-installer/releases/latest/download/install-php-extensions
RUN chmod +x /usr/local/bin/install-php-extensions

RUN install-php-extensions pdo_pgsql intl gd zip mcrypt imagick memcached memcache opcache

ADD .dctl/containers/php/conf/php.ini /usr/local/etc/php/php.ini
ADD .dctl/containers/php/conf/php.ini /usr/local/etc/php/php.ini
ADD .dctl/containers/php/conf/www.conf /usr/local/etc/php-fpm.d/www.conf

ADD . /var/www/html
RUN sed -i -e "s/www-data:x:82:82:Linux User,,,:\/home\/www-data:\/sbin\/nologin/www-data:x:${USER_ID}:${GROUP_ID}:Linux User,,,:\/home\/www-data:\/bin\/bash/g" /etc/passwd
RUN sed -i -e "s/www-data:x:82:www-data/www-data:x:${GROUP_ID}:www-data/g" /etc/group

RUN apk add bash

RUN cd /bin && curl --silent --show-error https://getcomposer.org/installer | php && mv composer.phar composer && chmod 777 composer

WORKDIR /var/www/html
RUN chown -R www-data:www-data /var/www/html
COPY .dctl/containers/php/entrypoint.sh /entrypoint.sh

CMD ["/entrypoint.sh"]

EXPOSE 9000
