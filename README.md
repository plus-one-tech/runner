# runner

runner is a lightweight command runner for scripts and tasks.

It provides a simple way to execute source files and reusable tasks with a single command.

```bash
runner hello.py
runner build
runner
```

## Overview

`runner` is designed to make script execution and task running simple and predictable.

It focuses on three goals:

* **Unified execution** – run scripts and tasks using the same command
* **Transparency** – always show the actual command being executed
* **Minimal design** – avoid complex DSLs or heavy configuration

Unlike many task runners, `runner` does not introduce a new scripting language. Task files simply contain normal program code.

## Basic Usage

### Run a script

```bash
runner hello.py
```

### Run a task

If `build.run` exists:

```bash
runner build
```

### Run default task

If `runfile.run` exists:

```bash
runner
```

## `.run` Files

A `.run` file is a simple executable task file.

Example:

```text
#python
print("Hello Runner")
```

The first line specifies the runtime.

Supported header styles:

```text
#python
#program.py
#.py
```

The rest of the file is executed as normal program code.

## Configuration

`runner` can be configured using `runner.env`.

Example:

```text
runtime.python=python
runtime.bash=bash
runtime.node=node

ext.py=python
ext.js=node
ext.sh=bash
```

This maps file extensions and runtimes to actual commands.

## Specification

Full specification is available here:

**docs/runner-spec.md**

## Status

The project is currently in the **specification stage**.

Implementation will start after the specification stabilizes.

## License

TBD
