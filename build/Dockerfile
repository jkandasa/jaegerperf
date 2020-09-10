FROM alpine:3.12

LABEL maintainer="Jeeva Kandasamy <jkandasa@redhat.com>"

ENV APP_HOME /app/

EXPOSE 8080

# create a user and give permission for the locations
RUN mkdir -p ${APP_HOME}

# copy WEB files
COPY ./web ${APP_HOME}/web

# copy application bin file
COPY ./jaegerperf ${APP_HOME}/jaegerperf

# copy default files
COPY ./resources ${APP_HOME}/resources

WORKDIR ${APP_HOME}

# update permissions to all user
# temporery fix to work on OpenShift
RUN chmod 0777 ${APP_HOME} -R

CMD /app/jaegerperf
