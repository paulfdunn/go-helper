package cryptoh

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
)

// MD5Checksum provides a []byte with the MD5 hash (checksum) for the input.
func MD5Checksum(input []byte) [16]byte {
	return md5.Sum(input)
}

// MD5ChecksumBase64 provides a string with the MD5 hash (checksum) in base64 for the input.
func MD5ChecksumBase64(input []byte) string {
	s := MD5Checksum(input)
	return base64.StdEncoding.EncodeToString(s[:])
}

// SHA1Checksum provides a []byte with the MD5 hash (checksum) for the input.
func SHA1Checksum(input []byte) [20]byte {
	return sha1.Sum(input)
}

// SHA1ChecksumBase64 provides a string with the MD5 hash (checksum) in base64 for the input.
func SHA1ChecksumBase64(input []byte) string {
	s := SHA1Checksum(input)
	return base64.StdEncoding.EncodeToString(s[:])
}
