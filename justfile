alias on := start
alias off := stop

# Lists all available tasks
default:
    @just --list

# Start infrastructure
start:
    docker compose up -d

# Stop infrastructure
stop:
    docker compose down
