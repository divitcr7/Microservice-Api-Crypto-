FROM umputun/baseimage:buildgo-master as build 
ARG VERSION
WORKDIR /build
ADD . /build
COPY go.* ./
RUN go mod download

RUN apk update && apk upgrade
RUN go build -mod=readonly -o app -ldflags="-X github.com/streamdp/ccd/config.Version=$VERSION" ./

FROM umputun/baseimage:app

LABEL org.opencontainers.image.authors="Streamdp <@streamdp>"

ENV TIME_ZONE="Europe/Minsk"

COPY --from=build /build/app 	    /srv/app
COPY --from=build /build/site 	/srv/site/

WORKDIR /srv
EXPOSE 8080

CMD ["/srv/app"]
