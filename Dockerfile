# Актуализировать версию образа golang (при необходимости)
ARG BUILD_IMAGE_NAME=golang:1.25-alpine
ARG DEPLOY_IMAGE_NAME=alpine

# Образ с необходимыми пакетами
FROM  ${BUILD_IMAGE_NAME} AS packages_image
ENV PATH=/usr/local/go/bin:$PATH
WORKDIR /go/src
COPY go.mod go.sum ./
RUN echo "packages_image" && \
    go mod download

# Образ для сборки приложения   
FROM ${BUILD_IMAGE_NAME} AS build_image
# Метка. По имени метки можно найти и удалить образ.
LABEL temporary="" 
ARG APPLICATION_NAME=application_template
ARG BRANCH=master
ARG USERNAME
ARG USERPASSWD
ARG VERSION=0.1.1
WORKDIR /go/
COPY --from=packages_image /go ./
RUN echo -e "build_image" && \
    rm -r ./src && \
    apk update && \
    apk add --no-cache git && \
    # брать исходный код с репозитория на gitlab.cloud.gcm 
    #git clone -b ${BRANCH} http://${USERNAME}:${USERPASSWD}@gitlab.cloud.gcm/a.belyakov/${APPLICATION_NAME}.git ./src/${VERSION}/ && \
    #git clone -b ${BRANCH} http://${USERNAME}:${USERPASSWD}@192.168.9.33/a.belyakov/${APPLICATION_NAME}.git ./src/${VERSION}/ && \
    # брать исходный код с репозитория на github.com 
    git clone -b ${BRANCH} https://github.com/av-belyakov/${APPLICATION_NAME}.git ./src/${VERSION}/ && \
    go build -C ./src/${VERSION}/cmd/ -o ../app

# Основной рабочий образ     
FROM ${DEPLOY_IMAGE_NAME}
LABEL author="Artemij Belyakov"
ARG APPLICATION_NAME=application_template
# аргумент STATUS содержит режим запуска приложения prod или development
# если значение содержит запись development, то в таком режиме и будет
# работать приложение, во всех остальных случаях режим работы prod
ARG STATUS=prod
ARG USER_NAME=dockeruser
ARG LOGS_DIR=logs
ARG VERSION=0.1.1
#!!! здесь заменить переменную окружения на соответствующую имени приложения !!!
ENV GO_APPLICATION_TEMPLATE_MAIN=${STATUS}
WORKDIR /opt/${APPLICATION_NAME}
RUN addgroup --g 1500 groupcontainer && \
    adduser -u 1500 -G groupcontainer -D ${USERNAME} --home ${WORK_DIR}
USER ${USER_NAME}
RUN mkdir ./${LOGS_DIR}
COPY --from=build_image /go/src/${VERSION}/app ./
COPY --from=build_image /go/src/${VERSION}/README.md ./
COPY --from=build_image /go/src/${VERSION}/version ./ 
COPY config/* ./config/

ENTRYPOINT [ "./app" ]
