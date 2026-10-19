# Memory Hygiene and Signer Audit

## 📘 Overview

This document provides an audit of where secret keys, seeds, and private key material are held in memory within `soroauth`, what guarantees the library provides, and what Go's runtime prevents us from guaranteeing.

---

## 🔍 Key Material Audit

| Component / Path | Where held | Who owns the buffer | Zeroed on completion? | Notes |
|---|---|---|---|---|
| **CLI Environment Read Path (`cmd/soroauth`)** | `SOROAUTH_SIGNER_SEED` or `STELLAR_SIGNER_SEED` read into `[]byte` via `os.Getenv` or prompt | The CLI process / environment parser | **Yes** (zeroed immediately after parsing or keypair construction where owned) | Most common exposure point. The CLI reads the seed string/bytes, constructs the signer or signs, and zeroizes the local byte slice. Note that `os.Getenv` returns an immutable Go `string`, which Go's runtime/garbage collector may place in read-only or managed memory segments beyond our control. |
| **`NewEd25519Signer` / SDK Signer** | `keypair.Full` from `github.com/stellar/go-stellar-sdk/keypair` | The Stellar Go SDK (`keypair.Full`) | **No** (the library does not own `keypair.Full` internals) | `keypair.Full` encapsulates the raw seed/private key bytes inside exported or unexported struct fields owned by `go-stellar-sdk`. `soroauth` receives a pointer to `keypair.Full` and invokes its `.Sign()` method. Because we do not own the internal slice inside `keypair.Full`, we cannot zeroize it. |
| **Custom / Remote Signers (KMS, Wallets)** | AWS KMS, GCP KMS, WalletSDK, RPC signers | Remote services, hardware tokens, or third-party SDKs | **N/A** | Secret keys never reside in the process memory of `soroauth`; signing operations are delegated via cryptographic RPCs or remote API calls. |

---

## 🛡️ What Go Cannot Guarantee

When reasoning about cryptographic memory hygiene in Go, we must be strictly honest about the runtime's limitations. An overclaim here is worse than silence:

1. **Garbage Collector Copies:** Go is a managed language with a moving/copying garbage collector. During stack growth, heap allocation, or compaction, the runtime may copy byte slices containing secret keys to new memory locations without zeroing the old location.
2. **String Immutability:** In Go, strings are immutable and their backing arrays can be interned or garbage collected arbitrarily. Reading secrets from environment variables (`os.Getenv`) or command-line flags yields strings whose backing memory cannot be reliably zeroed.
3. **Swap Space & OS Paging:** The operating system may page anonymous memory containing key material out to disk (swap space) or NVMe wear-leveling blocks.
4. **Core Dumps:** If the process crashes or receives an unhandled signal (e.g., `SIGABRT`, `SIGSEGV`), the OS may write a core dump containing the full process address space, including live or unzeroed key buffers.

---

## 🧪 Verification & Testing

The library includes tests (`TestSeedZeroing`) verifying that buffers owned directly by soroauth (such as temporary seed slices parsed in CLI utilities and helper functions) are explicitly zeroed out (`for i := range buf { buf[i] = 0 }`) after their lifecycle ends.
