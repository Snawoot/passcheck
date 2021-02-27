passcheck
=========

Securely check list of passwords against [HIBP password database](https://haveibeenpwned.com/Passwords). Check is performed by sending 5 hex digits of password SHA-1 hash to HIBP servers and seeking match in retrieved list with requested hash prefix, leveraging [K-anonymity](https://en.wikipedia.org/wiki/K-anonymity) approach.

Program accepts CSV (RFC 4180) with `login,password` pairs via STDIN. Outputs list of breached accounts via STDOUT and log via STDERR.

## Installation

#### Binary download

Pre-built binaries available on [releases](https://github.com/Snawoot/passcheck/releases/latest) page.

#### From source

Alternatively, you may install passcheck from source. Run within source directory

```
make install
```

## Extras

### scan-passwordstore.sh

Shell script which scans password saved in [`pass`](https://www.passwordstore.org/) for breached passwords. It automatically pipes all passwords from password store into passcheck utility.

Usage:

```
./scan-passwordstore.sh
```
