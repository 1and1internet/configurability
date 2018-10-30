# configurability

### Background
We host containers running pre-defined images which we maintain on behalf
of our customers.

The containers are often running apps like apache, nginx, php or mysql.

Our customers do not have root access or the facility to supply the 
image themselves (as we manage the images). However, they wish to be 
able to customise the configuration of these applications.

In order to allow customers to specify how their applications are
configured, we mount config maps into the image which have their
specific requirements to be applied.

This application is built into each image and runs prior to the
application being started. It uses it's own configuration (also baked
into the image) to tell it which customer config map to read and where
to apply it.

This application was originally written in python but in order to 
minimise the bloat in customer images (as often python and pip were 
installed only to satisfy support apps the customer wasn't interested 
in) it was re-written in go.

### configurability config
Files describing what needs to be configured are typically stored in
```bash
/etc/configurability/<app>.ini
```

The location of this folder can be altered by setting
```bash
CONFIGURABILITY_INTERNAL=/new location/
```

A php.ini in that folder might contain the following:
```ini
[php]
enabled = true
ini_file_path = /etc/php/7.1/fpm/php.ini
configuration_file_name = configuration-php.json
```

Where

* 'enabled' means this app should apply these customisations (if true)
* 'ini_file_path' is the file that needs editing and
* 'configuration_file_name' is the name of the file with the custom values
in

A full set of example configurability config files can be found in
```bash
testfiles/etc_configuration
```

Nginx and Apache don't specify config files as we only allow customers to
edit gzip level and the document root.

Configuration files are searched for in the folder determined by:
```bash
$CONFIGURABILITY_DIR
```

A typical 'configuration-apache2.json' file would contain:
```json
{"gzip":"3","document_root":"mysite/html"}
```

A full set of example customisation files can be found in
```bash
testfiles/customisation
```

### Building and Testing

It should be possible to type
```bash
make
```
which should create
```bash
bin/
├── configurator
└── plugins
    ├── apache2.so
    ├── basic.so
    ├── mysql.so
    ├── nginx.so
    └── php.so

```

If you want to test the configurator it will be necessary to set the
following environment variables:

```bash
export CONFIGURABILITY_DIR=testfiles/customisations
export CONF_PLUGIN_FOLDER=bin/plugins
export CONFIGURABILITY_INTERNAL=testfiles/etc_configuration
export TEST_INPUT_FOLDER=testfiles/source_config
export TEST_OUTPUT_FOLDER=testfiles/output
```

**NOTE**: Under normal operation $TEST_OUTPUT_FOLDER and
$TEST_INPUT_FOLDER should not be set.

The values in the **testfiles/customisations/** files can be edited to 
see the effect on the resultant files. Sample source files are provided
in the **testfiles/source_config** folder.

Finally, just run the configurator and compare the contents of
**/tmp/output** to the original source files and the customisations.

### Vending
Vending is done using dep for go `https://golang.github.io/dep/`
