rd /s /q "../protobeans"
mkdir "../protobeans"
protoc --go_out=../protobeans *.proto