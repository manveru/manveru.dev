with import <nixpkgs> {};
mkShell {
  buildInputs = [ go ];
  GO111MODULE = "on";
}
