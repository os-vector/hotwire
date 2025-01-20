go build cmd/main.go
sudo setcap 'cap_net_bind_service=+ep' main
DEBUG_LOGGING=true ./main
