#  Loading Gleaner Output into Triplestores

## About

Nabu primarily is a tool for reading from an S3 object store and writing 
to a triplestore.  The object stores can be any S3 compliant object stores
so Google Cloud Storage, Wasabi, or others.   For most cases I am using
Minio, an open source S3 object store. 

Similarly, the triplestore can be any standards compliant triplestore.  Here
the primary standards we need implemented include

* SPARQL 1.1 with Update support
* SPARQL 1.1 over HTTP

## Basic operations

Config file example

```
{
  "sparql": {
    "endpoint": "https://graph.example.org/sparql"
  },
  "objects": {
      "bucket": "sourcebucket",
      "prefix": ["bucket, "prefixes", "to", "load"],
      "domain": "us-east-1"
  },
  "minio": {
    "address": "objectstore.example.org",
    "port": "",
    "accesskey": "worldsbestaccesskey",
    "secretkey": "worldsbestsecretkey"
  }
}
```

Commands are like the following:

### Prefix: load all objects in the bucket/prefixes

The mode "prefix" in Nabu is used for loading a S3 object prefix 
path into a triplestore
```
nabu -cfg file -mode prefix
```

### Prune: load all objects in the bucket/prefixes

The mode "prune" in Nabu is to sync a prefix to the graph (remove graphs no longer in use, add new ones)

Note that updated graphs become new objects, since the object name is the SHA256 of the object

```
nabu -cfg file -mode prune
```

## TODO

- overview here of the basic config, source, app sink flow
- The quad as prov pattern where the path to the triples, becomes
the quad value 
 
