#!/bin/bash
# Scripts to manual test the API using curl

BASE_URL="http://localhost:8085/v1"

echo "1. Applying a Temporary Ban..."
RESPONSE=$(curl -s -X POST "$BASE_URL/restrictions" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user_123",
    "type": "TEMP_BAN",
    "reason": "Violation of TOS",
    "duration": "48h",
    "created_by": "admin_01"
  }')

echo "Response: $RESPONSE"
RESTRICTION_ID=$(echo $RESPONSE | jq -r '.id')
echo "Restriction ID: $RESTRICTION_ID"
echo -e "\n"

echo "2. Checking Active Restrictions for user_123..."
curl -s -X GET "$BASE_URL/restrictions/user_123" | jq
echo -e "\n"

if [ "$RESTRICTION_ID" != "null" ] && [ -n "$RESTRICTION_ID" ]; then
  echo "3. Submitting an Appeal..."
  APPEAL_RES=$(curl -s -X POST "$BASE_URL/appeals" \
    -H "Content-Type: application/json" \
    -d '{
      "restriction_id": "'$RESTRICTION_ID'",
      "user_id": "user_123",
      "reason": "I did not check logs."
    }')
  echo "Appeal Response: $APPEAL_RES"
  
  APPEAL_ID=$(echo $APPEAL_RES | jq -r '.id')
  echo "Appeal ID: $APPEAL_ID"
  echo -e "\n"

  if [ "$APPEAL_ID" != "null" ]; then
      echo "4. Reviewing Appeal (Approving)..."
      curl -s -X POST "$BASE_URL/appeals/$APPEAL_ID/review" \
        -H "Content-Type: application/json" \
        -d '{
            "reviewer_id": "admin_99",
            "status": "APPROVED",
            "notes": "Valid point, revoked."
        }' | jq
      echo -e "\n"
      
      echo "5. Verifying Restriction is Revoked..."
      curl -s -X GET "$BASE_URL/restrictions/user_123" | jq
  fi
else
  echo "Failed to capture Restriction ID, skipping appeal."
fi
