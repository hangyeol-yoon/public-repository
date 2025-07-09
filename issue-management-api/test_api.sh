#!/bin/bash

# API 테스트 스크립트
BASE_URL="http://localhost:8080"

echo "=== 이슈 관리 API 테스트 ==="
echo

# 1. 헬스 체크
echo "1. 헬스 체크"
curl -X GET $BASE_URL/health
echo -e "\n"

# 2. 담당자 없는 이슈 생성
echo "2. 담당자 없는 이슈 생성 (PENDING 상태)"
curl -X POST $BASE_URL/issue \
  -H "Content-Type: application/json" \
  -d '{
    "title": "버그 수정 필요",
    "description": "로그인 페이지에서 오류 발생"
  }'
echo -e "\n"

# 3. 담당자가 있는 이슈 생성
echo "3. 담당자가 있는 이슈 생성 (IN_PROGRESS 상태)"
curl -X POST $BASE_URL/issue \
  -H "Content-Type: application/json" \
  -d '{
    "title": "새 기능 개발",
    "description": "사용자 프로필 페이지 개발",
    "userId": 1
  }'
echo -e "\n"

# 4. 전체 이슈 조회
echo "4. 전체 이슈 조회"
curl -X GET $BASE_URL/issues
echo -e "\n"

# 5. 상태별 필터링 (PENDING)
echo "5. PENDING 상태 이슈만 조회"
curl -X GET "$BASE_URL/issues?status=PENDING"
echo -e "\n"

# 6. 이슈 상세 조회
echo "6. 이슈 상세 조회 (ID: 1)"
curl -X GET $BASE_URL/issue/1
echo -e "\n"

# 7. 이슈 수정 - 담당자 할당
echo "7. 이슈 수정 - 담당자 할당 (PENDING -> IN_PROGRESS)"
curl -X PATCH $BASE_URL/issue/1 \
  -H "Content-Type: application/json" \
  -d '{
    "userId": 2
  }'
echo -e "\n"

# 8. 이슈 수정 - 상태 변경
echo "8. 이슈 수정 - 상태를 COMPLETED로 변경"
curl -X PATCH $BASE_URL/issue/1 \
  -H "Content-Type: application/json" \
  -d '{
    "status": "COMPLETED"
  }'
echo -e "\n"

# 9. 완료된 이슈 수정 시도 (에러 발생)
echo "9. 완료된 이슈 수정 시도 (에러 발생해야 함)"
curl -X PATCH $BASE_URL/issue/1 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "제목 변경 시도"
  }'
echo -e "\n"

# 10. 존재하지 않는 사용자 할당 시도 (에러 발생)
echo "10. 존재하지 않는 사용자 할당 시도 (에러 발생해야 함)"
curl -X POST $BASE_URL/issue \
  -H "Content-Type: application/json" \
  -d '{
    "title": "테스트 이슈",
    "userId": 999
  }'
echo -e "\n"

# 11. 잘못된 상태값 사용 (에러 발생)
echo "11. 잘못된 상태값 사용 (에러 발생해야 함)"
curl -X POST $BASE_URL/issue \
  -H "Content-Type: application/json" \
  -d '{
    "title": "테스트 이슈"
  }' | head -c 0

curl -X PATCH $BASE_URL/issue/2 \
  -H "Content-Type: application/json" \
  -d '{
    "status": "INVALID_STATUS"
  }'
echo -e "\n"

# 12. 필수 필드 누락 (에러 발생)
echo "12. 필수 필드(title) 누락 (에러 발생해야 함)"
curl -X POST $BASE_URL/issue \
  -H "Content-Type: application/json" \
  -d '{
    "description": "제목이 없는 이슈"
  }'
echo -e "\n"

# 13. 존재하지 않는 이슈 조회 (에러 발생)
echo "13. 존재하지 않는 이슈 조회 (에러 발생해야 함)"
curl -X GET $BASE_URL/issue/999
echo -e "\n"

# 14. 최종 이슈 목록 조회
echo "14. 최종 이슈 목록 조회"
curl -X GET $BASE_URL/issues
echo -e "\n"

echo "=== 테스트 완료 ==="