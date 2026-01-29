#!/bin/bash
# Scripts to manual test the API using curl

BASE_URL="http://localhost:8080/v1"

echo "1. Applying a Temporary Ban..."
curl -X POST "$BASE_URL/restrictions" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user_123",
    "type": "TEMP_BAN",
    "reason": "Violation of TOS",
    "duration": "48h",
    "created_by": "admin_01"
  }'
echo -e "\n"

echo "2. Checking Active Restrictions for user_123..."
curl -X GET "$BASE_URL/restrictions/user_123"
echo -e "\n"

# Note: You need the Restriction ID from step 1 to run the appeal. 
# Copy it and replace below if running manually.

echo "3. Submitting an Appeal (Example)..."
# REPLACE RESTRICTION_ID WITH REAL UUID
RESTRICTION_ID="REPLACE_WITH_UUID" 
# curl -X POST "$BASE_URL/appeals" \
#   -H "Content-Type: application/json" \
#   -d '{
#     "restriction_id": "'$RESTRICTION_ID'",
#     "user_id": "user_123",
#     "reason": "I did not check logs."
#   }'
echo " (Skipped: Restriction ID needed)"
echo -e "\n"
