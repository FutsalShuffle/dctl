FROM python:3.11.0-alpine
ENV PYTHONDONTWRITEBYTECODE=1
ENV PYTHONUNBUFFERED=1
ARG APP_ENV='prod'
ARG APP_ENV=${APP_ENV}
ENV APP_ENV=${APP_ENV}

ARG USER_ID='1000'
ARG USER_ID=${USER_ID}
ENV USER_ID=${USER_ID}

ARG GROUP_ID='1000'
ARG GROUP_ID=${GROUP_ID}
ENV GROUP_ID=${GROUP_ID}

RUN apk update && apk add bash python3-dev gcc libc-dev ssmtp

WORKDIR /app

RUN pip install --upgrade pip
RUN pip install gunicorn
ADD . ./
RUN ls -la
ADD ./requirements.txt ./
RUN pip install -r requirements.txt
CMD ["gunicorn", "main.wsgi:application", "--bind", "0.0.0.0:8000"]
#CMD ["python", "manage.py", "runserver" ,"0.0.0.0:8000"]
