@echo off

setlocal

set "command=%1"
set "migration_name=%~2"

if "%command%" == "migrate-create" (
    echo Creating migrations...
    if "%migration_name%"=="" (
        echo Please provide a migration name.
        exit /b 1
    )
    migrate create -ext sql -dir cmd/migrate/migrations "%migration_name%"
    echo Migrations created successfully!
    exit /b
)

if "%command%" == "migrate-up" (
    echo Applying migrations...
    go run cmd/migrate/main.go up
    echo Migrations applied successfully!
    exit /b
)

if "%command%" == "migrate-down" (
    echo Reverting migrations...
    go run cmd/migrate/main.go down
    echo Migrations reverted successfully!
    exit /b
)

echo Invalid option. Use:
echo   migrate-create [name] - to create a migration
echo   migrate-up - to apply migrations
echo   migrate-down - to revert migrations
exit /b 1
