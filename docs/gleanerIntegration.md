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
minio:
  address: localhost
  port: 9000
  ssl: false
  accesskey: akey
  secretkey: skey
  bucket: gleaner2
objects:
  bucket: gleaner2
  domain: us-east-1
  prefix:
  - milled/iris
  - prov/iris
  - org
  prefixoff:
  - summoned/aquadocs
  - prov/aquadocs
  - milled/opentopography
  - prov/opentopography
sparql:
  endpoint: http://localhost/blazegraph/namespace/earthcube/sparql
  authenticate: false
  username: ""
  password: ""
txtaipkg:
  endpoint: http://0.0.0.0:8000
```

Commands are like the following:
### Help
nabu --help

### Prefix: load all objects in the bucket/prefixes

The mode "prefix" in Nabu is used for loading a S3 object prefix 
path into a triplestore
```
nabu prefix --cfg file
```
eg
```
nabu prefix --cfg ../gleaner/configs/nabu
```
and gleaner generated configuration: 
```
nabu prefix --cfgPath directory --cfgName name 
```
eg use generated
```
nabu prefix --cfgPath ../gleaner/configs --cfgName local 
```
### Prune: load all objects in the bucket/prefixes

The mode "prune" in Nabu is to sync a prefix to the graph (remove graphs no longer in use, add new ones)

Note that updated graphs become new objects, since the object name is the SHA256 of the object

```
nabu prune --cfg file 
```
and gleaner generated configuration:
```
nabu prefix --cfgPath directory --cfgName name 
```
eg use generated
```
nabu prefix --cfgPath ../gleaner/configs --cfgName local 
```

### Object: load one object in the bucket/prefixes

The mode "object" in Nabu is used for loading a S3 object 
path into a triplestore
```
nabu object --cfg file objectId
```
eg
```
nabu object --cfg ../gleaner/configs/nabu milled/opentopography/ffa0df033bb3a8fc9f600c80df3501fe1a2dbe93.rdf
```

and gleaner generated configuration:
```
nabu object --cfgPath directory --cfgName name  objectId
```
eg use generated
```
nabu object --cfgPath ../gleaner/configs --cfgName local milled/opentopography/ffa0df033bb3a8fc9f600c80df3501fe1a2dbe93.rdf
```

## TODO

- overview here of the basic config, source, app sink flow
- The quad as prov pattern where the path to the triples, becomes
the quad value 
 
