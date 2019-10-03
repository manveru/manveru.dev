let
  pkgs = import <nixpkgs> {};
  site = import ./.;

  deploy = site: (pkgs.writeShellScriptBin "deploy.sh" ''
    set -ex

    source /keybase/private/manveru/manveru.dev.sh

    mkdir deploy
    cp -r ${site}/* deploy
    chmod u+rwx -R deploy
    netlify deploy --prod
  '');
in {
  inherit site;

  meta = {
    impure = [
      (deploy site)
    ];
  };
}
