DOCKERVER :=`cat VERSION`
.DEFAULT_GOAL := nabu
 
  
nabu:
	cd cmd/nabu; \
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 env go build -o nabu

releases:    nabu

