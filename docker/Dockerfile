FROM alpine:3.11

ENV APP_HOME /app/

EXPOSE 8080

# copy go build file
COPY ./appbin $APP_HOME/jaegerperf

WORKDIR $APP_HOME

CMD /app/jaegerperf