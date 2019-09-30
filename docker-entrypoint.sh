#!/bin/sh
set -e

if [ $(echo "$1" | cut -c1) = "-" ]; then
  echo "$0: assuming arguments for dcrd"

  set -- dcrd "$@"
fi

if [ $(echo "$1" | cut -c1) = "-" ] || [ "$1" = "dcrd" ]; then
  mkdir -p "$DCRD_DATA"
  chmod 0755 "$DCRD_DATA"
  chown -R decred "$DCRD_DATA"

  echo "$0: setting appdata directory to $DCRD_DATA"

  set -- "$@" --appdata="$DCRD_DATA"
fi

if [ "$1" = "dcrd" ] || \
   [ "$1" = "dcrctl" ] || \
   [ "$1" = "gencerts" ] || \
   [ "$1" = "findcheckpoint" ] || \
   [ "$1" = "addblock" ] || \
   [ "$1" = "promptsecret" ]; then
  echo
  exec su-exec decred "$@"
fi

echo
exec "$@"
