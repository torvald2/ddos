FROM golang as gobuild 
WORKDIR /app

COPY go.mod ./
COPY *.go ./


RUN go build -o /app

ENTRYPOINT /app

FROM alpine
USER root
WORKDIR /
COPY *.txt ./

COPY --from=gobuild /app /myapp
RUN ["chmod", "+x", "/myapp"]

ENTRYPOINT ["/myapp"]