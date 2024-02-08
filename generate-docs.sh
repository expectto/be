#!/bin/bash


# I'm not satisfied with how godocdown generates READMEs by default
# but i'm ok with it for now. To be improved later

godocdown be_ctx > be_ctx/README.md
godocdown be_http > be_http/README.md
godocdown be_json > be_json/README.md
godocdown be_jwt > be_jwt/README.md
godocdown be_math > be_math/README.md
godocdown be_reflected > be_reflected/README.md
godocdown be_string > be_string/README.md
godocdown be_time > be_time/README.md
godocdown be_url > be_url/README.md
godocdown . > core-be-matchers.md