with import ./nix/nixpkgs.nix;
pkgs.mkShell {
  buildInputs = [
    cacert
    yarn
    yarn2nix
    watchexec
    go
    gotools
    gocode
    godef
    goimports
  ];

  LOCALE_ARCHIVE = "${buildPackages.glibcLocales}/lib/locale/locale-archive";
  LC_ALL = "en_US.UTF-8";
  GO111MODULE = "on";

  shellHook = ''
    unset preHook # fix for lorri
    unset GOPATH
  '';
}
