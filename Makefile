DOCKERVER :=`cat VERSION`
.DEFAULT_GOAL := loadv2 
 
loadGraphs:
	cd cmd/loadGraphs; \
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 env go build -o loadGraphs
 
  
loadv2:
	cd cmd/loadv2; \
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 env go build -o loadv2

releases:    loadGraphs loadv2

