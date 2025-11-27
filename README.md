# clone 💾

> [!Note]
> This is just a tiny convenience tool for when I'm cloning Gists and GitHub/Codeberg repos.

## Installation

```console
go install github.com/peterhellberg/clone@latest
```

## Usage

```console
clone -h

Usage of clone:
  -base string
    	The base directory to clone into (default "/home/peter/Code")
  -name string
    	The name to use for the project, if blank use repo name
```

> [!Tip]
> You can specify both the base directory where you keep
> your cloned repositories, as well as the name to use for
> the directory of the project being cloned.

```console
clone -base=/tmp/repos -name=trooper git@github.com:peterhellberg/clone.git
Cloning into '/tmp/repos/GitHub/peterhellberg/trooper'...
remote: Enumerating objects: 15, done.
remote: Counting objects: 100% (15/15), done.
remote: Compressing objects: 100% (11/11), done.
remote: Total 15 (delta 3), reused 15 (delta 3), pack-reused 0 (from 0)
Receiving objects: 100% (15/15), done.
Resolving deltas: 100% (3/3), done.
```

> [!Important]
> You will get an error if there is already a directory

```console
clone -name aisnake git@gist.github.com:124fb025981b9b167f942c39b87e6624.git
fatal: destination path '/home/peter/Code/GitHub/Gists/aisnake' already exists and is not an empty directory.
exit status 128
```
