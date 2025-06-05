FROM golang:1.24.3 AS build

# Dependency install
RUN apt-get update && apt-get install -y vim && apt-get install systemctl -y
RUN apt install dbus-x11 -y