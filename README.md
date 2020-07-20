<p align="center">
    <a alt="GoReport" href="https://goreportcard.com/report/github.com/otaviof/hosts">
        <img alt="GoReport" src="https://goreportcard.com/badge/github.com/otaviof/hosts">
    </a>
    <a href="https://godoc.org/github.com/otaviof/hosts/pkg/hosts">
        <img alt="GoDoc Reference" src="https://godoc.org/github.com/otaviof/hosts/pkg/hosts?status.svg">
    </a>
    <a alt="CI Status" href="https://travis-ci.com/otaviof/hosts">
        <img alt="CI Status" src="https://travis-ci.com/otaviof/hosts.svg?branch=master">
    </a>
</p>

# `hosts`

Command line utility to generate your `/etc/hosts`, based in a combination files. It also supports
reading an external HTTP resource to populate local definitions, like block-lists for instance.

The objective of `hosts` is to allow keeping individual projects, or contexts, hosts definitions
into their own files, and also move this type of data back to user home.

## Installing

The easiest way to install `hosts` is via `go get`:

```sh
go get -u github.com/otaviof/hosts/cmd/hosts
```

Alternatively, when cloning this repository, execute:

```sh
make install
```

## Configuration

The following example configuration has fields description as comments, please consider.

```yml
---
hosts:
  input:
    # external data sources
    sources:
      # input name
      - name: uBlockOrigin
        # external resource location
        uri: https://github.com/uBlockOrigin/uAssets/raw/master/thirdparties/www.malwaredomainlist.com/hostslist/hosts.txt
        # destination file
        file:  99-blocks.host
    transformations:
      # search by regular-expression
      - search: "127.0.0.1"
        # replace with string
        replace: "0.0.0.0"
      # search without replace blocks, are skipped
      - search: ^.*?(local|localhost|broadcasthost|ip6).*?$
      - search: ^\s+#.*?$
  # output files generated
  output:
    # output name
    - name: etc-hosts
      # output file path
      path: /etc/hosts
      # without files matching regular expression
      without: 99-blocks.*
    - name: dnsmasq-blocks
      path: /etc/dnsmasq.d/blocks.conf
      # only with files matching regular expression
      with: 99-blocks.*
```

The default location the configuration file is at `~/.hosts/hosts.yaml` (`~/.hosts`), or
alternatively you can employ `--base-dir` parameter to inform a different base directory.

To start, copy the example configuration:

```sh
mkdir ~/.hosts
cp -v test/hosts-dir/hosts.yaml ~/.hosts
${EDITOR} ~/.hosts/hosts.yaml
```

### Host Files

This application will look for `.host` files in the base directory location. You can find example
of those files in [`test/hosts-dir`](https://github.com/otaviof/hosts/tree/master/test/hosts-dir),
the formatting is the same [than `/etc/hosts` file](https://man7.org/linux/man-pages/man5/hosts.5.html).

For instance:

```
127.0.0.1 hostname.local hostname
```

### External Resource

It's a common use-case to map malicious or advertising related addresses in `/etc/hosts` to
`0.0.0.0` or `127.0.0.1`, therefore you block any communication from your device to those endpoints.

Online communities like for instance [SomeOneWhoCares.org](https://someonewhocares.org) and
[uBlock Assets](https://github.com/uBlockOrigin/uAssets), are providing a up-to-date list of
[hosts](https://github.com/uBlockOrigin/uAssets/tree/master/thirdparties) which users can adopt,
although, in many cases you may need modifications, and may want to skip certain entries as well.

Therefore, `hosts` provide a way to load the external resource and apply regular expression based
transformation, which can modify contents and skip certain lines. Please consider
[configuration](#configuration) section.

### DNSMasq

Alternatively, instead of generating `/etc/hosts` format, this application can generate the format
employed on [`dnsmasq.conf` files ](http://www.thekelleys.org.uk/dnsmasq/docs/dnsmasq-man.html).
The flag to turn this formatter on is under `hosts.output` as the following example:

```yml
---
hosts:
  output:
    - name: dnsmasq-blocks
      path: /etc/dnsmasq.d/blocks.conf
```

## Usage

For the command line parameters, please use `hosts --help` and inspect the sub-commands `update` and
`apply` as well.
