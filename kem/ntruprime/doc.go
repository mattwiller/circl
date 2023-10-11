//go:generate go run gen.go

// Package ntruprime implements the  NTRU Prime
// key encapsulation mechanism (KEM) as submitted to round 3 of the NIST
// Post-Quantum Cryptography Standardization Project.
//
// https://ntruprime.cr.yp.to/nist.html
//
// Both Streamlined NTRU Prime (i.e. sntrup, or Quotient NTRU) and LPR
// (i.e. ntrulpr, the Ring-LWE cryptosystem, or Product NTRU) variants are
// provided.  As of 2023, neither variant has been shown to be superior to
// the other.
//
// The Go implementation is based on the C reference implementation.
package ntruprime
