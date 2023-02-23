# Demo Issuer

[![License: AGPL v3](https://img.shields.io/badge/License-AGPL_v3-blue.svg)](https://www.gnu.org/licenses/agpl-3.0)

The Demo Issuer is a light weight implementation of the issuer actor in PolygonID solution. This repo contain several example use-cases to simpley demonstrate the cummunication between the different building blocks of PolygonID, and how to construct end to end process of issuance and verification. This project is strictly for education purposes.

You can find further information on our [associated documentation](https://demoissuer.gitbook.io/demoissuer/).

## Deprecation warning!
Demo issuer is not compatible with the latest version of the protocol. Use Issuer Node instead: https://github.com/0xPolygonID/sh-id-platform


## Usage

### Prerequisites
- [Golang](https://go.dev/doc/install)
- [Ngrok](https://ngrok.com/download)
- [Make](https://www.gnu.org/software/make/)*
- [Yarn](https://classic.yarnpkg.com/)*
- [npm](https://docs.npmjs.com/downloading-and-installing-node-js-and-npm)*
- [jq](https://stedolan.github.io/jq/download/)*

*used only in the KYC demo

### Run the Age KYC demo

If you are on Unix based system, you can run the following:
1. Setup the configuration in the [config file](issuer/issuer_config.default.yaml). ([docs](https://polygon-id.gitbook.io/demoissuer/getting-started#3.-setup-the-config))
2. Run ```make -j4``` on the root of the repo.
3. Issuer available at [localhost:3001](http://localhost:3001), Verifer available at [localhost:3002](http://localhost:3002).
4. Follow the instructions on screen along with the [KYC Demo documentation](https://demoissuer.gitbook.io/demoissuer/kyc-age-demo).

If you are on Windows or encounter issues with the one-line script, you can run it manually by following the steps in the [getting started](https://polygon-id.gitbook.io/demoissuer/kyc-age-demo#getting-started) section
### Run the Demo Issuer separately

1. Setup the configuration in the [config file](issuer/issuer_config.default.yaml). ([docs](https://polygon-id.gitbook.io/demoissuer/getting-started#3.-setup-the-config))
2. Run the demo issuer with ```go run cmd/main.go``` from the [/issuer](issuer) directory.

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

