## Usage

```
$ cd folder/with/photos  
$ sortPhotos <numberOfChars>
```

`numberOfChars` indicate number of first characters in filename which will be used as folder name, where file will be moved.

For example, given the following file list:

```
20180501.jpg
20181005.jpg
20191005.jpg
20191010.jpg
20191014.jpg
20191119.jpg
```

if you run:

```
sortPhotos 6
```

then files will be moved to the following folders:

```
201805
201810
201910
201911
```
