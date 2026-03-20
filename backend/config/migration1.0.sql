-- Migration for vocabulary table: merge grade and textbook into book field

-- Step 1: Add book column
ALTER TABLE ab_vocabulary ADD COLUMN book VARCHAR(150) DEFAULT NULL AFTER type;

-- Step 2: Merge grade and textbook into book
UPDATE ab_vocabulary SET book = CONCAT(IFNULL(grade, ''), ' ', IFNULL(textbook, '')) WHERE book IS NULL;

-- Step 3: Remove grade and textbook columns
ALTER TABLE ab_vocabulary DROP COLUMN grade, DROP COLUMN textbook;
