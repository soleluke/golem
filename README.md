#Golem - a simple discord utility bot
## Setup
## Configuration
All configuration is done in config/config.json. Contents should be as follows
```
{
	//Required fields
	"token":"<DISCORD_BOT_TOKEN>",
	//Optional Fields
	"prefix":"<single character to be used as the prefix for commands>",
	//default value is !
	"commands":<array of string containing the name of all commands to be used>
	//if commands is not specified, all commands will be allowed

}
```
## Triggers
Triggers are stored in config/triggers.json with the following format:
```
{
	"trigger":"response"
}
```
Triggers are loaded at runtime and the bot must be restarted to load any changes

