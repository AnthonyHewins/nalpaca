FROM docker.io/library/golang:1.24-bullseye AS build

ARG target

WORKDIR /app
ADD ./ /app
RUN make $target && mv ./bin/$target ./bin/server

# Use this tag for debug purpose only
#FROM gcr.io/distroless/base-debian11:debug
FROM gcr.io/distroless/base-debian12
WORKDIR /root/
COPY --from=build /app/bin/server /root

EXPOSE 8088 7654
ENTRYPOINT ["/root/server"]
