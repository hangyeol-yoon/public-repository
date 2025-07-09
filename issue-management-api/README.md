# 이슈 관리 API

백엔드 개발자 채용 과제를 위한 이슈 관리 REST API입니다.

## 실행 방법

### 1. 사전 요구사항
- Go 1.21 이상 설치
- 포트 8080번이 사용 가능해야 함

### 2. 실행 단계
```bash
# 1. 프로젝트 클론
git clone <repository-url>
cd issue-management-api

# 2. 의존성 설치
go mod tidy

# 3. 서버 실행
go run main.go
```

서버가 성공적으로 실행되면 다음과 같은 메시지가 출력됩니다:
```
Server starting on port 8080...
```

## API 테스트 방법

### 헬스 체크
```bash
curl -X GET http://localhost:8080/health
```

### 1. 이슈 생성 (POST /issue)

#### 담당자 없는 이슈 생성
```bash
curl -X POST http://localhost:8080/issue \
  -H "Content-Type: application/json" \
  -d '{
    "title": "버그 수정 필요",
    "description": "로그인 페이지에서 오류 발생"
  }'
```

#### 담당자가 있는 이슈 생성
```bash
curl -X POST http://localhost:8080/issue \
  -H "Content-Type: application/json" \
  -d '{
    "title": "새 기능 개발",
    "description": "사용자 프로필 페이지 개발",
    "userId": 1
  }'
```

### 2. 이슈 목록 조회 (GET /issues)

#### 전체 이슈 조회
```bash
curl -X GET http://localhost:8080/issues
```

#### 상태별 필터링
```bash
# PENDING 상태 이슈만 조회
curl -X GET "http://localhost:8080/issues?status=PENDING"

# IN_PROGRESS 상태 이슈만 조회
curl -X GET "http://localhost:8080/issues?status=IN_PROGRESS"
```

### 3. 이슈 상세 조회 (GET /issue/:id)
```bash
curl -X GET http://localhost:8080/issue/1
```

### 4. 이슈 수정 (PATCH /issue/:id)

#### 제목과 설명 수정
```bash
curl -X PATCH http://localhost:8080/issue/1 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "로그인 버그 수정",
    "description": "로그인 폼 검증 로직 개선"
  }'
```

#### 담당자 할당
```bash
curl -X PATCH http://localhost:8080/issue/1 \
  -H "Content-Type: application/json" \
  -d '{
    "userId": 2
  }'
```

#### 상태 변경
```bash
curl -X PATCH http://localhost:8080/issue/1 \
  -H "Content-Type: application/json" \
  -d '{
    "status": "COMPLETED"
  }'
```

#### 복합 수정
```bash
curl -X PATCH http://localhost:8080/issue/1 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "로그인 버그 수정 완료",
    "status": "COMPLETED",
    "userId": 2
  }'
```

## API 엔드포인트

| 메서드 | 경로 | 설명 |
|--------|------|------|
| POST | /issue | 이슈 생성 |
| GET | /issues | 이슈 목록 조회 |
| GET | /issue/:id | 이슈 상세 조회 |
| PATCH | /issue/:id | 이슈 수정 |
| GET | /health | 헬스 체크 |

## 데이터 모델

### 사용자 (User)
```json
{
  "id": 1,
  "name": "김개발"
}
```

### 이슈 (Issue)
```json
{
  "id": 1,
  "title": "버그 수정 필요",
  "description": "로그인 페이지에서 오류 발생",
  "status": "PENDING",
  "user": {
    "id": 1,
    "name": "김개발"
  },
  "createdAt": "2025-06-02T10:00:00Z",
  "updatedAt": "2025-06-02T10:00:00Z"
}
```

## 비즈니스 규칙

### 이슈 상태
- `PENDING`: 대기 중
- `IN_PROGRESS`: 진행 중
- `COMPLETED`: 완료
- `CANCELLED`: 취소됨

### 상태 변경 규칙
1. 담당자가 없는 경우 `PENDING` 또는 `CANCELLED` 상태만 가능
2. 담당자가 할당되면 기본적으로 `IN_PROGRESS` 상태로 설정
3. `COMPLETED` 또는 `CANCELLED` 상태의 이슈는 수정 불가
4. 담당자 제거 시 상태는 자동으로 `PENDING`으로 변경

### 기본 사용자
시스템에 기본으로 등록된 사용자:
- ID: 1, 이름: "김개발"
- ID: 2, 이름: "이디자인"  
- ID: 3, 이름: "박기획"

## 에러 응답 형식

모든 에러는 다음 형식으로 응답됩니다:
```json
{
  "error": "에러 메시지",
  "code": 400
}
```

### 주요 에러 상황
- 400 Bad Request: 잘못된 요청 데이터
- 404 Not Found: 존재하지 않는 리소스
- 422 Unprocessable Entity: 비즈니스 규칙 위반

## 테스트 시나리오

### 시나리오 1: 기본 이슈 생성 및 조회
1. 담당자 없는 이슈 생성 → `PENDING` 상태
2. 이슈 목록 조회로 생성 확인
3. 이슈 상세 조회

### 시나리오 2: 담당자 할당 및 상태 변경
1. 담당자 할당 → 자동으로 `IN_PROGRESS` 상태
2. 상태를 `COMPLETED`로 변경
3. 완료된 이슈 수정 시도 → 에러 발생

### 시나리오 3: 에러 처리
1. 존재하지 않는 사용자 할당 시도
2. 잘못된 상태값 사용
3. 필수 필드 누락

## 주요 구현 특징

- **REST API 설계**: 표준 HTTP 메서드와 상태 코드 사용
- **에러 처리**: 일관된 에러 응답 형식
- **데이터 검증**: 입력 데이터 유효성 검사
- **비즈니스 로직**: 요구사항에 맞는 상태 변경 규칙 구현
- **코드 구조**: 핸들러, 모델, 유틸리티 함수 분리