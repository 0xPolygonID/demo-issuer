# Demo Issuer

[![License: AGPL v3](https://img.shields.io/badge/License-AGPL_v3-blue.svg)](https://www.gnu.org/licenses/agpl-3.0)

This project is a light weight implementation of the issuer with several examples use cases to simple demonstrate the different building block of PolygonID and how to build and end to end process of issuance to verification in the simplest terms. This project is strictly for education purposes.

You can find further information on our [associated documentation](https://demoissuer.gitbook.io/demoissuer/).

## Usage

### Prerequisites
- [Golang](https://go.dev/doc/install)
- [Ngrok](https://ngrok.com/download)
- [Make](https://www.gnu.org/software/make/)*
- [Yarn](https://classic.yarnpkg.com/)*
- [Npm](https://docs.npmjs.com/downloading-and-installing-node-js-and-npm)*

*used only in the KYC demo

### Run the Age KYC demo

- Run ```make``` on the root of the repo to run the age KYC demo.
- Follow the instructions on screen.

### Run the Demo Issuer separately

- Setup the configuration in the [config file](issuer/issuer_config.default.yaml).
- Run the demo issuer with ```go run cmd/issuer/main.go``` from the [/issuer](issuer) directory.

[//]: # (### Run issuer/verifier webpage separately )

[//]: # (- Setup the configuration in the [config file]&#40;issuer/issuer_config.default.yaml&#41;.)

[//]: # (- The following steps should be executed for the [issuer-webpage]&#40;examples/kycAge/issuerClient&#41; and [verifier-webpage]&#40;examples/kycAge/verifierClient&#41; separately:)

[//]: # (  - Run ```yarn``` to install all dependencies)

[//]: # (  - Run the ```yarn dev```)

[//]: # (  - Open browser on deployed address &#40;[localhost:3001]&#40;https://localhost:3001&#41; for issuer webpage, or [localhost:3002]&#40;https://localhost:3002&#41; for verifier webpage.)


## Contributions
TBD


## License

Demo Issuer is released under the terms of the AGPL-3.0 license. See [LICENSE](LICENSE) for more information.


## A word of caution
This project was created primarily for education purposes. You should **NOT USE THIS CODE IN PRODUCTION SYSTEMS**.


## References

[1] [Iden3 repos](https://github.com/orgs/iden3/repositories)

[2] [Iden3 Documentation](https://docs.iden3.io/)

[3] [PolygonID docs](https://0xpolygonid.github.io/tutorials/)

