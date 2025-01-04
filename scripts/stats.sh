#!/bin/bash

# Function to calculate the percentage and round it to one decimal
calculate_percentage() {
  local lines=$1
  local total=$2
  if [ "$total" -gt 0 ]; then
    echo "scale=2; ($lines * 100) / $total" | bc
  else
    echo "0.00"
  fi
}

# Count lines of Go files
go_lines=$(find . -name '*.go' | xargs wc -l 2>/dev/null | grep -i ' total' | awk '{print $1}')

# Count lines of GDScript files
gdscript_lines=$(find . -name '*.gd' | xargs wc -l 2>/dev/null | grep -i ' total' | awk '{print $1}')

# Count shell lines
shell_lines=$(find . -name '*.sh' | xargs wc -l 2>/dev/null | grep -i ' total' | awk '{print $1}')

# If there are no matches, default to 0
go_lines=${go_lines:-0}
gdscript_lines=${gdscript_lines:-0}
shell_lines=${shell_lines:-0}

# Calculate total lines
total_lines=$((go_lines + gdscript_lines + shell_lines))

# Calculate and round percentages
go_percentage=$(calculate_percentage $go_lines $total_lines)
gd_percentage=$(calculate_percentage $gdscript_lines $total_lines)
shell_percentage=$(calculate_percentage $shell_lines $total_lines)

# Print the results with proper percentage formatting
printf "%-10s %-10s %s%%\n" "Go" "$go_lines" "$go_percentage"
printf "%-10s %-10s %s%%\n" "GDScript" "$gdscript_lines" "$gd_percentage"
printf "%-10s %-10s %s%%\n" "Shell" "$shell_lines" "$shell_percentage"
printf "%-10s %-10s\n" "Total" "$total_lines"
