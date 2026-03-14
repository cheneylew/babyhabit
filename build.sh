#!/bin/bash

APP_NAME="babyhabit"
APP_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PID_FILE="$APP_DIR/$APP_NAME.pid"
LOG_FILE="$APP_DIR/$APP_NAME.log"

start() {
    if [ -f "$PID_FILE" ]; then
        PID=$(cat "$PID_FILE")
        if ps -p "$PID" > /dev/null 2>&1; then
            echo "$APP_NAME is already running (PID: $PID)"
            exit 1
        else
            rm -f "$PID_FILE"
        fi
    fi

    cd "$APP_DIR"
    nohup ./$APP_NAME > "$LOG_FILE" 2>&1 &
    echo $! > "$PID_FILE"
    echo "$APP_NAME started (PID: $(cat $PID_FILE))"
}

stop() {
    if [ -f "$PID_FILE" ]; then
        PID=$(cat "$PID_FILE")
        if ps -p "$PID" > /dev/null 2>&1; then
            kill "$PID"
            sleep 2
            if ps -p "$PID" > /dev/null 2>&1; then
                kill -9 "$PID"
            fi
            echo "$APP_NAME stopped"
        else
            echo "$APP_NAME is not running"
        fi
        rm -f "$PID_FILE"
    else
        echo "PID file not found"
    fi
}

restart() {
    stop
    sleep 1
    start
}

status() {
    if [ -f "$PID_FILE" ]; then
        PID=$(cat "$PID_FILE")
        if ps -p "$PID" > /dev/null 2>&1; then
            echo "$APP_NAME is running (PID: $PID)"
            exit 0
        else
            echo "$APP_NAME is not running (stale PID file)"
            exit 1
        fi
    else
        echo "$APP_NAME is not running"
        exit 1
    fi
}

case "$1" in
    start)
        start
        ;;
    stop)
        stop
        ;;
    restart)
        restart
        ;;
    status)
        status
        ;;
    *)
        echo "Usage: $0 {start|stop|restart|status}"
        exit 1
esac
