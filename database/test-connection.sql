-- Test PostgreSQL Connection
-- Run this in VS Code to verify your connection works

-- Check PostgreSQL version
SELECT version();

-- Check if database exists
SELECT current_database();

-- List all tables (should be empty if database just created)
SELECT
    table_schema,
    table_name,
    table_type
FROM information_schema.tables
WHERE table_schema NOT IN ('pg_catalog', 'information_schema')
ORDER BY table_schema, table_name;

-- After running migrations, check table count
SELECT
    COUNT(*) as table_count
FROM information_schema.tables
WHERE table_schema = 'public';

-- View all tables with row counts (after loading data)
SELECT
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY tablename;
