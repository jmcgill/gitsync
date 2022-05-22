# GitSync

GitSync is a simple application to automatically sync remote changes to a repository and execute a command once the
sync is complete.

It allows rapid edit-test cycles on embedded devices (e.g. Raspberry Pi) without needing to edit the code on the
device itself.

## Requirements
Git Sync must be able to receive Github webhooks on port 80.

To enable this:

1. Run ngrok to open a named tunnel to this device: `ngrok http --subdomain qpi 80`
2. Add this ngrok domain as a Webhook destination for the repository you're working on
3. Create a Procfile in the root directory of the git repository
3. Run `gitsync` from the root directory of the git repository

IMAGE HERE

## Procfiles
As on Heroku, the Procfile is used to instruct GitSync how to re-run the application after a sync is completed.

For a python application, the Procfile might look like this:

```
python3 main.py
```
