# runner testdata

Suggested mapping:
- hello.py / hello.run / build.run / runfile.run -> AT-010..AT-013, AT-020..AT-021
- bad-no-header.run / bad-invalid-header.run / bad-header-leading-blank.run -> AT-053..AT-055
- ok.run -> AT-056
- install.run -> AT-032, AT-057, AT-064
- script-missing-runtime.run -> AT-058
- script-invalid-outer-body.run -> AT-059
- script-duplicate-os.run -> proposed AT-059A
- script-unknown-os.run -> proposed AT-059B
- script-no-os-block.run -> proposed AT-059C
- script-side-effects.run -> proposed AT-064B
- env/*.runner.env -> AT-070..AT-085
- script-vars.run + env/vars.runner.env -> variable expansion tests
- utf8-bom.run / lf.run / crlf.run -> AT-112..AT-114
