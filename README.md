# Radicle Woodpecker Addon

A [radicle](https://radicle.xyz/) forge [addon](https://woodpecker-ci.org/docs/administration/forges/addon) for [Woodpecker CI](https://woodpecker-ci.org/).

## Getting Started

### Build the addon

The minimum requirement for building the addon are:
- Go v1.21 (or newer)
- make

In order to build the **Radicle Woodpecker Addon** execute the provided makefile under project's root directory:
```bash
make build
```
This will generate the executable `/tmp/bin/radicle-woodpecker-addon` which then should be added to
`WOODPECKER_ADDON_FORGE` configuration of Woodpecker CI.

> This forge addon is compatible with Woodpecker CI v3.

### Woodpecker with Radicle Requirements

Woodpecker CI and Radicle Addon requiro the following to be set up:

* [Radicle node](https://radicle.xyz/) (v1.2.0 or newer)
* [RIT's Radicle http API rad:z3pDipgU1YJRBPfZzNsW3s31dgnAR](https://explorer.radicle.gr/nodes/seed.radicle.gr/rad:z3pDipgU1YJRBPfZzNsW3s31dgnAR) (v1.0.0 or newer)
* [Radicle CI Broker rad:zwTxygwuz5LDGBq255RA2CbNGrz8](https://app.radicle.xyz/nodes/seed.radicle.garden/rad:zwTxygwuz5LDGBq255RA2CbNGrz8) (v0.11.0 or newer)
* [Radicle Webhooks Adapter rad:z2hCUNw2T1qU31LyGy7VPEiS7BkxW](https://app.radicle.xyz/nodes/seed.radicle.garden/rad:z2hCUNw2T1qU31LyGy7VPEiS7BkxW) (v0.7.0 or newer)

Radicle Webhooks Adapter should use the following configuration:
```
WEBHOOKS_SETTINGS_DB_PATH=/path/to/.radicle/httpd/webhooks.db
```

that will allow the adapter to retrieve extra webhooks configuration registered by Woodpecker CI.

### Configuration

The Radicle Woodpecker Addon uses configuration through Environment Variables. Here is a list with the details and the
default value for each one of them:

| EnvVar                | Description                            | Default Value |
|-----------------------|----------------------------------------|---------------|
| `WOODPECKER_HOST_URL` | The URL of Woodpecker CI installation. | ""            |
| `RADICLE_URL`         | The URL of the radicle HTTP API        | ""            |
| `RADICLE_HOOK_SECRET` | The secret to use for the webhooks.    | ""            |

### Running Woodpecker CI

In order to deploy you Woodpecker CI instance please refer to the official 
[documentation](https://woodpecker-ci.org/docs/administration/getting-started) for the deployment procedure. 
The supported Woodpecker CI version is v3.X.

The Radicle Addon should be setup as an [addon forge](https://woodpecker-ci.org/docs/administration/forges/addon).
For setting up the addon forge you must add the following to your configuration:
```
WOODPECKER_ADDON_FORGE=/path/to/radicle-woodpecker-addon-binary
```

## Issues / Feedback / Questions

For any sort of inquiry or issue regarding the Radicle Woodpecker Addon please open an issue at this repo.

Please note that this addon is still in its early stages, so any feedback and issues would be very-very welcome. 
You can also reach out to us on`#integrations` over in [Radicle's Zulip](https://radicle.zulipchat.com) 
with your thoughts / comments / questions.
