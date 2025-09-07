-- down migration

-- 1. 인덱스 삭제
DROP INDEX IF EXISTS accounts_owner_currency_idx;

-- 2. 외래키 제약 조건 삭제
ALTER TABLE "accounts" DROP CONSTRAINT IF EXISTS accounts_owner_fkey;

-- 3. users 테이블 삭제
DROP TABLE IF EXISTS "users";