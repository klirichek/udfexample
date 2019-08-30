## Introduction

UDFexample is a set of demo functions and demo plugins which can be plugged to Manticoresearch daemon
and provide different enhanced actions with your data, like custom functions, custom tokenization
during indexing and querying, and custom ranking.

Manticore Search is an open source search server designed to be fast, scalable and with powerful and accurate full-text search capabilities. It is a fork of popular search engine Sphinx.

## Examples
* Stateless UDF example 'strtoint' which transforms string to number using Go fmt parameter
* Stateful UDF example 'avgmva' which calculates average of provided MVA attribute
* Stateless UDF example 'inspect' which just parses parameters and report them back to daemon
* Stateless UDF example 'curl' which downloads given resource and returns it, if it is text
* Stateful UDF example 'sequence' which returns monotonically growing integers
* tokenizer plugin example 'hideemail' converts token 'any@space.io' to 'mailto:any@space.io' and drop any other emails
* query tokenizer plugin example 'queryshow' just displays back all calls and parameters
* ranker plugin example 'myrank' also just displays back all calls and parameters

## Installation
First, clone or download the repo.
Also you need `sphinxudf.h` from manticore sources (look at /src there). Add it to cloned sources
Then run `go build -buildmode=c-shared -o udfexample.so .`

### Plugging

Place built `udfexample.so` somewhere in the system. Add param `plugin_dir = /path/to/udf/dir` to your
config into 'common' section.

### Usage

Usage is different from the type of example you want to work with.
UDFs usually need to be first loaded via mysql console, for example

```mysql
  CREATE FUNCTION avgmva RETURNS FLOAT SONAME 'udfexample.so';
```

Tokenfilter is added to 'index' section with the line like:

```
index_token_filter = udfexample.so:hideemail:opt=blabla;another=bar
```

Query token filter have to be attached to query in runtime, like:

```mysql
select * from ru where match ('сталь') OPTION token_filter='udfexample.so:queryshow:bla';l
```

Ranker plugin has first to be plugged, and then used via query in runtime, like:

```mysql
CREATE PLUGIN myrank TYPE 'ranker' SONAME 'udfexample.so';
SELECT * from ru WHERE match ('сталь') OPTION ranker=myrank('option1=1');
```

## Structure
Examples are written on Go with CGo package. Base file `udfhelpers.go` contains C binding skeleton of the
library and some helper functions for converting/transferring data between go and C.

Each example, in turn, located in separate file and may be used standalone together with helpers.
