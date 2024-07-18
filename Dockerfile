FROM golang:1.22.5 as builder
WORKDIR /app
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -v -o q .

FROM gcr.io/google.com/cloudsdktool/google-cloud-cli:slim
COPY --from=builder /app/q /usr/local/bin/q
RUN apt-get install jq vim -y
