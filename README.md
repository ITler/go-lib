# go-lib

Go library that provides opinionated abstractions and convenience code

## Prerequisites

1. [Install Go] programming language for contributions and maintenance to the repository.
2. [Install NodeJS] and implicitly [NPM] for NodeJS based repository dependencies.
3. [Install Mage] for task automation when contributing and maintaining the repository.

## Contributing

### Install tools

1. from package.json via `npm install`
2. and if you are good with the installation methods provided via mage, which
   you could check by looking at `externalDependencies` variable
   in the [Magefile](./magefiles/mage.go), follow the instructions
   from the subsequent section

### Install in CI environment

For setting up a CI environment use:

```sh
mage ci
```

<!-- prettier-ignore-start -->
[Install Go]: https://go.dev/doc/install
[Install NodeJS]: https://nodejs.org
[Install Mage]: https://magefile.org/
[NPM]: https://docs.npmjs.com/about-npm
<!-- prettier-ignore-end -->
