BENS
====
Build environment security

A simple secure variable store.

Many systems have secure vaults:
* Ansible has vault
* Hashicorp has vault (as well)
* most ci systems let you store environmental variables on your build nodes
* k8s has secrets
* macOS has keychain

All of these secure variables store systems don't integrate with each other, few of them are cross platform, and all of them bring along a lot of dependencies.

I went looking, but I couldn't find a simple secure variable store that would work on both Windows and Linux without any dependences. So, I decided to build one.

Setup
-----
BENS requires a yaml configuration file to store the environment, a public and private key, and a pass file to store the password for the private key. You can generate all of these by running the `scripts\init.sh` script. Note: this script depends on openssl and bas64 so they must be installed on your system before running the script.

Running
-------
If you're working directory contains the keys and yaml files all you have to do to run bens is the following:

    bens environment

The `init.sh` generates a dummy variable in the default environment so you should see an environment formated for shell. To load it into your shell run `eval $(bens environment)`.  The `environment` commands requires the private key and pass file to run. If you don't have access to those the command will fail.

The other bens command `add` only requires the public key and yaml files. To add the variable FOO with the value "bar" you simply run the following:

    bens add BAR "bar"

Now that you've added run `environment` against to verify it was added.

Environment Variables
---------------------
* `BENS_PASS`: If set read the pass from this environmental variable, unless `--ask-pass` is specified on the command line. This environment variable isn't required, if it's unset the pass is read from the `pass.txt` file.

Limitations
-----------
* Values for environmental variables are limited to the size of an RSA block
