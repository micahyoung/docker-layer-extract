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
On Windows, the tarfile doesn't not appear to be a fully standard TAR format. You can still extract most of the relevant data but symlinks and other metadata are will cause issues for some TAR clients. I found the best tool for extracting is `tar` from the Windows, but `7-zip` also seems to extract well enough if you say "No" to all override prompts.


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
