nats -s localhost:4225 kv get nalpaca positions --raw | \
 protoc --proto_path ./api/proto -I /usr/local/include --decode tradesvc.v0.Positions api/proto/tradesvc/v0/tradesvc.proto