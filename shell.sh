#!/bin/bash

clear

# update chmod for scripts
find ./scripts/ -type f -exec chmod +x {} \;

alias help="moony:help"
alias moony:dev="bash ./scripts/dev.sh"
alias moony:build="bash ./scripts/build.sh"
alias moony:start="./build/main"

moony:help() {
  echo "\nHelp"
  echo ""
  printf "%-20s %s\n" "help" "Show help"
  printf "%-20s %s\n" "moony:help" ""
  echo ""
  printf "%-20s %s\n" "moony:dev" "Run server in dev mode"
  printf "%-20s %s\n" "" "(basically will try to build and run)"
  echo ""
  printf "%-20s %s\n" "moony:build" "Build server"
  echo ""
  printf "%-20s %s\n" "moony:start" "Run server at 127.0.0.1"
  printf "%-20s %s\n" "-host" "with host flag it will run at 0.0.0.0"
  printf "%-20s %s\n" "" "run moony:build or moony:dev before start command"
  echo ""
}

echo "moony environment loaded"