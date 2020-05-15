with import ./nix {};
pkgs.mkShell {
  buildInputs = [
    cacert
    watchexec
    niv
    euphenix.euphenix
  ];
}
