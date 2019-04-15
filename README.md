# docker-layer-extract

Tool to extract individual layers from a saved docker image

## Usage

### Extract your docker image to a tarball
```
$ docker save <my image tag> -o <my image file>.tar
```

### Extract newest layer from image tarball
```
$ docker-layer-extract --imagefile <my image file>.tar extract \
--newest \
--layerfile <my extract layer>.tar
```


### List layers in image tarball
```
$ docker-layer-extract --imagefile <my image file>.tar list 
...
Layer 3:  ID: e51c8d4beda7dffeeb0b0b38fdae6a22e53377207f8c089cb24e35771ebb1506
  Command: `cmd /S /C C:\vc_redist.x64.exe /quiet /install`
```

### Extract a layer from image tarball
```
$ docker-layer-extract --imagefile <my image file>.tar extract \
--layerid e51c8d4beda7dffeeb0b0b38fdae6a22e53377207f8c089cb24e35771ebb1506 \
--layerfile <my extract layer>.tar
```

## Windows Tarball format
On Windows, the extracted layer embeds Windows-specific values in a PAX header. These values are useful for docker buy not compatible with most tar implementations. To strip these headers when extracting, use the `--strip-pax` option.

## Viewing Hive Delta Entries
Each `*_Delta` file is a registry hive file and can be viewed using `regedit`
* Open Regedit
* Click on `HKEY_LOCAL_MACHINE`
* Click `File -> Load Hive`
* Navigate to your delta file example `System_Delta` and open
* Choose a memorable, non-conflicting Key Name to load it under (ex: `Temp_System_Delta`)

To close:
* Click on Key Name (ex: `Temp_System_Delta`)
* Click `File -> Unload Hive`
* Confirm 
