FROM golang:1.22.5-bullseye AS build

WORKDIR /app

COPY . ./

# install all package
RUN go mod download


# build and put on path /bin/app (app.exe)
RUN CGO_ENABLED=0 go build -o /bin/app


# next stage: install linux(defian)
FROM debian:bullseye-slim

# copy form build stage: from /bin/app to /bin
COPY --from=build /bin/app /bin

# Copy the config file to the appropriate location
COPY ./config.yaml ./config.yaml
COPY ./images ./images

# update and get all necessary pkg, install ca-certificates for HTTPS, remove unnecessary file
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# for http
EXPOSE 8080

# for run app and use .env
RUN dir -s    

CMD ["/bin/app"]