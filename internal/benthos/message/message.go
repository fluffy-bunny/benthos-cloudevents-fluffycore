package message

import (
	"github.com/gogo/status"
	"google.golang.org/grpc/codes"
)

//------------------------------------------------------------------------------

/*
Internal message blob format:

- Four bytes containing number of message parts in big endian
- For each message part:
    + Four bytes containing length of message part in big endian
    + Content of message part

                                         # Of bytes in message part 2
                                         |
# Of message parts (u32 big endian)      |           Content of message part 2
|                                        |           |
v                                        v           v
| 0| 0| 0| 2| 0| 0| 0| 5| h| e| l| l| o| 0| 0| 0| 5| w| o| r| l| d|
  0  1  2  3  4  5  6  7  8  9 10 11 13 14 15 16 17 18 19 20 21 22
              ^           ^
              |           |
              |           Content of message part 1
              |
              # Of bytes in message part 1 (u32 big endian)
*/

// DeserializeBytes rebuilds a 2D byte array from a binary serialized blob.
func DeserializeBytes(b []byte) ([][]byte, error) {
	if len(b) < 4 {
		return nil, status.Error(codes.InvalidArgument, "message is not long enough to contain number of parts")
	}

	numParts := uint32(b[0])<<24 | uint32(b[1])<<16 | uint32(b[2])<<8 | uint32(b[3])
	if numParts >= uint32(len(b)) {
		return nil, status.Error(codes.InvalidArgument, "message is not long enough to contain all parts")
	}

	b = b[4:]

	parts := make([][]byte, numParts)

	for i := uint32(0); i < numParts; i++ {
		if len(b) < 4 {
			return nil, status.Error(codes.InvalidArgument, "message is not long enough to contain part length")
		}
		partSize := uint32(b[0])<<24 | uint32(b[1])<<16 | uint32(b[2])<<8 | uint32(b[3])
		b = b[4:]

		if uint32(len(b)) < partSize {
			return nil, status.Error(codes.InvalidArgument, "message is not long enough to contain all parts")
		}

		parts[i] = b[:partSize]
		b = b[partSize:]
	}
	return parts, nil
}
