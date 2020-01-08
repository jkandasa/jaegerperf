FROM alpine:3.7

ENV APP_HOME /app/

# copy go build file
COPY ./main $APP_HOME/jaegerperf

EXPOSE 8080

WORKDIR $APP_HOME
CMD /app/jaegerperf