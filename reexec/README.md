## reexec

This is fork from [docker/reexec](https://github.com/moby/moby/tree/v1.13.1/pkg/reexec).

This package modify `reexec.Init` to use `path.Base(os.Args[0])` instead of 
directly `os.Args[0]`.
