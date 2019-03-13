#Secret-Deployer

Secret-Deployer is meant to allow developers to keep their Kubernetes secrets in source control by encrypting them with Google KMS.
It allows developers to pass a simple yaml file to encryption, store it in source control, and apply it to your Kubernetes cluster.


### Configuration

A config file is expected at `$HOME/.secret-deployer.yml`, or the `-config` flag can be passed for a different location.
The file expects the following patten

```yaml
projectId: <GOOGLE_PROJECT_ID>
location: <KEY_LOCATION.
keyringId: <KEYRING_ID>
keyId: <KEY_ID>
```
