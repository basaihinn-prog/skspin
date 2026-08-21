# build stage
FROM golang:1.14.4-alpine AS build-env

# It is important that these ARG's are defined after the FROM statement
ARG GIT_URL
ARG GIT_GROUP
ARG GIT_ACCOUNT
ARG GIT_PASSWORD
ARG BUILD_ENV

# git is required to fetch go dependencies
RUN apk update -qq && apk --no-cache add git gcc g++

# Create a netrc file using the credentials specified using --build-arg
RUN printf "machine ${GIT_URL}\n\
    login ${GIT_ACCOUNT}\n\
    password ${GIT_PASSWORD}\n"\
    >> /root/.netrc
RUN chmod 600 /root/.netrc

COPY . /src/
RUN go env -w GOPRIVATE=${GIT_URL}/${GIT_GROUP}
RUN cd /src && go build -tags ${BUILD_ENV} -o /src/dist/app && find -name '*.json' -exec cp {} /src/dist/ \;
RUN cp /src/*.crt /src/dist/
RUN cp /src/*.key /src/dist/

# final stage
FROM alpine:3.12.0
WORKDIR /app/
COPY --from=build-env /src/dist/ /app/
RUN addgroup -S playstar && adduser -S playstar -G playstar && chown -R playstar:playstar /app/
USER playstar
EXPOSE ${PORT}
ENTRYPOINT ["./app"]