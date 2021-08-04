FROM golang as builder

WORKDIR /app
COPY . .

RUN go build -o /app/bot smuel1414/gcloud.vms

CMD ["/app/bot"]