FROM node:18.7.0-alpine3.16

ARG APP_ENV='prod'
ARG APP_ENV=${APP_ENV}
ENV APP_ENV=${APP_ENV}

ARG USER_ID='1000'
ARG USER_ID=${USER_ID}
ENV USER_ID=${USER_ID}

ARG GROUP_ID='1000'
ARG GROUP_ID=${GROUP_ID}
ENV GROUP_ID=${GROUP_ID}
ADD . /var/www/html

RUN npm install -g pm2 gulp
COPY .dctl/containers/node/entrypoint.sh /entrypoint.sh
CMD ["/entrypoint.sh"]

WORKDIR /var/www/html
EXPOSE 9000
