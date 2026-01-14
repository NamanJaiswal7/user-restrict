#!/bin/bash
# Helper script to commit with a date in January 2026

MSG="$1"
if [ -z "$MSG" ]; then
  echo "Usage: ./commit.sh <message>"
  exit 1
fi

# Generate a random day (1-31) and time for Jan 2026
DAY=$(( ( RANDOM % 31 ) + 1 ))
HOUR=$(( RANDOM % 24 ))
MINUTE=$(( RANDOM % 60 ))
SECOND=$(( RANDOM % 60 ))

# Format: YYYY-MM-DD HH:MM:SS
DATE_STR="2026-01-$(printf "%02d" $DAY) $(printf "%02d" $HOUR):$(printf "%02d" $MINUTE):$(printf "%02d" $SECOND)"

export GIT_AUTHOR_DATE="$DATE_STR"
export GIT_COMMITTER_DATE="$DATE_STR"

git add .
git commit -m "$MSG"
echo "Committed with date: $DATE_STR"
