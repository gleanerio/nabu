# 1. Record architecture decisions

Date: 05-01-2023

## Status

Proposed

## Context

URNs for the data graph URI are set in the file internal/graph/mintURN.go
This is a decision regarding the graph URN for the data graphs, not the 
prov graphs.  

current
```
urn:{bucket}:{docstatus}:{provider}:{sha}
```

proposed
```
urn:{bucket}:{provider}:{sha}
```

## Decision

So under the old approach we had URNs like

```rdf
urn:gleaner.io:summoned:edmo:0255293683036aac2a95a2479cc841189c0ac3f8
```
or
```rdf
urn:gleaner.io:milled:edmo:0255293683036aac2a95a2479cc841189c0ac3f8
```

The milled and summoned elements were pointless and led to confusion and were not 
really important in terms of getting to the object.  

The new desired URN pattern would then look like.  These would likely always be pulled
from the summoned prefix, and as such be JSON-LD.  

```rdf
urn:gleaner.io:edmo:0255293683036aac2a95a2479cc841189c0ac3f8
```

## Consequences

This impacts gleaner in the generation of prov which will need to use this same pattern
to fill out the prov records.  

Also, this means the URN does not actually represent the location of the object.  Rather the 
client must know to go looking in summoned and or milled.  As noted, the use of milled is 
not really compelling.  That aside, it is a case where the URN is now just an identifier and 
does not represent a resolvable object.  

Given this point, it is up for discussion if the URN might be better as:

alternative
```
urn:{bucket}:summoned:{provider}:{sha}
```

Note also here, the sha is now based on the more expressive approach coded into Gleaner and 
not just the sha of the data graph by default. 
