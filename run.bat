go build -o bookings.exe ./cmd/web/. || exit /b
bookings.exe -production=false -cache=false -dbname=bookings -dbuser=postgres -dbpass=India@100