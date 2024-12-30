#!/bin/bash

function show_menu() {
    echo "Select Environment:"
    echo "1) Development"
    echo "2) Production"
    read -p "Enter choice [1-2]: " env_choice

    case $env_choice in
        1)
            COMPOSE_FILE="docker-compose.dev.yml"
            ;;
        2)
            COMPOSE_FILE="docker-compose.prod.yml"
            ;;
        *)
            echo "Invalid choice!"
            exit 1
            ;;
    esac

    echo "Select Action:"
    echo "1) Up (Start and Build)"
    echo "2) Down (Stop and Remove)"
    echo "3) Start"
    echo "4) Stop"
    echo "5) Build"
    read -p "Enter choice [1-5]: " action_choice

    case $action_choice in
        1)
            ACTION="up --build"
            ;;
        2)
            ACTION="down"
            ;;
        3)
            ACTION="start"
            ;;
        4)
            ACTION="stop"
            ;;
        5)
            ACTION="build"
            ;;
        *)
            echo "Invalid choice!"
            exit 1
            ;;
    esac
}

show_menu
docker compose -f $COMPOSE_FILE $ACTION