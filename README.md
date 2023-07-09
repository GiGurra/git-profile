# git-profile
switch between git profiles


## Installation

```
go install github.com/gigurra/git-profile@v0.0.2
```

## Usage

Make a git profiles file `~/.ssh/git-profiles.json` like so:
```
{
  "profile1": {
    "ssh_config": {
      "User": "git",
      "IdentityFile": "~/.ssh/profile1.pub",
      "IdentityAgent": "..../something/something/agent.sock"
    },
    "git_config": {
      "user.name": "Name1",
      "user.email": "email1@something.com"
    }
  },
  "ingrid": {
    "ssh_config": {
      "User": "git",
      "IdentityFile": "~/.ssh/profile2.pub",
      "IdentityAgent": "..../something/something/agent.sock"
    },
    "git_config": {
      "user.name": "Name2",
      "user.email": "email2@something.com"
    }
  }
}
```

Then you can switch between profiles like so:
```
git-profile profile1
git-profile profile2
```

Which create a new `~/.ssh/config` file and set global git config values for user.name and user.email.
```
