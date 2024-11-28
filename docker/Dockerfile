FROM docker.io/library/golang:1.22-bullseye AS build

WORKDIR /app
ADD ./ /app
RUN make nalpaca

# Use this tag for debug purpose only
#FROM gcr.io/distroless/base-debian11:debug
FROM gcr.io/distroless/base-debian12
WORKDIR /root/
COPY --from=build /app/bin/* /root

EXPOSE 8088 7654
ENTRYPOINT ["/root/orders"]
