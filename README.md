# Highspot Mixtape Assignment

### Assumptions
I made a few assumptions in order to proceed with the design of this command line app.
- Each type of object in mixtape.json (user, song, playlist) is uniquely identified by its respective ID
- IDs are unique so that a remove operation would have exactly one playlist to remove
- The list of song IDs belonging to a playlist do not contain duplicates (the add song to playlist operation preserves this)
- This is a command line app that would potentially be scripted against which makes logging useful for debugging

### Design
Here are some key design decisions that I made.

This is a command line app. A likely use case is it would be scripted against and used to process large batches of changes. The three operations add playlist, remove playlist, and add song to playlist would be the most frequently executed. It made sense to decrease the runtime complexity of these three operations as much as possible, with justifiable tradeoff to space complexity. With this in mind, I chose an implementation where at the start, I read all of mixtape into memory and create hashmaps to index the arrays of objects in mixtape. A hashmap each to lookup users, songs, playlists, and songs belonging to a playlist. These hashmaps provide constant runtime complexity to lookup the indices of these objects in their corresponding arrays. The exception is the last hashmap which is used as a set datastructure to check if a song exists in a playlist. With this choice:
- Add playlist has a runtime complexity of O(s) where s is the number of songs belonging to the playlist, because it's checked if these songs are in mixtape.
- Remove playlist has a runtime complexity of O(1). After looking up the index of the playlist to remove (via hashmap), it is swapped with the last playlist in the array and the length of the array is decreased by 1 (with pointers, not actually resizing the underlying array). One tradeoff is that ordering is not maintained, however, assuming these operations happen a lot, I optimized for runtime cost. It would be relatively cheaper to sort the array by playlist ID before writing the output to a file.
- Add song(s) to playlist has a runtime complexity of O(s), s is the number of songs being added. This is because its checked if the song being added is in mixtape and is not already in the existing playlist.

I decided to implement the changes file in JSON because JSON is easy to read and work with. There are other formats/protocols that are much more space efficient, which should be considered at larger scales.

I decided to add logging so it could be used for debugging purposes. Currently, all logging is sent to stdout. If there was more time, I would separate logs to stdout and stderr. For example, a log indicating a song could not be added to a playlist because it already belongs to the playlist would go to stderr.

Another design choice was to keep processing changes whenever an invalid scenario like the one above was encountered. In the real world, this decision whether to keep going or stop when a class of error is encountered would be driven by product.

There are comments throughout the code with additional design details.

### How to Build and Run
1. Install golang. There are many ways. Here is a download page with instructions: https://golang.org/doc/install
2. Create a go development directory which will be used to store go dependencies. eg.
```
mkdir -p $HOME/workspace/go
export GOPATH=$HOME/workspace/go
```
3. Clone this github repo. From the base directory, run the following to build the command line binary:
```
go build -o highspot main.go
```
4. Run the command. eg.
```
./highspot -m mixtape.json -c changes.json
```

The original mixtape.json is in the `./json` directory. There is a sample changes.json file in there too.

### How to Run Tests
1. Install `go` and set `GOPATH` env variable with steps 1 and 2 from "How to Build and Run".
2. Install the Ginkgo test framework to run tests:
```
go get github.com/onsi/ginkgo/ginkgo
go get github.com/onsi/gomega/...
```
3. Add `$GOPATH/bin` to your `PATH`. Step 3 installed a `ginkgo` binary in `$GOPATH/bin` which is needed to run tests.
4. Clone this github repo and run `ginkgo -r .` from the base directory to run all tests.

### Ideas for Running at Larger Scale
The current implementation reads mixtape into memory, performs operations asynchronously, and writes it back to a file. There are two big factors at play. If the size of mixtape is very large, we run into memory constraints. If the amount of changes we want to apply is very large, we run into CPU constraints.

With respect to large amounts of changes, one could use concurrent processing to process more in the same amount of time. Since the order of these operations matter, one would need to find a way to divide them such that either processing them becomes order agnostic or the divided unit is self contained. If such is the case, this strategy could similarly be applied to multiple machines/VMs via sharding.

With respect to a very large mixtape size, exceeding practical physical memory, one could use a database or distributed databases to store mixtape.

There are improvements that can be made to the implementation. Right now, it's expensive to create the initial hashmaps for lookup. It would be ideal to do this once for as many changes as we can apply. Perhaps instead of using file(s) to input changes, we could use a file stream.

### Known Issues
- integration tests are a bit bare, however the unit tests make up for it
- logs for invalid cases should be sent to stderr, not stdout (didn't get around to implementing this)