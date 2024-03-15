package cryptoh

import "fmt"

func ExampleMD5ChecksumBase64() {
	fmt.Printf("%s", MD5ChecksumBase64([]byte("admin:Western Digital Corporation:admin")))
	// Output:
	// l+uthS0Nq/1rca4m//Yfow==
}

func ExampleMD5Checksum() {
	fmt.Printf("% 02x", MD5Checksum([]byte("admin:Western Digital Corporation:admin")))
	// Output:
	// 97 eb ad 85 2d 0d ab fd 6b 71 ae 26 ff f6 1f a3
}

func ExampleSHA1ChecksumBase64() {
	fmt.Printf("%s", SHA1ChecksumBase64([]byte("admin")))
	// Output:
	// 0DPiKuNIrrVmD8IUCuw1hQxNqZc=
}

func ExampleSHA1Checksum() {
	fmt.Printf("% 02x", SHA1Checksum([]byte("admin")))
	// Output:
	// d0 33 e2 2a e3 48 ae b5 66 0f c2 14 0a ec 35 85 0c 4d a9 97
}
