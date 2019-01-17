#!/bin/bash
case "$1" in
    parse)

        ;;
    run)
        echo "run"
        go run main.go
        ;;
    build)
        echo "build"
        ;;
     *)
        echo $"Usage: $0 {run|parse|build|}"
        exit 1

esac