go-oci8
=======

[![Build Status](https://travis-ci.org/mattn/go-oci8.svg)](https://travis-ci.org/mattn/go-oci8)

Description
-----------

Oracle driver conforming to the built-in database/sql interface

Installation
------------

This package can be installed with the go get command:

    go get github.com/mattn/go-oci8

You need to put `oci8.pc` like into your `$PKG_CONFIG_PATH`. `oci8.pc` should be like below.

### Example for Windows

```
prefix=/devel/target/XXXXXXXXXXXXXXXXXXXXXXXXXX
exec_prefix=${prefix}
libdir=c:/oraclexe/app/oracle/product/11.2.0/server/oci/lib/msvc
includedir=c:/oraclexe/app/oracle/product/11.2.0/server/oci/include/include

glib_genmarshal=glib-genmarshal
gobject_query=gobject-query
glib_mkenums=glib-mkenums

Name: oci8
Description: oci8 library
Libs: -L${libdir} -loci
Cflags: -I${includedir}
Version: 11.2
```

### Example for Linux

```
prefix=/devel/target/XXXXXXXXXXXXXXXXXXXXXXXXXX
exec_prefix=${prefix}
libdir=/usr/lib/oracle/11.2/client64/lib
includedir=/usr/include/oracle/11.2/client64

glib_genmarshal=glib-genmarshal
gobject_query=gobject-query
glib_mkenums=glib-mkenums

Name: oci8
Description: oci8 library
Libs: -L${libdir} -lclntsh
Cflags: -I${includedir}
Version: 11.2
```

Documentation
-------------

API documentation can be found here: http://godoc.org/github.com/mattn/go-oci8

Examples can be found under the `./_example` directory

License
-------

MIT: http://mattn.mit-license.org/2014

ToDo
----

* LastInserted is not int64
* Fetch number is more improvable

Author
------

Yasuhiro Matsumoto (a.k.a mattn)

Special Thanks
--------------

Jamil Djadala
