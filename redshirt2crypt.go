// Program redshirt2crypt encrypts and decrypts the redshirt2 files used in
// Darwinia.
package main

/*
 * redshirt2crypt.go
 * Encrypt and decrypt Darwinia Redshirt2 files.
 * By J. Stuart McMurray
 * Created 20220515
 * Last Modified 20220515
 */

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

/* Algorithm stolen from
https://forums.introversion.co.uk/viewtopic.php?p=479470&sid=08b726f0e5eb1793912baa2e74a90ace#p479470
*/

const (
	magic  = "redshirt2"
	buflen = 1024
)

/* table contains the offsets for each encrypted character. */
var table = []int8{
	0x1f, 0x07, 0x09, 0x01,
	0x0b, 0x02, 0x05, 0x05,
	0x03, 0x11, 0x28, 0x0c,
	0x23, 0x16, 0x1b, 0x02,
}

func main() {
	flag.Usage = func() {
		fmt.Fprint(
			os.Stderr,
			`Usage: %s

Reads stdin and if it finds a redshirt2 file, decrypts to stdout.  Otherwise
redshirt2-encrypts stdin to stdout.
`,
			os.Args[0],
		)
		flag.PrintDefaults()
	}
	flag.Parse()

	/* Work out whether to encrypt or decrypt. */
	var (
		transform func(*int, []byte)
		ctx       int
	)
	rb := make([]byte, 9)
	n, err := os.Stdin.Read(rb)
	if errors.Is(err, io.EOF) { /* Not an error we care about */
		err = nil
	}
	if nil != err {
		log.Fatalf("Looking for redshirt2 magic: %s", err)
	}
	if magic == string(rb) { /* Got the magic, decrypt. */
		transform = decrypt
	} else { /* Encrypting, write the header and first bit. */
		transform = encrypt
		if _, err := fmt.Printf("%s", magic); nil != err {
			log.Fatalf("Writing magic: %s", err)
		}
		encrypt(&ctx, rb[:n])
		if _, err := os.Stdout.Write(rb[:n]); nil != err {
			log.Fatalf("Writing first chunk: %s", err)
		}
	}

	/* Read chunks of the file, transform, and write back. */
	buf := make([]byte, buflen) /* TODO: Unhardcode. */
	for {
		/* Grab a chunk. */
		n, err = os.Stdin.Read(buf)
		if 0 == n {
			break
		}
		/* Crypt it. */
		transform(&ctx, buf[:n])
		os.Stdout.Write(buf[:n])
	}

	/* EOFs are normal */
	if errors.Is(err, io.EOF) {
		return
	}
	log.Printf("Read error: %s", err)
}

/* decrypt unredshirt2's the bytes in b. */
func decrypt(ctx *int, b []byte) {
	for i, v := range b {
		l := int8(v)
		/* Only works on printable characters. */
		if 0x20 >= l {
			continue
		}
		*ctx++
		l -= table[*ctx%len(table)]
		if 0x20 > l {
			l += 0x5f
		}
		b[i] = byte(l)
	}
}

/* encrypt applies redshirt2's algorithm to b. */
func encrypt(ctx *int, b []byte) {
	for i, v := range b {
		l := int8(v)
		/* Only works on printable characters. */
		if 0x20 >= l {
			continue
		}
		*ctx++
		l += table[*ctx%len(table)]
		if 0x20 > l {
			l -= 0x5f
		}
		b[i] = byte(l)
	}
}
