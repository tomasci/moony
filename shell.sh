#!/bin/bash

clear

# update chmod for scripts
find ./scripts/ -type f -exec chmod +x {} \;

# run commands
alias moony:dev="bash ./scripts/dev.sh"
alias moony:build="bash ./scripts/build.sh"
alias moony:start="./build/main"

# migrations
alias moony:migrate:up="bash ./scripts/migrate_up.sh"
alias moony:migrate:down="bash ./scripts/migrate_down.sh"
alias moony:migrate:create="bash ./scripts/migrate_create.sh"
alias moony:migrate:generate="bash ./scripts/migrate_generate.sh"

# redis
alias moony:redis:setup="bash ./scripts/run_redis_setup.sh"
alias moony:redis:cli="bash ./scripts/run_redis_cli.sh"
alias moony:redis:dev="bash ./scripts/run_redis_dev.sh"

# display amount of written go-lines in project
alias moony:stats="bash ./scripts/stats.sh"

echo "moony environment loaded"