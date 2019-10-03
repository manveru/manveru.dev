with import ./nix {};
pkgs.mkShell {
  buildInputs = [
    cacert
    watchexec
    niv
  ];

  shellHook = ''
    unset preHook # fix for lorri
  '';
}
