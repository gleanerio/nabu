# Nabu

## About
Loads graphs


## Notes

For now running is just running main..  there is not env for flag
settings at this time.  

### Curl load commands

```
 curl -X POST -d @test.rq --header "Content-Type:application/sparql-update" clear.local:3030/doatika/update
```

for

```
INSERT DATA {  graph <http://opencoredata.org/objectgraph/id/suffix> {<http://example.org/s> <http://example.org/p> "text here "  }}
```

### Trigger on type

The trigger code can act on all objects and look into the minio metadata to resolve actions.
Routing based on bucket, metadata and object content  ????? (what do I mean with all this?) 

## mirror commands for objects


