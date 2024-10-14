# Exercise 2
<code>Shred(path)</code> function that overwrites any given file 3 times with random data, and deletes the file afterwards.

# Summary
The module <code>shred.go</code> contains the function <code>Shred(path)</code> which overwrites the given file with random data obtained using the module <code>crypto/rand</code>, flushing the file, closing it, and reopening it to perform these actions 2 more times.

I've decided to use the module <code>bufio</code> for buffered writing, as well as using a randomization block size of 4KB. By using these two techniques I've managed to improve speed performance and also reduce the memory footprint.

# Test Cases
I've implemented test cases for <code>Shred(path)</code> which includes:
- testing the <code>path</code> parameter for other than a regular file
- testing the byte contents to be randomized before deletion
- testing for correct file deletion after randomization

I've included <code>shred_test.go</code> module for test coverage, which gives the following results:
```
pablo@machine:~/shred$ go test -v -cover
=== RUN   TestShredDir
--- PASS: TestShredDir (0.00s)
=== RUN   TestShredNothing
--- PASS: TestShredNothing (0.00s)
=== RUN   TestShredRegularFile
--- PASS: TestShredRegularFile (0.00s)
=== RUN   TestOverwriteRegularFile
--- PASS: TestOverwriteRegularFile (0.00s)
PASS
coverage: 93.8% of statements
ok  	github.com/pabloandresm/shred	0.005s
```

# Use cases, pros and cons
The main use case of <code>Shred(path)</code> function is security, so if you delete a file with this function, it will be useless if undeleted.

- The main drawback using this function is _**speed**_:
 - Since the file is overwritten 3 times, then its size will impact duration.
 - Since the file could be located on different media types (HDD, SSD, network shares, flash drives), then that technology will impact duration.

 - On data-journaling filesystems the data might still be present on the underlying medium.

 - Another drawback could be _**free storage space**_. Some filesystems support the _**sparse**_ feature, which means that the blocks with no data (zeros) are not physically allocated. By using <code>Shred(path)</code> the entire size of the file will be allocated before being deleted, which can lead to no free space before completing the task.
 - If the function is used on a file located on a solid state device, then _**it will not succeed in hidding its content**_. This is because of the wear-level mechanisms that these type of devices use, which will translate and hide the mapping between a logical block to a physical one. So after shredding a file, if you read raw blocks from the device you could read content that was supposed to be shredded.

# Usage
```
package main

import (
    "github.com/pabloandresm/shred"
    "fmt"
    )

func main() {
//
    err := shred.Shred("tmp_file")
    if err != nil {
        fmt.Println(err)
        }
//
}
```

# Installation
```
go get -u github.com/pabloandresm/shred
```

---
Pablo Martikian

23 May 2022
