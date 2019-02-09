# docker-layer-extract

Tool to extract individual layers from a saved docker image

## Usage

### Extract your docker image to a tarball
```
$ docker save <my image tag> -o <my image file>.tar
```

### List layers in image tarball
```
$ docker-layer-extract list -i <my image file>.tar
...
Layer 3:  ID: e51c8d4beda7dffeeb0b0b38fdae6a22e53377207f8c089cb24e35771ebb1506
  Command: `cmd /S /C C:\vc_redist.x64.exe /quiet /install`
```

### Extract a layer from image tarball
```
$ docker-layer-extract extract -l 3 -i <my image file>.tar -o <my extract layer.tar>
```