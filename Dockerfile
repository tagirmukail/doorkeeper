FROM alpine
ENV LANGUAGE="en"

ADD doorkeeper /doorkeeper
ADD config.json /config.json
#CMD CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

RUN apk add --no-cache ca-certificates

CMD ["chmod", "+x", "/doorkeeper"]
CMD [ "/doorkeeper" ]

