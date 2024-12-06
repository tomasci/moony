#!/bin/bash

clear

# update chmod for scripts
find ./scripts/ -type f -exec chmod +x {} \;

alias help="moony:help"
alias moony:dev="bash ./scripts/dev.sh"

moony:help() {
  echo "\nHelp"
  echo ""
  printf "%-20s %s\n" "help" "Show help"
  printf "%-20s %s\n" "moony:help" ""
  echo ""
  printf "%-20s %s\n" "moony:dev" "Run server in dev mode"
  printf "%-20s %s\n" "" "(basically will try to build and run)"
  echo ""
}

echo "moony environment loaded"