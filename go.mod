module github.com/ablanchetmd/gator

go 1.23.1

replace github.com/ablanchetmd/gator/internal/config => ../internal/config

require (
	github.com/google/uuid v1.6.0
	github.com/lib/pq v1.10.9
)
