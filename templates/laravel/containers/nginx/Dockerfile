FROM zedsh/nginx

ARG USER_ID='1000'
ARG USER_ID=${USER_ID}
ENV USER_ID ${USER_ID}

ARG GROUP_ID='1000'
ARG GROUP_ID=${GROUP_ID}
ENV GROUP_ID ${GROUP_ID}

ARG DOCUMENT_ROOT='web'
ARG DOCUMENT_ROOT=${DOCUMENT_ROOT}
ENV DOCUMENT_ROOT=${DOCUMENT_ROOT}


RUN apk --no-cache add shadow
RUN groupmod -g $USER_ID nginx
RUN usermod -u $USER_ID nginx

#
COPY .dctl/containers/nginx/conf/nginx.conf /etc/nginx/nginx.conf
COPY .dctl/containers/nginx/conf/default.conf /etc/nginx/conf.d/default.conf
RUN sed -i "s/#DOCUMENT_ROOT#/${DOCUMENT_ROOT}/g" /etc/nginx/conf.d/default.conf
COPY .dctl/containers/nginx/conf/upstream.conf /etc/nginx/conf.d/upstream.conf

EXPOSE 80

STOPSIGNAL SIGTERM

CMD ["nginx", "-g", "daemon off;"]
