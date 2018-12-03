#!/bin/bash
case "$1" in
    parse)
        echo "parse"
        go run parse_convert.go
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