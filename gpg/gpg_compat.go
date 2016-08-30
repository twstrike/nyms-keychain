package main

import (
	"crypto/dsa"
	"crypto/ecdsa"
	"crypto/rsa"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
	"time"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/packet"
)

//XXX This is copied all over. It doesnt seem to be part of openpgp spec
//but seems to be the de facto standard. Should we create a "convention"
//(or gnupg) packet and put things like this in it?
func primaryIdentity(e *openpgp.Entity) *openpgp.Identity {
	var firstIdentity *openpgp.Identity
	for _, ident := range e.Identities {
		if firstIdentity == nil {
			firstIdentity = ident
		}
		if ident.SelfSignature.IsPrimaryId != nil && *ident.SelfSignature.IsPrimaryId {
			return ident
		}
	}

	return firstIdentity
}

const (
	dateFormat = "2006-01-02"
)

/*
This will output like gpg2 version 2.0.x (the format has changed on gpg 2.1.x)

pub   2048D/0x76D78F0500D026C4 2010-08-19 [expires: 2018-08-19]
			Key fingerprint = 85E3 8F69 046B 44C1 EC9F  B07B 76D7 8F05 00D0 26C4
uid                 [ultimate] GPGTools Team <team@gpgtools.org>
uid                 [ultimate] GPGMail Project Team (Official OpenPGP Key) <gpgmail-devel@lists.gpgmail.org>
uid                 [ultimate] GPGTools Project Team (Official OpenPGP Key) <gpgtools-org@lists.gpgtools.org>
uid                 [ultimate] [jpeg image of size 5871]
sub   2048g/0x07EAE49ADBCBE671 2010-08-19 [expires: 2018-08-19]
sub   4096R/0xE8A664480D9E43F5 2014-04-08 [expires: 2024-01-02]
*/
func publicKeyringFormat(w io.Writer, entities openpgp.EntityList) {
	for _, e := range entities {
		expiration := ""
		primaryIdentity := primaryIdentity(e)
		if primaryIdentity != nil && primaryIdentity.SelfSignature != nil {
			expiration = signatureExpirationFormat(primaryIdentity.SelfSignature)
		}

		fmt.Fprintf(w, "pub   %s %s\n", publicKeyFormat(e.PrimaryKey), expiration)
		fmt.Fprintf(w, "      Key fingerprint = %s\n", fingerprintFormat(e.PrimaryKey.Fingerprint))

		for _, i := range e.Identities {
			//XXX Should also show if the identity is expired or revoked
			fmt.Fprintf(w, "uid                 [ unknown] %s\n", i.Name)
		}

		for _, sub := range e.Subkeys {
			fmt.Fprintf(w, "sub   %s %s\n", publicKeyFormat(sub.PublicKey), signatureExpirationFormat(sub.Sig))
		}

		fmt.Fprintln(w, "")
	}
}

func signatureExpirationFormat(sig *packet.Signature) string {
	if sig.KeyLifetimeSecs == nil {
		return ""
	}

	lifetime := time.Duration(*sig.KeyLifetimeSecs) * time.Second
	expirationTime := sig.CreationTime.Add(lifetime)
	return fmt.Sprintf("[expires: %s]", expirationTime.Format(dateFormat))
}

func publicKeyFormat(pub *packet.PublicKey) string {
	bitLen, err := pub.BitLength()
	if err != nil {
		return ""
	}

	keyType := "?"
	switch pub.PublicKey.(type) {
	case *rsa.PublicKey:
		keyType = fmt.Sprintf("%dR", bitLen)
	case *dsa.PublicKey:
		keyType = fmt.Sprintf("%dD", bitLen)
	case *ecdsa.PublicKey:
		keyType = fmt.Sprintf("%d ec?", bitLen)
	}

	return fmt.Sprintf("%s/0x%16X %s", keyType, pub.KeyId,
		pub.CreationTime.Format(dateFormat))
}

// 579E BCB2 6C97 72CD B7A8  96F2 9737 2B21 1CAD F401
func fingerprintFormat(fp [20]byte) string {
	f := strings.ToUpper(hex.EncodeToString(fp[:]))

	ret := ""
	for i := 0; i < 10; i++ {
		if i > 0 {
			ret += " "
		}

		if i == 5 {
			ret += " "
		}

		ret += f[i*4 : i*4+4]
	}

	return ret
}
