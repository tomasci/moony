#!/bin/bash

# count lines of Go files
go_lines=$(find . -name '*.go' | xargs wc -l 2>/dev/null | grep -i ' total' | awk '{print $1}')

# count lines of GDScript files
gdscript_lines=$(find . -name '*.gd' | xargs wc -l 2>/dev/null | grep -i ' total' | awk '{print $1}')

# if there are no matches, default to 0
go_lines=${go_lines:-0}
gdscript_lines=${gdscript_lines:-0}

# calculate total lines
total_lines=$((go_lines + gdscript_lines))

# calculate percentage with rounding
if [ "$total_lines" -ne 0 ]; then
  go_percentage=$(((go_lines * 100 + total_lines / 2) / total_lines))
  gd_percentage=$(((gdscript_lines * 100 + total_lines / 2) / total_lines))
else
  go_percentage=0
  gd_percentage=0
fi

# print the results with proper percentage formatting
printf "%-10s %-5s %d%%\n" "Go" "$go_lines" "$go_percentage"
printf "%-10s %-5s %d%%\n" "GDScript" "$gdscript_lines" "$gd_percentage"
printf "%-10s %-5s\n" "Total" "$total_lines"
