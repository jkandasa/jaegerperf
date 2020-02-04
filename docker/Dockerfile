FROM alpine:3.11

LABEL maintainer="Jeeva Kandasamy <jkandasa@redhat.com>"

ENV APP_HOME /app/

EXPOSE 8080

# copy WEB files
COPY ./web $APP_HOME/web

# copy application bin file
COPY ./appbin $APP_HOME/jaegerperf

WORKDIR $APP_HOME

CMD /app/jaegerperf