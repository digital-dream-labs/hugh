# hugh 

Hugh is an _extremely_ minimal set of utilities to minimize boilerplate code in [12 factor apps](https://12factor.net/).  

It is _not_ intended to be a framework, because frameworks are (in my opinion) far too opinionated and tend to spider throughout the codebase.

## Package structure

|Directory| Description |
|--|--|
| config | A wrapper for viper config loading |
| database | reduce database connection boilerplate |
| grpc | client/server/interceptor libraries |
| log | because no package such as this is complete without yet another log implementation |
| testing | dockerized instances of some handy things for to be used in testing packages |

