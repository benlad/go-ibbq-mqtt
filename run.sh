go build
sudo setcap cap_net_admin,cap_net_raw+eip go-ibbq-mqtt
LOGXI=* HA_AUTO_DISCOVERY=TRUE ./go-ibbq-mqtt