{ pkgs }:
pkgs.mkShell {
  buildInputs = with pkgs; [
    cacert
    watchexec
    euphenix.euphenix
  ];
}
