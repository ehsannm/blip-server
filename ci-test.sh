#!/usr/bin/env bash

go version

skipedPackages=(
#  "pkg/admin"
)
packages=$(go list ./...);

for pkg in $packages; do
  skipTest=false;
  for spkg in "${skipedPackages[@]}"; do
    x=$(echo "$pkg" | grep -c "$spkg");
    if [ "$x" -eq 1 ]; then
      skipTest=true
    fi
  done
  if [ $skipTest = false ]; then
    x=$(go test -v "$pkg");
    if [ $? -ne 0 ]; then
      curl https://notifier.nstd.me/log/Git%20-%20River%20Test%20Error/${x}
      echo "\033[0m${x}";
      exit 1
    fi
    echo "\033[0;37m${pkg} \033[0;32mPASSED!.";
  fi
done

