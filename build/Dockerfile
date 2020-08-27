FROM alpine:3.12

LABEL maintainer="Jeeva Kandasamy <jkandasa@redhat.com>"

ENV APP_HOME /app/

EXPOSE 8080

# copy WEB files
COPY ./web $APP_HOME/web

# copy application bin file
COPY ./jaegerperf $APP_HOME/jaegerperf

# copy default files
COPY ./resources $APP_HOME/resources

WORKDIR $APP_HOME

CMD /app/jaegerperf