# rustup for Go

🦀 rustup packaged as a `go install`-able module

<table align=center><td>

```go
//go:generate go tool cargo build
```

</table>

## Installation

```sh
go get -tool github.com/jcbhmr/go-rustup/cmd/...@latest
```

## Usage

```sh
go tool cargo build
```

## Development

### How it works

First, let's review the official installation guide for installing rustup:

<dl>
  <dt>Windows
  <dd>Download and run <code>rustup-init.exe</code> and follow the onscreen instructions.
  <dt>Linux & macOS
  <dd>Run <code>curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh</code> in your terminal and follow the onscreen instructions.
</dl>

The `https://sh.rustup.rs` shell script doesn't really do much; "It just does platform detection, downloads the installer and runs it." The "installer" that's being referred to can be seen in the Windows instructions: `rustup-init.exe`. The Linux & macOS shell script downloads a platform-specific `rustup-init` binary and runs it.

This `rustup-init` binary then reads its own `argv[0]` value to see that it should act like `rustup`. Then it copies itself to `<somewhere>/rustup`, `<somewhere>/cargo`, `<somewhere>/rustc`, etc. and sets up the user's `PATH` and a few other tidbits. Notice how it **copies itself** and **reads its own `argv[0]` to see how it should act**. Those `rustup`, `cargo`, `rustc`, etc. binaries are all the same binary code, just with different file names. The `rustup-init` binary is a chimara that changes its behavior based on the name of the binary.

> The rustup binary is a chimera, changing its behavior based on the
> name of the binary. This is used most prominently to enable
> Rustup's tool 'proxies' - that is, rustup itself and the rustup
> proxies are the same binary: when the binary is called 'rustup' or
> 'rustup.exe' it offers the Rustup command-line interface, and
> when it is called 'rustc' it behaves as a proxy to 'rustc'.
>
> This scheme is further used to distinguish the Rustup installer,
> called 'rustup-init', which is again just the rustup binary under a
> different name.

We use this fact to our advantage. We want `./cmd/rustup`, `./cmd/cargo`, `./cmd/rustc` Go main packages that all proxy to such an appropriately named extracted & cached official `rustup-init` binary. But we don't actually need to include 14 different copies of the same code; we can just `//go:embed` one copy of the `rustup-init` binary and then extract & rename it to match all the expected binary names.

TODO: Remove all the Git history with big binaries when I was iterating on this.

TODO: Externalize some of the internal deps.