Diego-Beta CLI Plugin
=====================

**This plugin is for CLI v6.12.4 and older. It has been replaced by [CLI v6.13.0+](https://github.com/cloudfoundry/cli/releases), which incorporates some of the commands from this plugin, and the [Diego-Enabler](https://github.com/cloudfoundry-incubator/diego-enabler) plugin, which includes the `enable-diego` and `disable-diego` commands.**

This plugin enable Diego-specific features, and also allows pushing docker image. For more detail information of running apps on Diego, see [here](https://github.com/cloudfoundry-incubator/diego-design-notes/blob/master/migrating-to-diego.md)

##Installation

#####Install from Url (v.6.8.0+)
OSX
  ```
  cf install-plugin https://github.com/cloudfoundry-incubator/diego-cli-plugin/raw/master/bin/osx/diego-beta.osx
  ```

linux64:
  ```
  cf install-plugin https://github.com/cloudfoundry-incubator/diego-cli-plugin/raw/master/bin/linux64/diego-beta.linux64
  ```

windows64:
  ```
  cf install-plugin https://github.com/cloudfoundry-incubator/diego-cli-plugin/raw/master/bin/win64/diego-beta.win64
  ```

#####Install from Binary file (v.6.7.0)


- Download the binary [`win64`](https://github.com/cloudfoundry-incubator/diego-cli-plugin/raw/master/bin/win64/diego-beta.win64) [`linux64`](https://github.com/cloudfoundry-incubator/diego-cli-plugin/raw/master/bin/linux64/diego-beta.linux64) [`osx`](https://github.com/cloudfoundry-incubator/diego-cli-plugin/raw/master/bin/osx/diego-beta.osx)
- Install plugin `$ cf install-plugin <binary_name>`


##Full Command List

| command | usage | description|
| :--------------- |:---------------| :------------|
|`enable-diego`| `cf enable-diego App_Name` |enable diego for an app|
|`disable-diego`| `cf disable-diego App_Name` |disable diego for an app|
|`has-diego-enabled`| `cf has-diego-enabled App_Name` |check if diego is enabled for an app|
|`set-health-check`|`cf set-health-check App_Name port`|set health_check_type flag to either `port` or `none`|
|`get-health-check`|`cf get-health-check App_Name`|get value of health_check_type flag for an app|
|`docker-push`|`cf docker-push APP_NAME DOCKER_IMAGE [Options]`<br><br>Usage:<br>`cf docker-push test-app user/docker_path -c run.sh`|push a docker image as an app<br><br>Options:<br>`-c`: set start command<br>`--no-start`: Do not start app after push<br>`--no-route`: Do not create route for app|
